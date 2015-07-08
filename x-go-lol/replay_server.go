package xlol

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type lastChunkInfoGenerator func() LastChunkInfo

// A ReplayServer can serve over http a replay.
type ReplayServer struct {
	loader             ReplayDataLoader
	r                  *Replay
	startStreamChunk   ChunkID
	metadataRequester  chan GameMetadata
	chunkInfoRequester chan lastChunkInfoGenerator
	finish             chan struct{}
	listener           net.Listener
	TimeDivisor        DurationMs
}

// NewReplayServer initializes a ReplayServer from data located somewhere
func NewReplayServer(loader ReplayDataLoader) (*ReplayServer, error) {
	res := &ReplayServer{
		loader: loader,
	}
	if loader == nil {
		return nil, fmt.Errorf("Empty loader")
	}
	var err error
	res.r, err = LoadReplay(loader)
	if err != nil {
		return nil, err
	}

	for _, c := range res.r.Chunks {
		if int(c.ID) < res.r.MetaData.StartGameChunkID {
			continue
		}
		if c.isAssociated() == false {
			continue
		}
		res.startStreamChunk = c.ID
		break
	}

	res.TimeDivisor = 4
	return res, nil
}

// EncryptionKey returns the encryption key used to encrypt the game
// data.
func (h *ReplayServer) EncryptionKey() string {
	return h.r.EncryptionKey
}

func (h *ReplayServer) checkGameKey(parts []string) bool {
	if len(parts) < 3 {
		return false
	}
	return parts[1] == h.r.MetaData.GameKey.PlatformID &&
		parts[2] == fmt.Sprintf("%d", h.r.MetaData.GameKey.ID)
}

func (h *ReplayServer) getParam(parts []string) (int, bool) {
	if h.checkGameKey(parts) == false {
		return 0, false
	}

	if len(parts) != 5 || parts[4] != "token" {
		return 0, false
	}
	res, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		return 0, false
	}
	return int(res), true
}

type restFunctionHandler func(*ReplayServer, []string, http.ResponseWriter, *http.Request)

func handleVersion(h *ReplayServer, parts []string, w http.ResponseWriter, req *http.Request) {
	w.Header()["Content-Type"] = []string{"text/plain"}
	_, err := io.Copy(w, bytes.NewBuffer([]byte(h.r.Version)))
	if err != nil {
		panic(err)
	}
}

func (h *ReplayServer) generateMetadata(currentChunk ChunkID) GameMetadata {
	res := h.r.MetaData
	res.ClientBackFetchingEnabled = true
	var kfid KeyFrameID = -1
	for cid := h.startStreamChunk; cid <= currentChunk; cid++ {
		c := h.r.Chunks[h.r.chunksByID[cid]]
		res.PendingAvailableChunkInfo = append(res.PendingAvailableChunkInfo, c.ChunkInfo)
		if c.KeyFrame != kfid {
			kf := h.r.KeyFrames[h.r.keyframeByID[c.KeyFrame]]
			kfid = kf.ID
			res.PendingAvailableKeyFrameInfo = append(res.PendingAvailableKeyFrameInfo, kf.KeyFrameInfo)
		}
	}
	return res
}

func toDurationMs(d time.Duration) DurationMs {
	return DurationMs(d / time.Millisecond)
}

