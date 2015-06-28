package lol

import "fmt"

type GameID uint64 //this is a 64 bit, EUW reached limit of int32 EUW > NA !

type TeamID int

type Game struct {
	CreateDate EpochMillisecond `json:"createDate"`
	Id         GameID           `json:"gameId"`
	Champion   ChampionID       `json:"championId"`

	Fellows []struct {
		Champion ChampionID `json:"championID"`
		Name     string
		Summoner SummonerID `json:"summonerID"`
		Team     TeamID     `json:"teamID"`
	} `json:"fellowPlayers"`

	Stats struct {
		Win     bool `json:"win"`
		Kills   int  `json:"championsKilled"`
		Death   int  `json:"numDeaths"`
		Assists int  `json:"assists"`
	} `json:"stats"`
}

type RecentGames struct {
	Games []Game `json:"games"`
}

func (a *APIRegionalEndpoint) GetSummonerRecentGames(id SummonerID) ([]Game, error) {
	resp := RecentGames{}
	err := a.Get(fmt.Sprintf("/v1.3/game/by-summoner/%d/recent", id), nil, &resp)
	if err != nil {
		return nil, err
	}
	//we collect all names in a second call

	return resp.Games, nil
}

func (g Game) String() string {
	return fmt.Sprintf("GameID:%d Champion Played: %d Won:%v KDA:%d/%d/%d",
		g.Id,
		g.Champion,
		g.Stats.Win,
		g.Stats.Kills,
		g.Stats.Death,
		g.Stats.Assists)
}
