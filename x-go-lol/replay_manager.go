package xlol

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	lol ".."
)

// A ReplayDownloader is able to Download (spectate) a game currently
// played on a LoL server.
type ReplayDownloader interface {
	Download(region *lol.Region, id lol.GameID) error
}

// A ReplayHandler is able to server over HTTP a game, as if it would
// be on a LoL server.
type ReplayHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

// A ReplayManager is both a ReplayHandler and a ReplayDownloader
type ReplayManager interface {
	ReplayDownloader
	ReplayHandler
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

// Download fetches from the lol spectator server data of a game
// identified by its region and ID, and save it on the local hardrive
func (m *LocalManager) Download(region *lol.Region, id lol.GameID) error {
	d, err := newReplayDataDir(m.datadir, region, id)
	if err != nil {
		return err
	}

	api, err := NewSpectateAPI(region, id)
	if err != nil {
		return err
	}

	var metadata GameMetadata

	//now we get the data, we should in a loop :
	//1. getLastChunkInfo
	//2. compute from it next time there will be chunk available
	//3. get All Chunk and Keyframe available
	//4. if last chunk is available, downlaod it and break to 6
	//5. wait until specified time, repeat from 1
	//6. Compute the metadata data to connect from game starting 0:00

	//to serve, we should :
	//serve getMetaData, making it believe that
	nextChunkToDownload := 0
	nextKeyframeToDownload := 0
	i := -1
	for {
		i++

		err = api.Get(GetGameMetaData, 1, &metadata)
		if err != nil {
			return err
		}

		//saves the metadata
		mPath := d.metaDataPath()
		if i > 0 {
			mPath = fmt.Sprintf("%s.%d", mPath, i)
		}
		err = m.saveJSON(mPath, metadata)
		if err != nil {
			return err
		}

		var cInfo LastChunkInfo
		err := api.Get(GetLastChunkInfo, 1, &cInfo)
		if err != nil {
			return err
		}
		nextAvailableChunkDate := time.Now().Add(cInfo.NextAvailableChunk.Duration() + cInfo.Duration.Duration()/10)

		err = m.saveJSON(path.Join(path.Dir(mPath), fmt.Sprintf("chunkInfo.%d.json", i)), cInfo)
		if err != nil {
			return err
		}

		for ; nextChunkToDownload <= int(cInfo.ID); nextChunkToDownload++ {
			chunkPath := d.chunkPath(ChunkID(nextChunkToDownload))
			err = os.MkdirAll(path.Dir(chunkPath), 0755)
			if err != nil {
				return err
			}
			f, err := os.Create(chunkPath)
			if err != nil {
				return err
			}
			defer f.Close()

			err = api.ReadAll(GetGameDataChunk, nextChunkToDownload, f)
			if rerr, ok := err.(lol.RESTError); ok == true {
				if rerr.Code == 404 {
					log.Printf("Skipped %s", chunkPath)
					continue
				}
			}
			if err != nil {
				return err
			}
			log.Printf("Saved %s", chunkPath)
		}

		for ; nextKeyframeToDownload <= int(cInfo.AssociatedKeyFrameID); nextKeyframeToDownload++ {
			keyFramePath := d.keyFramePath(KeyFrameID(nextKeyframeToDownload))
			err = os.MkdirAll(path.Dir(keyFramePath), 0755)
			if err != nil {
				return err
			}
			f, err := os.Create(keyFramePath)
			if err != nil {
				return err
			}
			defer f.Close()
			err = api.ReadAll(GetKeyFrame, nextKeyframeToDownload, f)
			if rerr, ok := err.(lol.RESTError); ok == true {
				if rerr.Code == 404 {
					log.Printf("Skipped %s", keyFramePath)
					continue
				}
			}
			if err != nil {
				return err
			}
			log.Printf("Saved %s", keyFramePath)
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

	f, err := os.Create(d.endOfGameDataPath())
	if err != nil {
		return err
	}
	defer f.Close()
	return api.ReadAll(EndOfGameStats, NullParam, f)
}
