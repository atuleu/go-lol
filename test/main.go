package main

import (
	"log"

	lol ".."
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

	log.Printf("%d", fgames.RefrehInterval)
	for i, g := range fgames.Games {
		log.Printf("%d: {%s}", i, g)
	}

	return nil
}

func main() {
	if err := Execute(); err != nil {
		log.Fatalf("Exited after error: %s\n", err)
	}
}