func (h *ReplayServer) generateLastChunkInfo(currentChunk ChunkID, emitDate time.Time) (time.Time, lastChunkInfoGenerator) {
	c, ok := h.r.ChunkByID(currentChunk)
	if ok == false {
		panic(fmt.Sprintf("Missing chunk %d", currentChunk))
	}

	kf, ok := h.r.KeyFrameByID(c.KeyFrame)
	if ok == false {
		panic(fmt.Sprintf("Missing keyframe %d", c.KeyFrame))
	}

	res := LastChunkInfo{
		ID:                   currentChunk,
		AssociatedKeyFrameID: c.KeyFrame,
		NextChunkID:          kf.NextChunkID,
		EndStartupChunkID:    ChunkID(h.r.MetaData.EndStartupChunkID),
		StartGameChunkID:     ChunkID(h.r.MetaData.StartGameChunkID),
		EndGameChunkID:       ChunkID(h.r.MetaData.EndGameChunkID),
		Duration:             c.Duration,
	}

	nextDate := emitDate.Add((res.Duration / h.TimeDivisor).Duration())
	return nextDate, func() LastChunkInfo {
		now := time.Now()
		res.AvailableSince = toDurationMs(now.Sub(emitDate))
		if now.Before(nextDate) {
			res.NextAvailableChunk = toDurationMs(nextDate.Sub(now))
		} else {
			res.NextAvailableChunk = 0
		}
		return res
	}
}

