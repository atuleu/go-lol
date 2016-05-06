package lol

import (
	"fmt"
	"log"
	"net/http"
)

// A ChampionMastery represents the mastery level of a payer with a
// given champion
type ChampionMastery struct {
	Champion             ChampionID       `json:"championID"`
	Level                int              `json:"championLevel"`
	Points               int              `json:"championPoints"`
	PointsSinceLastLevel int              `json:"championPointsSinceLastLevel"`
	PointsUntilNextLevel int              `json:"championPointsUntilNextLevel"`
	ChestGranted         bool             `json:"chestGranted"`
	HighestGrade         string           `json:"highestGrade"`
	LastPlayTime         EpochMillisecond `json:"lastPlayTime"`
	Player               SummonerID       `json:"playerId"`
}

func (a *APIEndpoint) formatChampionMasteryURL(url string, options map[string]string) string {
	res := fmt.Sprintf("https://%s/championmastery/location/%s%s?api_key=%s", a.region.url, a.region.platformID, url, a.key)
	for k, v := range options {
		res = fmt.Sprintf("%s&%s=%s", res, k, v)
	}
	return res
}

// GetChampionMastery returns the champion mastery of a given player for a given champion, ni
func (a *APIEndpoint) GetChampionMastery(playerID SummonerID, championID ChampionID) (*ChampionMastery, error) {

	res := &ChampionMastery{}
	url := a.formatChampionMasteryURL(fmt.Sprintf("/player/%d/champion/%d", playerID, championID), nil)
	err := a.g.Get(url, res)
	if err != nil {
		if rerr, ok := err.(RESTError); ok == true {
			if rerr.Code == http.StatusNoContent {
				return nil, nil
			}
		}
		return nil, err
	}
	return res, nil
}

// GetChampionMasteries returns all of the ChampionMastery of a given
// player, ordered by decreasing points
func (a *APIEndpoint) GetChampionMasteries(playerID SummonerID) ([]ChampionMastery, error) {
	res := []ChampionMastery{}
	url := a.formatChampionMasteryURL(fmt.Sprintf("/player/%d/champions", playerID), nil)
	err := a.g.Get(url, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetChampionMasteryScore return sthe total score of a player
func (a *APIEndpoint) GetChampionMasteryScore(playerID SummonerID) (int, error) {
	res := -1
	url := a.formatChampionMasteryURL(fmt.Sprintf("/player/%d/score", playerID), nil)
	err := a.g.Get(url, &res)
	if err != nil {
		return -1, err
	}
	return res, nil
}

// GetChampionMasteryTopChampions returns the count top champion of a player
func (a *APIEndpoint) GetChampionMasteryTopChampions(playerID SummonerID, count int) ([]ChampionMastery, error) {
	res := []ChampionMastery{}
	options := map[string]string{}
	if count != 3 {
		options["count"] = fmt.Sprintf("%d", count)
	}
	url := a.formatChampionMasteryURL(fmt.Sprintf("/player/%d/topchampions", playerID), options)
	log.Printf(url)
	err := a.g.Get(url, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
