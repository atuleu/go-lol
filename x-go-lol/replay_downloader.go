package xlol

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	lol ".."
)

// A ReplayDownloader is able to download lol game replay, and save
// them on the local hardrive
type ReplayDownloader struct {
	datadir *replaysDataDir
}

// NewReplayDownloader creates a new replay downloader, whos data will
// be stored in basedir
func NewReplayDownloader(basedir string) (*ReplayDownloader, error) {
	res := &ReplayDownloader{}
	var err error
	res.datadir, err = newReplaysDataDir(basedir)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Download fetches from the lol spectator server data of a game
// identified by its region and ID, and save it on the local hardrive
func (d *ReplayDownloader) Download(region *lol.Region, id lol.GameID) error {
	_, err := newReplayDataDir(d.datadir, region, id)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s/observer-mode/rest/consumer/getGameMetaData/%s/%d/1/token",
		region.SpectatorUrl(),
		region.PlatformID(),
		id)

	resp, err := http.Get(url)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return lol.RESTError{Code: resp.StatusCode}
	}

	data, err := ioutil.ReadAll(resp.Body)

	buffer := bytes.NewBuffer(data)

	dec := json.NewDecoder(buffer)
	var indentedBuffer bytes.Buffer
	json.Indent(&indentedBuffer, data, "", "  ")
	log.Printf("%s", indentedBuffer.String())
	var metadata GameMetadata
	err = dec.Decode(&metadata)
	if err != nil {
		return err
	}

	//now we get the data, we should in a loop :
	//1. getLastChunkInfo
	//2. compute from it next time there will be chunk available
	//3. get All Chunk and Keyframe available
	//4. if last chunk is available, downlaod it and break to 6
	//5. wait until specified time, repeat from 1
	//6. Compute the metadata data to connect from game starting 0:00

	//to serve, we should :
	//serve getMetaData, making it believe that

	return nil
}