func handleGetMetadata(h *ReplayServer, parts []string, w http.ResponseWriter, req *http.Request) {
	_, ok := h.getParam(parts)
	if ok == false {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	//request metadata
	md, ok := <-h.metadataRequester
	if ok == false {
		//internal loop is not running
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header()["Content-Type"] = []string{"application/json"}
	data, err := json.Marshal(md)
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(w, bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
}

func handleGetLastChunkInfo(h *ReplayServer, parts []string, w http.ResponseWriter, req *http.Request) {
	_, ok := h.getParam(parts)
	if ok == false {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	ciGen, ok := <-h.chunkInfoRequester
	if ok == false {
		//internal loop is not running
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header()["Content-Type"] = []string{"application/json"}
	ci := ciGen()
	data, err := json.Marshal(ci)
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(w, bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
}

func handleGetDataChunk(h *ReplayServer, parts []string, w http.ResponseWriter, req *http.Request) {
	param, ok := h.getParam(parts)
	if ok == false {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	r, err := h.loader.OpenChunk(ChunkID(param))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer r.Close()
	w.Header()["Content-Type"] = []string{"application/octet-stream"}
	_, err = io.Copy(w, r)
	if err != nil {
		panic(err)
	}
}

func handleGetKeyFrame(h *ReplayServer, parts []string, w http.ResponseWriter, req *http.Request) {
	param, ok := h.getParam(parts)
	if ok == false {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	r, err := h.loader.OpenKeyFrame(KeyFrameID(param))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer r.Close()

	w.Header()["Content-Type"] = []string{"application/octet-stream"}
	_, err = io.Copy(w, r)
	if err != nil {
		panic(err)
	}
}

func handleGetEndOfGame(h *ReplayServer, parts []string, w http.ResponseWriter, req *http.Request) {
	if h.checkGameKey(parts) != true || len(parts) != 4 || parts[3] != "null" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	f, err := h.loader.OpenEndOfGameStats()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer f.Close()
	w.Header()["Content-Type"] = []string{"application/octet-stream"}
	_, err = io.Copy(w, f)
	if err != nil {
		panic(err)
	}
}

var restMapping = map[SpectateFunction]restFunctionHandler{
	Version:          handleVersion,
	GetGameMetaData:  handleGetMetadata,
	GetLastChunkInfo: handleGetLastChunkInfo,
	EndOfGameStats:   handleGetEndOfGame,
	GetKeyFrame:      handleGetKeyFrame,
	GetGameDataChunk: handleGetDataChunk,
}

// version > text/plain
// json : application/json
// endofgame,chunk, keyframe application/octet-stream

func (h *ReplayServer) handle(w http.ResponseWriter, req *http.Request) {
	log.Printf("Got request %s %s %s from %s", req.Proto, req.Method, req.URL.Path, req.RemoteAddr)

	URL := req.URL.Path
	if strings.HasPrefix(URL, Prefix) == false {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	parts := strings.Split(strings.TrimPrefix(URL, Prefix), "/")
	if len(parts) == 0 {
		w.WriteHeader(http.StatusNotFound)
	}
	function := SpectateFunction(parts[0])

	handler, ok := restMapping[function]
	if ok == false {
		w.WriteHeader(http.StatusNotFound)
	}

	handler(h, parts, w, req)
}

// Close stops a running ReplayServer
func (h *ReplayServer) Close() error {
	if h.listener == nil {
		return fmt.Errorf("Server is not listening")
	}
	// this will close the intern loop, making new dynamic request
	// returnning 404 and exiting, we defer it because we want it to
	// happen after listener is closed.
	defer close(h.finish)
	// This will close the listening for new connection
	return h.listener.Close()
}

func (h *ReplayServer) internLoop() {
	shouldContinue := true
	currentChunkID := h.startStreamChunk
	currentMetadata := h.generateMetadata(currentChunkID)
	c, ok := h.r.ChunkByID(currentChunkID)
	if ok == false {
		panic(fmt.Sprintf("Missing chunk %d", currentChunkID))
	}
	kf, ok := h.r.KeyFrameByID(c.KeyFrame)
	if ok == false {
		panic(fmt.Sprintf("Missing kf %d", c.KeyFrame))
	}

	startChunkInfo := LastChunkInfo{
		ID:                   c.ID,
		AvailableSince:       0,
		NextAvailableChunk:   c.Duration / h.TimeDivisor,
		AssociatedKeyFrameID: kf.ID,
		NextChunkID:          kf.NextChunkID,
		EndStartupChunkID:    ChunkID(h.r.MetaData.EndStartupChunkID),
		StartGameChunkID:     ChunkID(h.r.MetaData.StartGameChunkID),
		EndGameChunkID:       ChunkID(h.r.MetaData.EndGameChunkID),
		Duration:             c.Duration,
	}
	currentGenerator := func() LastChunkInfo {
		return startChunkInfo
	}

	//internal channel
	tick := make(chan bool)
	defer close(tick)

	bootstrapped := false
	bootstrap := func() {
		if bootstrapped == false {
			bootstrapped = true
			log.Printf("Started data increment loop")
			go func() {
				<-time.After(c.Duration.Duration() / 10)
				tick <- true
			}()
		}
	}

	for shouldContinue {
		select {
		case <-h.finish:
			// finish is closed, likely by a call to Close, so we stop
			// the loop.
			shouldContinue = false
		case h.metadataRequester <- currentMetadata:
			bootstrap()
		case h.chunkInfoRequester <- currentGenerator:
			bootstrap()
		case <-tick:
			if currentChunkID >= ChunkID(h.r.MetaData.EndGameChunkID) {
				// internal error catchup
				continue
			}
			currentChunkID++
			emitDate := time.Now()
			log.Printf("Incrementing to chunk %d", currentChunkID)
			currentMetadata = h.generateMetadata(currentChunkID)
			var nextDate time.Time
			nextDate, currentGenerator = h.generateLastChunkInfo(currentChunkID, emitDate)
			if currentChunkID < ChunkID(h.r.MetaData.EndGameChunkID) {
				//not at the end, we will increment the counter in the future
				go func() {
					<-time.After(nextDate.Sub(emitDate))
					tick <- true
				}()
			}
		}
	}

	// this will make request pending on this data exiting with a
	// http.StatusNotFound
	close(h.metadataRequester)
	close(h.chunkInfoRequester)
}

// ListenAndServe starts an http Server on the given address to show
// the replay
func (h *ReplayServer) ListenAndServe(addr string) error {
	//we must start intern loop for serving data over time
	var err error
	h.listener, err = net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	h.metadataRequester = make(chan GameMetadata)
	h.chunkInfoRequester = make(chan lastChunkInfoGenerator)
	h.finish = make(chan struct{})
	go h.internLoop()

	err = http.Serve(h.listener, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		h.handle(w, req)
	}))

	// Closing the connection will lead to this kind of nasty thing. We
	// should maybe use some kind of framework or implement a closable
	// Listener, that will catch the condition
	if err != nil && strings.HasSuffix(err.Error(), " use of closed network connection") {
		return nil
	}
	return err
}
