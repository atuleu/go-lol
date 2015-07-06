package xlol

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	lol ".."
)

// A ReplayDownloader is able to Download (spectate) a game currently
// played on a LoL server.
type ReplayDownloader interface {
	Download(region *lol.Region, id lol.GameID, encryptionKey string) error
}

// A ReplayGetHandler is able to server over HTTP a game, as if it would
// be on a LoL server.
type ReplayGetHandler interface {
	GetHandler(region *lol.Region, id lol.GameID) (http.Handler, string, error)
}

// A ReplayManager is both a ReplayHandler and a ReplayDownloader
type ReplayManager interface {
	ReplayDownloader
	ReplayGetHandler
}

// A LocalManager is a ReplayManager that store its data in a
// perticular location on the FileSystem
type LocalManager struct {
	datadir *replaysDataDir
}

// NewLocalManager creates a new LocalManager, who data will
// be stored in basedir
func NewLocalManager(basedir string) (*LocalManager, error) {
	res := &LocalManager{}
	var err error
	res.datadir, err = newReplaysDataDir(basedir)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *LocalManager) saveJSON(path string, v interface{}) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	_, err = io.Copy(f, bytes.NewBuffer(data))
	return err
	// enc := json.NewEncoder(f)
	// return enc.Encode(v)
}

func (m *LocalManager) downloadBinary(api *SpectateAPI, fn SpectateFunction, id int, filepath string) error {
	err := os.MkdirAll(path.Dir(filepath), 0755)
	if err != nil {
		return err
	}
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	err = api.ReadAll(fn, id, f)
	if err == nil {
		log.Printf("Saved %s", filepath)
		return nil
	}

	if rerr, ok := err.(lol.RESTError); ok == true {
		if rerr.Code == 404 {
			log.Printf("Skipped %s", filepath)
			return nil
		}
	}
	return err
}

// Download fetches from the lol spectator server data of a game
// identified by its region and ID, and save it on the local hardrive
func (m *LocalManager) Download(region *lol.Region, id lol.GameID, encryptionKey string) error {
	// d, err := newReplayDataDir(m.datadir, region, id)
	// if err != nil {
	// 	return err
	// }

	// api, err := NewSpectateAPI(region, id)
	// if err != nil {
	// 	return err
	// }

	//now we get the data, we should in a loop :
	//1. getLastChunkInfo
	//2. compute from it next time there will be chunk available
	//3. get All Chunk and Keyframe available
	//4. if last chunk is available, downlaod it and break to 6
	//5. wait until specified time, repeat from 1
	//6. Compute the metadata data to connect from game starting 0:00

	//to serve, we should :
	//serve getMetaData, making it believe that
	/*	nextChunkToDownload := 0
		nextKeyframeToDownload := 0

		replay := NewEmptyReplay()
		replay.EncryptionKey = encryptionKey

		for {
			var metadata GameMetadata
			err = api.Get(GetGameMetaData, 1, &metadata)
			if err != nil {
				return err
			}

			var cInfo LastChunkInfo
			err := api.Get(GetLastChunkInfo, 1, &cInfo)
			if err != nil {
				return err
			}

			nextAvailableChunkDate := time.Now().Add(cInfo.NextAvailableChunk.Duration() + cInfo.Duration.Duration()/10)

			replay.MergeFromMetaData(metadata)
			replay.MergeFromLastChunkInfo(cInfo)
			replay.Consolidate()

			for ; nextChunkToDownload <= int(cInfo.ID); nextChunkToDownload++ {

				chunkPath := d.chunkPath(ChunkID(nextChunkToDownload))
				if err := m.downloadBinary(api, GetGameDataChunk, nextChunkToDownload, chunkPath); err != nil {
					return err
				}
			}

			for ; nextKeyframeToDownload <= int(cInfo.AssociatedKeyFrameID); nextKeyframeToDownload++ {
				keyFramePath := d.keyFramePath(KeyFrameID(nextKeyframeToDownload))
				if err := m.downloadBinary(api, GetKeyFrame, nextKeyframeToDownload, keyFramePath); err != nil {
					return err
				}
			}

			//saves the metadata
			//erases the pending info, we recompute it at this end
			metadata.PendingAvailableChunkInfo = []ChunkInfo{}
			metadata.PendingAvailableKeyFrameInfo = []KeyFrameInfo{}
			err = m.saveJSON(d.metaDataPath(), metadata)
			if err != nil {
				return err
			}

			err = m.saveJSON(d.managerDataPath(), replayData)
			if err != nil {
				return err
			}

			if cInfo.EndGameChunkID > 0 && nextChunkToDownload > int(cInfo.EndGameChunkID) {
				log.Printf("End of game detected and reached at %d", cInfo.EndGameChunkID)
				break
			}

			cTime := time.Now()
			if cTime.After(nextAvailableChunkDate) == true {
				continue
			}
			log.Printf("Waiting until %s", nextAvailableChunkDate)
			time.Sleep(nextAvailableChunkDate.Sub(cTime))

		}
		err = replayData.check(d)
		if err != nil {
			return err
		}
		f, err := os.Create(d.endOfGameDataPath())
		if err != nil {
			return err
		}
		defer f.Close()

		return api.ReadAll(EndOfGameStats, NullParam, f)
	*/
	return nil
}

