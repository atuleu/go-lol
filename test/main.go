package main

import (
	"fmt"
	"log"
	"os"
	"path"

	lol ".."
	xlol "../x-go-lol"
)

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

	dl, err := xlol.NewReplayDownloader(path.Join(os.TempDir(), "go-lol"))
	return dl.Download(region, fgames.Games[0].Id)

}

func main() {
	if err := Execute(); err != nil {
		log.Fatalf("Exited after error: %s\n", err)
	}
}
