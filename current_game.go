package lol

import (
	"fmt"
	"strings"
	"time"
)

// CurrentGame represent a Game that a Summoner is currently playing
type CurrentGame struct {
	BannedChampion []struct {
		Champion ChampionID `json:"championId"`
		PickTurn int        `json:"pickTurn"`
		Team     int        `json:"teamID"`
	} `json:"bannedChampions"`

	ID         GameID `json:"gameId"`
	GameLength int64  `json:"GameLength"`

	Participants []struct {
		ID   SummonerID `json:"summonerId"`
		Name string     `json:"summonerName"`
	} `json:"participants"`

	Observer struct {
		EncryptionKey string `json:"encryptionKey"`
	} `json:"observers"`
}

// GetCurrentGame return the CurrentGame of a Summoner identified by
// its SummonerID. It returns nil,nil if the user is not currently
// playing a game.
func (a *APIEndpoint) GetCurrentGame(id SummonerID) (*CurrentGame, error) {
	res := &CurrentGame{}
	err := a.g.Get(fmt.Sprintf("https://%s/observer-mode/rest/consumer/getSpectatorGameInfo/%s/%d?api_key=%s",
		a.region.url,
		a.region.platformID,
		id,
		a.key), res)
	if rerr, ok := err.(RESTError); ok == true {
		if rerr.Code == 404 {
			//user is just not in a game currently
			return nil, nil
		}
	}

	if err != nil {
		return res, fmt.Errorf("Could not get current game infor for %d: %s", id, err)
	}
	return res, nil
}

func (g CurrentGame) String() string {
	participantName := make([]string, 0, len(g.Participants))
	for _, v := range g.Participants {
		participantName = append(participantName, v.Name)
	}
	res := fmt.Sprintf("GameID:%d GameLength:%s Participants:[%s]", g.ID, time.Duration(g.GameLength)*time.Second, strings.Join(participantName, " , "))
	return res
}
