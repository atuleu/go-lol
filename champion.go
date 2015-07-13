package lol

import "fmt"

// ChampionID is a unique identifier for a Champion
type ChampionID int //for sure there will never even be more than a thousand champions/

// A Champion can be controlled by a Summoner in a Game
type Champion struct {
	AllyTips  []string `json:"allytips"`
	Blurb     string   `json:"blurb"`
	EnemyTips []string `json:"enemytips"`
	ID        int      `json:"id"`
	Image     ImageDto `json:"image"`
	Info      struct {
		Attack     int
		Defense    int
		Difficulty int
		Magic      int
	} `json:"info"`
	Key     string `json:"key"`
	Lore    string `json:"lore"`
	Name    string `json:"name"`
	Partype string `json:"partype"`
	Passive struct {
		Description          string   `json:"description"`
		Image                ImageDto `json:"image"`
		Name                 string   `json:"name"`
		SanitizedDescription string   `json:"sanitizedDescription"`
	} `json:"passive"`
	Recommended []struct{} `json:"recommended"`
	Skins       []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Num  int    `json:"num"`
	} `json:"skins"`
	Spells []struct{} `json:"spells"`
	Stats  struct{}   `json:"stats"`
	Tags   []string   `json:"tags"`
	Title  string     `json:"title"`
}

// GetChampion returns the champion data for the current patch
func (a *StaticAPIEndpoint) GetChampion(id ChampionID) (*Champion, error) {
	res := &Champion{}
	err := a.cachedGet(fmt.Sprintf("/champion/%d", id), map[string]string{"champData": "all"}, res)
	if err != nil {
		return nil, err
	}
	return res, err
}
