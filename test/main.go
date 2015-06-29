package main

import (
	"fmt"
	"log"
	"time"

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

	summoners, err := api.GetSummonerByName([]string{"YakaVerkyll"})
	if err != nil {
		return err
	}

	if len(summoners) != 1 {
		return fmt.Errorf("Invalid response size!")
	}
	//display matches

	games, err := api.GetSummonerRecentGames(summoners[0].Id)
	if err != nil {
		return err
	}

	for _, g := range games {
		log.Printf("%s", g)
	}

	for {
		cg, err := api.GetCurrentGame(summoners[0].Id)
		if err != nil {
			return err
		}
		if cg == nil {
			log.Printf("User is not in game")
		} else {
			log.Printf("User is in game: %v", cg)
		}
		time.Sleep(10 * time.Second)
	}

	return nil
}

func main() {
	if err := Execute(); err != nil {
		log.Fatalf("Exited after error: %s\n", err)
	}
}
