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

	lol ".."
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
	if len(region.PlatformID()) == 0 || len(region.SpectatorURL()) == 0 {
		return nil, fmt.Errorf("Invalid static region")
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

func (a *SpectateAPI) readBinary(fn SpectateFunction, id int, onSuccess func(), onUnreachable func()) ([]byte, error) {
	var res bytes.Buffer
	if err := a.ReadAll(fn, id, &res); err != nil {
		if rerr, ok := err.(lol.RESTError); ok == true {
			if rerr.Code == 404 {
				onUnreachable()
				return nil, nil
			}
		}
		return nil, err
	}
	onSuccess()
	return res.Bytes(), nil
}

// SpectateGame is spectating a Game from the SpectateAPI endpoint. It
// is fetching all data needed to spectate the Replay again and checks
// for its integrity.
func (a *SpectateAPI) SpectateGame(encryptionKey string) (*Replay, error) {

	replay := NewEmptyReplay()
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
	nextChunkToDownload := ChunkID(0)
	nextKeyframeToDownload := KeyFrameID(0)

	chunks := make(map[ChunkID][]byte)
	keyFrames := make(map[KeyFrameID][]byte)

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

		nextAvailableChunkDate := time.Now().Add(cInfo.NextAvailableChunk.Duration() + cInfo.Duration.Duration()/10)

		replay.MergeFromMetaData(metadata)
		replay.MergeFromLastChunkInfo(cInfo)
		replay.Consolidate()

		// thoses parts are nasty and growas by itself. Basically i
		// did not had the knowledge of how the client fetches data.
		// Therefore it would be nicer, fo a replay to work
		for ; nextChunkToDownload <= cInfo.ID; nextChunkToDownload++ {
			chunks[nextChunkToDownload], err = a.readBinary(GetGameDataChunk,
				int(nextChunkToDownload),
				func() {
					log.Printf("Downloaded Chunk %d", nextChunkToDownload)
				},
				func() {
					log.Printf("Skipped Chunk %d", nextChunkToDownload)
				})
			if err != nil {
				return nil, err
			}
			//makes sure data is downlaoded for all successful chunk download
			if len(chunks[nextChunkToDownload]) > 0 {
				c := Chunk{
					ChunkInfo: ChunkInfo{
						ID: nextChunkToDownload,
					},
				}
				replay.addChunk(c)
			}
		}

		for ; nextKeyframeToDownload <= cInfo.AssociatedKeyFrameID; nextKeyframeToDownload++ {
			keyFrames[nextKeyframeToDownload], err = a.readBinary(GetKeyFrame,
				int(nextKeyframeToDownload),
				func() {
					log.Printf("Downloaded KeyFrame %d", nextKeyframeToDownload)
				},
				func() {
					log.Printf("Skipped KeyFrame %d", nextKeyframeToDownload)
				})
			if err != nil {
				return nil, err
			}
			if len(keyFrames[nextKeyframeToDownload]) > 0 {
				kf := KeyFrame{
					KeyFrameInfo: KeyFrameInfo{
						ID: nextKeyframeToDownload,
					},
				}
				replay.addKeyFrame(kf)
			}
		}

		if cInfo.EndGameChunkID > 0 && nextChunkToDownload > cInfo.EndGameChunkID {
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

	var eog bytes.Buffer
	err = a.ReadAll(EndOfGameStats, NullParam, &eog)
	if err != nil {
		return nil, err
	}
	replay.endOfGameStats = eog.Bytes()

	// we put all the data where it is needed
	for idx, c := range replay.Chunks {
		data, ok := chunks[c.ID]
		if ok == false {
			continue
		}
		replay.Chunks[idx].data = data
	}

	for idx, kf := range replay.KeyFrames {
		data, ok := keyFrames[kf.ID]
		if ok == false {
			continue
		}
		replay.KeyFrames[idx].data = data
	}

	//we check data integrity
	err = replay.check(nil)
	if err != nil {
		return nil, fmt.Errorf("New downloaded replay is inconsistent: %s", err)
	}

	return replay, nil
}
