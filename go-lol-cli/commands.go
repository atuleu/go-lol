package main

import (
	"fmt"
	"log"
	"path"
	"time"

	"launchpad.net/go-xdg"

	lol ".."
	xlol "../x-go-lol"
)

type WatchSummonerCommand struct {
	RegionCode string `long:"region" short:"r" description:"region to use for looking up summoners" default:"euw"`
	Interval   string `long:"interval" short:"n" description:"Interval (300s 2m30s 10m) to wait between game check, min: 10s" default:"150s"`
}

func (x *WatchSummonerCommand) Execute(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("watch-summoner require the Summoner To Watch")
	}

	region, err := lol.NewRegionByCode(x.RegionCode)
	if err != nil {
		return err
	}

	sleepDuration, err := time.ParseDuration(x.Interval)
	if err != nil {
		return err
	}
	if sleepDuration < time.Duration(10)*time.Second {
		return fmt.Errorf("Interval between checks (%s) is too small", sleepDuration)
	}

	s, err := lol.NewXdgAPIKeyStorer()
	if err != nil {
		return err
	}

	key, ok := s.Get()
	if ok == false {
		return fmt.Errorf("It seems that there are no API key store, did you use set-api-key? ")
	}

	err = key.Check()
	if err != nil {
		return err
	}

	api := lol.NewAPIEndpoint(region, key)

	ids, err := api.GetSummonerByName(args)
	if err != nil {
		if rerr, ok := err.(lol.RESTError); ok == true {
			if rerr.Code != 404 {
				return err
			}
		}
		return err
	}

	if len(ids) != 1 {
		return fmt.Errorf("Could not find Summoner '%s'", args[0])
	}
	summoner := ids[0]

	cachedir, err := xdg.Cache.Ensure("go-lol/versions")
	if err != nil {
		return err
	}

	manager, err := xlol.NewLocalManager(path.Dir(cachedir))
	if err != nil {
		return err
	}

	for {
		currentGame, err := api.GetCurrentGame(summoner.ID)
		if err != nil {
			return fmt.Errorf("Could not check if %s is in a game: %s", summoner.Name, err)
		}

		if currentGame != nil {
			log.Printf("%s is in game, we start download it", summoner.Name)
			err = manager.Download(region, currentGame.ID, currentGame.Observer.EncryptionKey)
			if err != nil {
				return err
			}
		}

		// we sleep
		log.Printf("Next check for %s in-game status at %s",
			summoner.Name,
			time.Now().Add(sleepDuration))
		time.Sleep(sleepDuration)
	}

}

func init() {
	parser.AddCommand("watch-summoner",
		"Watch a Summoner fro in Game status and download the game it plays",
		"This command regularly polls the LoL server to check if a summoner is in Game, and download the replay of the game",
		&WatchSummonerCommand{})
}
