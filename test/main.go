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

	GameServerAddress string            `json:"gameServerAddress"`
	Port              int               `json:"port"`
	EncryptionKey     string            `json:"encryptionKey"`
	ChunkTimeInterval time.Milliseconds `json:"chunkTimeInterval"`
	StartTime         time.Time         `json:"startTime"`
	LastChunkId       int               `json:"lastChunkId"`
	LastKeyFrameId    int               `json:"lastKeyFrameId"`
	EndStartupChunkId int               `json:"endStartupChunkId"`
	DelayTime         time.Milliseconds `json: "delayTime"`

	PendingAvailableChunkInfo []struct {
		Id           int               `json:"id"`
		Duration     time.Milliseconds `json:"duration"`
		ReceivedTime time.Time         `json:"receivedTime"`
	} `json:"pendingAvailableChunkInfo"`

	PendingAvailableKeyFrameInfo []struct {
		Id           int       `json:"id"`
		ReceivedTime time.Time `json:"receivedTime"`
		NextChunkId  int       `json:"nextChunkId"`
	} `json:"pendingAvailableKeyFrameInfo"`

	KeyFrameInterval          time.Milliseconds
	DecodedEncryptionKey      string            `json:"decodedEncryptionKey"`
	StartGameChunkId          int               `json:"startGameChunkId"`
	ClientAddedLag            time.Milliseconds `json:"clientAddedLag"`
	ClientBackFetchingEnabled bool              `json:"clientBackFetchingEnabled"`
	ClientBackFetchingFreq    int               `json:"clientBackFetchingFreq"`
	InterestScore             int               `json:"interestScore"`
	FeaturedGame              bool              `json:"featuredGame"`
	CreateTime                time.Time         `json:"createTime"`
	EndGameChunkId            int               `json:"endGameChunkId"`
	EndGameKeyFrameId         int               `json:"endGameKeyFrameId"`
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
