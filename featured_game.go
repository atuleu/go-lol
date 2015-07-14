package lol

import "fmt"

// FeaturedGameInfo is an information about a featured game
type FeaturedGameInfo struct {
	BannedChampion []struct {
		Champion ChampionID `json:"championId"`
		PickTurn int        `json:"pickTurn"`
		Team     int        `json:"teamID"`
	} `json:"bannedChampions"`

	ID            GameID  `json:"gameId"`
	GameLength    uint64  `json:"GameLength"`
	GameMode      string  `json:"gameMode"`
	GameQueue     QueueID `json:"gameQueueConfigId"`
	GameStartTime uint64  `json:"gameStartTime"`
	GameType      string  `json:"gameType"`
	Map           MapID   `json:"mapId"`

	Observer struct {
		EncryptionKey string `json:"encryptionKey"`
	} `json:"observers"`

	Participants []struct {
		ID          SummonerID `json:"summonerId"`
		Name        string     `json:"summonerName"`
		Bot         bool       `json:"bot"`
		Champion    ChampionID `json:"championId"`
		ProfileIcon uint64     `json:"profileIconId"`

		SummonerSpell1 SummonerSpellID `json:"spell1Id"`
		SummonerSpell2 SummonerSpellID `json:"spell2Id"`

		TeamID int64 `json:"teamId"`
	} `json:"participants"`

	Platform string `json:"platformId"`
}

// FeaturedGames is a list of games that Riot Game considers worth
// spectating
type FeaturedGames struct {
	Games          []FeaturedGameInfo `json:"gameList"`
	RefrehInterval int64              `json:"clientRefreshInterval"`
}

// GetFeaturedGames returns the currently played FeaturedGame on the region
func (a *APIEndpoint) GetFeaturedGames() (*FeaturedGames, error) {
	res := &FeaturedGames{}

	err := a.g.Get(fmt.Sprintf("https://%s/observer-mode/rest/featured?api_key=%s",
		a.region.url,
		a.key), res)
	if err != nil {
		return nil, fmt.Errorf("Could not fetch featured game on %s: %s", a.region.code, err)
	}
	return res, nil

}
