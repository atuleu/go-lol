package main

import (
	"fmt"
	"log"

	lol ".."
)

type GameMetadata struct {
	GameKey struct {
		Id         GameID `json:"gameId"`
		PlatformId string `json: "platformId"`
	} `json:"gameKey"`
	GameServerAddress         string            `json:"gameServerAddress"`
	Port                      int               `json:"port"`
	EncryptionKey             string            `json:"encryptionKey"`
	ChunkTimeInterval         time.Milliseconds `json:"chunkTimeInterval"`
	StartTime                 time.Time         `json:"startTime"`
	LastChunkId               int               `json:"lastChunkId"`
	LastKeyFrameId            int               `json:"lastKeyFrameId"`
	EndStartupChunkId         int               `json:"endStartupChunkId"`
	DelayTime                 time.Milliseconds `json: "delayTime"`
	PendingAvailableChunkInfo []struct {
		Id           int
		Duration     time.Milliseconds
		ReceivedTime time.Time
	}
	PendingAvailableKeyFrameInfo []struct {
		Id           int
		ReceivedTime time.Time
		NextChunkId  int
	}
	KeyFrameInterval time.Milliseconds
}

func Execute() error {
	storer, err := lol.NewXdgAPIKeyStorer()

	if err != nil {
		return err
	}

	key, ok := storer.Get()

	if ok == false {
		log.Printf("You don't have a key")
	} else {
		log.Printf("Your API key is %s\n", key)
	}

	region, err := lol.NewRegionByCode("euw")
	if err != nil {
		return err
	}

	api := lol.NewAPIRegionalEndpoint(region, key)

	fgames, err := api.GetFeaturedGames()
	if err != nil {
		return err
	}

	if len(fgames.Games) == 0 {
		return fmt.Errorf("No featured games available")
	}

	dl, err := lol.NewReplayDownloader(region)
	return dl.Download(fgames.Games[0].Id)

}

func main() {
	if err := Execute(); err != nil {
		log.Fatalf("Exited after error: %s\n", err)
	}
}