/*
// AvailableReplay parses all available replay on hardrive that are
// finished, and returns their GameMetadata, organiszed by regions
func (m *LocalManager) AvailableReplay() (map[string][]GameMetadata, error) {
	return m.datadir.allFinishedReplays()
}

type gameReplayHandler struct {
	d              *replayDataDir
	localData      *Replay
	metaData       *GameMetadata
	cinfo          LastChunkInfo
	currentChunkId ChunkID
	rx             *regexp.Regexp
}

const (
	functionIdx int = 1
	paramIdx    int = 6
	nullIdx     int = 4
)

func newGameReplayHandler(d *replayDataDir) (*gameReplayHandler, error) {
	res := &gameReplayHandler{
		d:         d,
		localData: &ReplayMetadata{},
		metaData:  &GameMetadata{},
	}
	err := res.loadJSON(d.metaDataPath(), res.metaData)
	if err != nil {
		return nil, err
	}

	err = res.loadJSON(d.managerDataPath(), res.localData)
	if err != nil {
		return nil, err
	}

	err = res.localData.check(d)
	if err != nil {
		return nil, err
	}

	for i := res.localData.StartGameChunkID; i <= res.localData.MaxChunk; i++ {
		res.currentChunkId = i
		if res.localData.chunks[i].KeyFrame > 0 {
			break
		}
	}

	res.metaData.ClientBackFetchingEnabled = true

	res.rx, err = regexp.Compile(fmt.Sprintf(`\A([a-zA-Z]+)(/%s/%d/((null)|(([0-9]+)/token)))?\z`,
		res.metaData.GameKey.PlatformID,
		res.metaData.GameKey.ID))
	return res, nil
}

func (h *gameReplayHandler) loadJSON(path string, v interface{}) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	return dec.Decode(v)
}

func (h *gameReplayHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		w.WriteHeader(404)
		return
	}

	url := req.URL.Path
	log.Printf("Handling %s", url)
	if strings.HasPrefix(url, Prefix) == false {
		w.WriteHeader(404)
		return
	}

	url = strings.TrimPrefix(url, Prefix)

	m := h.rx.FindStringSubmatch(url)
	if len(m) == 0 {
		w.WriteHeader(404)
		return
	}
	fn := m[functionIdx]
	null := m[nullIdx]
	param := m[paramIdx]

	switch SpectateFunction(fn) {
	case Version:
		h.handleVersion(null, param, w)
	case GetGameMetaData:
		h.handleGetMetaData(null, param, w)
	case GetLastChunkInfo:
		h.handleGetLastChunkInfo(null, param, w)
	case GetKeyFrame:
		h.handleGetKeyFrame(null, param, w)
	case GetGameDataChunk:
		h.handleGetChunk(null, param, w)
	case EndOfGameStats:
		h.handleGetEndOfGame(null, param, w)
	default:
		w.WriteHeader(404)
	}
}

func (h *gameReplayHandler) handleVersion(null, param string, w http.ResponseWriter) {
	if len(null) != 0 || len(param) != 0 {
		w.WriteHeader(404)
	}
	_, err := io.Copy(w, bytes.NewBuffer([]byte(h.localData.Version)))
	if err != nil {
		panic(err)
	}
}

func (h *gameReplayHandler) handleGetMetaData(null, param string, w http.ResponseWriter) {

}

func (h *gameReplayHandler) handleGetLastChunkInfo(null, param string, w http.ResponseWriter) {

}

func (h *gameReplayHandler) handleGetChunk(null, param string, w http.ResponseWriter) {

}

func (h *gameReplayHandler) handleGetKeyFrame(null, param string, w http.ResponseWriter) {

}

func (h *gameReplayHandler) handleGetEndOfGame(null, param string, w http.ResponseWriter) {

}

func (m *LocalManager) GetHandler(region *lol.Region, id lol.GameID) (http.Handler, string, error) {
	d, err := newReplayDataDir(m.datadir, region, id)
	if err != nil {
		return nil, "", err
	}

	_, err = os.Stat(d.endOfGameDataPath())
	if err != nil {
		if os.IsNotExist(err) == true {
			return nil, "", fmt.Errorf("No full data available for game %s:%d", region.Code(), id)
		}
		return nil, "", err
	}

	return nil, "", fmt.Errorf("Not yet implemented")
}
*/
