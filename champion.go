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
	Stats  champStats `json:"stats"`
	Tags   []string   `json:"tags"`
	Title  string     `json:"title"`
}

type champStats struct {
	Attackrange          float32 `json:"attackrange"`
	Mpperlevel           float32 `json:"mpperlevel"`
	Mp                   float32 `json:"mp"`
	Attackdamage         float32 `json:"attackdamage"`
	Hp                   float32 `json:"hp"`
	Hpperlevel           float32 `json:"hpperlevel"`
	Attackdamageperlevel float32 `json:"attackdamageperlevel"`
	Armor                float32 `json:"armor"`
	Mpregenperlevel      float32 `json:"mpregenperlevel"`
	Hpregen              float32 `json:"hpregen"`
	Critperlevel         float32 `json:"critperlevel"`
	Spellblockperlevel   float32 `json:"spellblockperlevel"`
	Mregen               float32 `json:"mregen"`
	Attackspeedperlevel  float32 `json:"attackspeedperlevel"`
	Spellblock           float32 `json:"spellblock"`
	Movespeed            float32 `json:"movespeed"`
	Attackspeedoffset    float32 `json:"attackspeedoffset"`
	Crit                 float32 `json:"crit"`
	Hpregenperlevel      float32 `json:"hpregenperlevel"`
	Armorperlevel        float32 `json:"armorperlevel"`
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
