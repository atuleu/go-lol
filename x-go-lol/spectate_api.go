package xlol

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/charlywazzap/go-lol"
)

// SpectateAPI is an helper to use the REST api for spectate mode
type SpectateAPI struct {
	region *lol.Region
	id     lol.GameID
	debug  bool
}

const (
	//NullParam should be passed for function without a parameter
	NullParam int = -1
)

// SpectateFunction is a function of the SpectateAPI
type SpectateFunction string

const (
	// GetGameMetaData is a function to fetch game metadata
	GetGameMetaData SpectateFunction = "getGameMetaData"
	// GetLastChunkInfo is a function that returns the LastChunkInfo
	GetLastChunkInfo SpectateFunction = "getLastChunkInfo"
	// GetGameDataChunk returns the binary compressed and encoded data for a Chunk
	GetGameDataChunk SpectateFunction = "getGameDataChunk"
	// GetKeyFrame returns the binary compressed and encoded data for a Keyframe
	GetKeyFrame SpectateFunction = "getKeyFrame"
	//EndOfGameStats returns a json data about game stats
	EndOfGameStats SpectateFunction = "endOfGameStats"
	//Version returns current version
	Version SpectateFunction = "version"
	// Prefix is the prefix of the spectate API
	Prefix string = "/observer-mode/rest/consumer/"
)

// NewSpectateAPI creates a new API endpoint dedicated to get data for
// the specified game (from lol.Region and lol.GameID)
func NewSpectateAPI(region *lol.Region, id lol.GameID) (*SpectateAPI, error) {
	if region.IsDynamic() == false {
		return nil, fmt.Errorf("SpectateAPI is only working with dynamic region")
	}

	return &SpectateAPI{
		region: region,
		id:     id,
		debug:  false,
	}, nil
}

// Format formats an URL appropriately for the API
func (a *SpectateAPI) Format(function SpectateFunction, param int) string {
	if param != NullParam {
		return fmt.Sprintf("http://%s%s%s/%s/%d/%d/token",
			a.region.SpectatorURL(),
			Prefix,
			function,
			a.region.PlatformID(),
			a.id,
			param)
	}
	return fmt.Sprintf("http://%s%s%s/%s/%d/null",
		a.region.SpectatorURL(),
		Prefix,
		function,
		a.region.PlatformID(),
		a.id)
}

// VersionURL is returning the URL for getting the SpectateAPI version
// used on the distant server.
func (a *SpectateAPI) VersionURL() string {
	return fmt.Sprintf("http://%s%s%s", a.region.SpectatorURL(), Prefix, Version)
}

// Get parses JSON data into the v param. Only GetGameMetaData,
// GetLastChunkInfo and EndOfGameStats should use it
func (a *SpectateAPI) Get(function SpectateFunction, param int, v interface{}) error {
	url := a.Format(function, param)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return lol.RESTError{Code: resp.StatusCode}
	}

	if a.debug == true {
		debugPath := path.Join(os.TempDir(),
			fmt.Sprintf("go-lol-debug.%s.%d.%d.json", function, param, time.Now().Unix()))
		f, err := os.Create(debugPath)
		if err != nil {
			return err
		}
		defer f.Close()
		log.Printf("Debugging data to %s", debugPath)
		dec := json.NewDecoder(io.TeeReader(resp.Body, f))
		return dec.Decode(v)
	}
	dec := json.NewDecoder(resp.Body)
	return dec.Decode(v)
}

// log the json data from the function. Mainly for reverse engineering purpose
func (a *SpectateAPI) logJSON(function SpectateFunction, param int) error {
	url := a.Format(function, param)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return lol.RESTError{Code: resp.StatusCode}
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var indented bytes.Buffer
	err = json.Indent(&indented, data, "", "  ")
	if err != nil {
		return err
	}

	log.Printf("%s:\n%s", url, indented.String())
	return nil
}

// ReadAll reads the entire response from the REST Api and copy it to w io.Writer
func (a *SpectateAPI) ReadAll(function SpectateFunction, param int, w io.Writer) error {
	url := a.Format(function, param)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return lol.RESTError{Code: resp.StatusCode}
	}
	_, err = io.Copy(w, resp.Body)
	return err
}

// Version is getting from the Server the version currently used.
func (a *SpectateAPI) Version() (string, error) {
	url := fmt.Sprintf("http://%s%s%s", a.region.SpectatorURL(), Prefix, Version)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return "", lol.RESTError{Code: resp.StatusCode}
	}
	d, err := ioutil.ReadAll(resp.Body)
	return string(d), err
}

func (a *SpectateAPI) readBinary(fn SpectateFunction, id int, onSuccess func([]byte) error, onUnreachable func()) error {
	var res bytes.Buffer
	if err := a.ReadAll(fn, id, &res); err != nil {
		if rerr, ok := err.(lol.RESTError); ok == true {
			if rerr.Code == http.StatusNotFound {
				onUnreachable()
				return nil
			}
		}
		return err
	}

	return onSuccess(res.Bytes())
}

