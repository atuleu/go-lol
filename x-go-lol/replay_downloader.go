package lol

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type ReplayDownloader struct {
	region *Region
}

func NewReplayDownloader(region *Region) (*ReplayDownloader, error) {
	if len(region.platformId) == 0 || len(region.spectatorUrl) == 0 {
		return nil, fmt.Errorf("Region does not have a spectator mode (static endpoint)")
	}
	return &ReplayDownloader{
		region: region,
	}, nil
}

func (d *ReplayDownloader) Download(id GameID) error {
	url := fmt.Sprintf("http://%s/observer-mode/rest/consumer/getGameMetaData/%s/%d/0/token",
		d.region.spectatorUrl,
		d.region.platformId,
		id)
	resp, err := http.Get(url)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return RESTError{Code: resp.StatusCode}
	}

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	var out bytes.Buffer
	json.Indent(&out, data, "", "  ")
	log.Printf("%s", out.Bytes())

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
