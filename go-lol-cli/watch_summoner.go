package main

import (
	"fmt"
	"log"
	"time"

	"github.com/charlywazzap/go-lol"
	"github.com/charlywazzap/go-lol/x-go-lol"
)

type WatchSummonerCommand struct {
	Interval string `long:"interval" short:"n" description:"Interval (300s 2m30s 10m) to wait between game check, min: 10s" default:"150s"`
}

func (x *WatchSummonerCommand) Execute(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("watch-summoner require the Summoner To Watch")
	}

	sleepDuration, err := time.ParseDuration(x.Interval)
	if err != nil {
		return err
	}
	if sleepDuration < time.Duration(10)*time.Second {
		return fmt.Errorf("Interval between checks (%s) is too small", sleepDuration)
	}

	i, err := NewInteractor(options)
	if err != nil {
		return err
	}

	ids, err := i.api.GetSummonerByName(args)
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

	for {
		currentGame, err := i.api.GetCurrentGame(summoner.ID)
		if err != nil {
			return fmt.Errorf("Could not check if %s is in a game: %s", summoner.Name, err)
		}

		if currentGame != nil {
			log.Printf("%s is in game, we start download it", summoner.Name)

			api, err := xlol.NewSpectateAPI(i.region, currentGame.ID)
			if err != nil {
				return err
			}

			// ft, err := i.manager.Create(i.region, currentGame.ID)
			// if err != nil {
			// 	return err
			// }

			//spectate the game
			replay, err := api.SpectateGame(currentGame.Observer.EncryptionKey, nil)
			if err != nil {
				return err
			}

			// adds replay info
			replay.AddGameInfo(*currentGame)

			err = i.manager.Store(replay)
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

var cachedir string

func init() {

	parser.AddCommand("watch-summoner",
		"Watch a Summoner fro in Game status and download the game it plays",
		"This command regularly polls the LoL server to check if a summoner is in Game, and download the replay of the game",
		&WatchSummonerCommand{})
}
