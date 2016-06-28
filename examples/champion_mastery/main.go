package main

import (
	"fmt"
	"log"

	"github.com/charlywazzap/go-lol"
)

func Execute() error {
	log.Printf("Opening dynamic api")
	keyStorer, err := lol.NewXdgAPIKeyStorer()
	if err != nil {
		return err
	}
	key, ok := keyStorer.Get()
	if ok == false {
		return fmt.Errorf("No API key found")
	}
	region, err := lol.NewRegion(lol.EUW)
	if err != nil {
		return err
	}
	api, err := lol.NewAPIEndpoint(region, key)
	if err != nil {
		return err
	}

	log.Printf("Opening static api")
	staticAPI, err := lol.NewStaticAPIEndpoint(region, key)
	if err != nil {
		return err
	}

	//YellowStar summoner ID
	var yellowID lol.SummonerID = 20637495

	score, err := api.GetChampionMasteryScore(yellowID)
	if err != nil {
		return err
	}
	log.Printf("YellowStar has a score of : %d", score)
	champions, err := api.GetChampionMasteryTopChampions(yellowID, 5)
	if err != nil {
		return err
	}
	for _, cm := range champions {
		champion, err := staticAPI.GetChampion(cm.Champion)
		if err != nil {
			return err
		}
		log.Printf("Champion: %s, score : %d, best note: %s, lastPlayed: %s",
			champion.Name, cm.Points, cm.HighestGrade, cm.LastPlayTime.Time())
	}
	return nil
}

func main() {
	if err := Execute(); err != nil {
		log.Fatalf("Unhandled error: %s", err)
	}

}
