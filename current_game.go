package lol

import (
	"fmt"
	"strings"
	"time"
)

// A ProfileIconID uniquely identifies a Profile Icon
type ProfileIconID int64

// CurrentGameInfo represent a Game that a Summoner is currently playing
type CurrentGameInfo struct {
	BannedChampion []struct {
		Champion ChampionID `json:"championId"`
		PickTurn int        `json:"pickTurn"`
		Team     int        `json:"teamID"`
	} `json:"bannedChampions"`

	ID            GameID           `json:"gameId"`
	GameLength    int64            `json:"GameLength"`
	GameMode      string           `json:"gameMode"`
	GameQueue     QueueID          `json:"gameQueueConfigId"`
	GameStartTime EpochMillisecond `json:"gameStartTime"`
	GameType      string           `json:"gameType"`
	Map           MapID            `json:"mapId"`

	Observer struct {
		EncryptionKey string `json:"encryptionKey"`
	} `json:"observers"`

	Participants []struct {
		ID          SummonerID    `json:"summonerId"`
		Name        string        `json:"summonerName"`
		Bot         bool          `json:"bot"`
		Champion    ChampionID    `json:"championId"`
		ProfileIcon ProfileIconID `json:"profileIconId"`

		Masteries []struct {
			ID   MasteryID `json:"matseryId"`
			Rank int       `json:"rank"`
		} `json:"masteries"`
		Runes []struct {
			ID    RuneID `json:"runeId"`
			Count int    `json:"count"`
		} `json:"runes"`

		SummonerSpell1 SummonerSpellID `json:"spell1Id"`
		SummonerSpell2 SummonerSpellID `json:"spell2Id"`

		TeamID int64 `json:"teamId"`
	} `json:"participants"`

	Platform string `json:"platformId"`
}

// GetCurrentGame return the CurrentGame of a Summoner identified by
// its SummonerID. It returns nil,nil if the user is not currently
// playing a game.
func (a *APIEndpoint) GetCurrentGame(id SummonerID) (*CurrentGameInfo, error) {
	res := &CurrentGameInfo{}
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

func (g CurrentGameInfo) String() string {
	participantName := make([]string, 0, len(g.Participants))
	for _, v := range g.Participants {
		participantName = append(participantName, v.Name)
	}
	res := fmt.Sprintf("GameID:%d GameLength:%s Participants:[%s]", g.ID, time.Duration(g.GameLength)*time.Second, strings.Join(participantName, " , "))
	return res
}