// SpectateGame is spectating a Game from the SpectateAPI endpoint. It
// is fetching all data needed to spectate the Replay again and checks
// for its integrity.
func (a *SpectateAPI) SpectateGame(encryptionKey string, w ReplayDataWriter) (*Replay, error) {

	replay := NewEmptyReplay()

	//saves the encryption key. We would need this information to
	//watch the replay again.
	replay.EncryptionKey = encryptionKey
	//Get the version
	var err error
	replay.Version, err = a.Version()
	if err != nil {
		return nil, err
	}

	// next should not use a loop, but some kind of Go channel stuff,
	// but we will be waiting most of the time with a decent
	// connection.
	//
	// We skip chunk 0 and KeyFrame 0, they do not Exists!
	nextChunkToDownload := ChunkID(1)
	nextKeyframeToDownload := KeyFrameID(1)

	// Callback when a download of a chunk is successful
	onSuccessChunkDownload := func(id ChunkID) func([]byte) error {
		return func(data []byte) error {
			log.Printf("Downloaded Chunk %d", id)
			//ensure that the replay will contains the Chunk
			c := Chunk{
				ChunkInfo: ChunkInfo{
					ID: id,
				},
			}
			replay.addChunk(c)
			cIdx, ok := replay.chunksByID[id]
			if ok == false {
				return fmt.Errorf("Internal error, Chunk %d should exists", id)
			}
			//saves the data in the replay
			replay.Chunks[cIdx].data = data
			if w != nil {
				return replay.saveChunk(w, replay.Chunks[cIdx])
			}
			return nil
		}
	}
	onUnreachableChunk := func(id ChunkID) func() {
		return func() {
			log.Printf("Skips download of Chunk %d, server returned 404", id)
		}
	}
	// Callback when a download of a keyframe is successful
	onSuccessKeyFrameDownload := func(id KeyFrameID) func([]byte) error {
		return func(data []byte) error {
			log.Printf("Downloaded KeyFrame %d", id)
			//ensure that the replay will contains the KeyFrame
			kf := KeyFrame{
				KeyFrameInfo: KeyFrameInfo{
					ID: id,
				},
			}
			replay.addKeyFrame(kf)
			kfIdx, ok := replay.keyframeByID[id]
			if ok == false {
				return fmt.Errorf("Internal error, KeyFrame %d should exists", id)
			}
			//saves the data in the replay
			replay.KeyFrames[kfIdx].data = data
			if w != nil {
				return replay.saveKeyFrame(w, replay.KeyFrames[kfIdx])
			}
			return nil
		}
	}
	onUnreachableKeyFrame := func(id KeyFrameID) func() {
		return func() {
			log.Printf("Skips download of KeyFrame %d, server returned HTTP 404", id)
		}
	}

	for {
		// 1. get all current metadata, and merge it in our replay structure
		var metadata GameMetadata
		err = a.Get(GetGameMetaData, 1, &metadata)
		if err != nil {
			return nil, err
		}

		var cInfo LastChunkInfo
		err := a.Get(GetLastChunkInfo, 1, &cInfo)
		if err != nil {
			return nil, err
		}

		// special case, when player are in the loading screen, the
		// Spectate API sends invalid chunk ID, so we wait 1 min to
		// poll it again, (anyway the first valid chunkInfo will make
		// us wait for the 3 min buffer for a game.
		if cInfo.ID == 0 {
			log.Printf("Received spurious LastChunkInfo %+v. waiting 1m", cInfo)
			cInfo.NextAvailableChunk = 60000 //ms
		}

		// at least wait 1s between loops
		if cInfo.NextAvailableChunk < 1000 {
			cInfo.NextAvailableChunk = 1000 //ms
		}

		nextAvailableChunkDate := time.Now().Add(cInfo.NextAvailableChunk.Duration() + cInfo.Duration.Duration()/10)

		//actually mergin the data
		replay.MergeFromMetaData(metadata)
		replay.MergeFromLastChunkInfo(cInfo)
		replay.Consolidate()

		//we save the replay we have so far
		if w != nil {
			//we save the replay
			replay.unsafeSave(w)
		}

		// We eagerly fetch all possible chunks and keyframe. Apparently chunks could be
		// missing. I do not know if they reappears after
		// sometime. Maybe ask kindly Riot ?
		for ; nextChunkToDownload <= cInfo.ID; nextChunkToDownload++ {
			if err := a.readBinary(GetGameDataChunk,
				int(nextChunkToDownload),
				onSuccessChunkDownload(nextChunkToDownload),
				onUnreachableChunk(nextChunkToDownload)); err != nil {
				return nil, err
			}
		}

		for ; nextKeyframeToDownload <= cInfo.AssociatedKeyFrameID; nextKeyframeToDownload++ {
			if err := a.readBinary(GetKeyFrame,
				int(nextKeyframeToDownload),
				onSuccessKeyFrameDownload(nextKeyframeToDownload),
				onUnreachableKeyFrame(nextKeyframeToDownload)); err != nil {
				return nil, err
			}
		}

		//checks for end of game
		if cInfo.EndGameChunkID > 0 && nextChunkToDownload > cInfo.EndGameChunkID {
			log.Printf("End of game detected and reached at %d", cInfo.EndGameChunkID)
			break
		}

		// wait until the next chunk info is available
		cTime := time.Now()
		if cTime.After(nextAvailableChunkDate) == true {
			continue
		}
		log.Printf("Waiting until %s", nextAvailableChunkDate)
		time.Sleep(nextAvailableChunkDate.Sub(cTime))
	}

	// fetching end of game stats
	var eog bytes.Buffer
	err = a.ReadAll(EndOfGameStats, NullParam, &eog)
	if err != nil {
		return nil, err
	}
	replay.endOfGameStats = eog.Bytes()

	// we check data integrity
	err = replay.check(nil)
	if err != nil {
		return nil, fmt.Errorf("New downloaded replay is inconsistent: %s", err)
	}

	// hourray !
	return replay, nil
}
