package lol

import "fmt"

type GameID uint64 //this is a 64 bit, EUW reached limit of int32 EUW > NA !

type Game struct {
	Id         GameID           `json:"gameId"`
	Invalid    bool             `json:"invalid"`
	Mode       string           `json:"gameMode"`
	Type       string           `json:"gameType"`
	SubType    string           `json:"subType"`
	MapId      int              `json:"mapId"`
	Team       int              `json:"teamId"`
	Champion   ChampionID       `json:"championId"`
	Spell1     int              `json:"spell1"`
	Spell2     int              `json:"spell2"`
	Level      int              `json:"level"`
	IpEarned   int              `json:"ipEarned"`
	CreateDate EpochMillisecond `json:"createDate"`

	Fellows []struct {
		Summoner SummonerID `json:"summonerID"`
		Team     int        `json:"teamID"`
		Champion ChampionID `json:"championID"`
	} `json:"fellowPlayers"`

	Stats struct {
		Level                           int  `json:"level"`
		GoldEarned                      int  `json:"goldEarned"`
		Death                           int  `json:"numDeaths"`
		TurretsKilled                   int  `json:"turretsKilled"`
		MinionsKilled                   int  `json:"minionsKilled"`
		Kills                           int  `json:"championsKilled"`
		GoldSpent                       int  `json:"goldSpent"`
		TotalDamageDealt                int  `json:"totalDamageDealt"`
		TotalDamageTaken                int  `json:"totalDamageTaken"`
		KillingSprees                   int  `json:"killingSprees"`
		LargestKillingSpree             int  `json:"largestKillingSpree"`
		Team                            int  `json:"team"`
		Win                             bool `json:"win"`
		NeutralMinionsKilled            int  `json:"neutralMinionsKilled"`
		LargestMultiKill                int  `json:"largestMultiKill"`
		PhysicalDamageDealtPlayer       int  `json:"physicalDamageDealtPlayer"`
		MagicDamageDealtPlayer          int  `json:"magicDamageDealtPlayer"`
		PhysicalDamageTaken             int  `json:"physicalDamageTaken"`
		MagicDamageTaken                int  `json:"magicDamageTaken"`
		TimePlayed                      int  `json:"timePlayed"`
		TotalHeal                       int  `json:"totalHeal"`
		TotalUnitsHealed                int  `json:"totalUnitsHealed"`
		Assists                         int  `json:"assists"`
		Item0                           int  `json:"item0"`
		Item1                           int  `json:"item1"`
		Item2                           int  `json:"item2"`
		Item3                           int  `json:"item3"`
		Item4                           int  `json:"item4"`
		Item5                           int  `json:"item5"`
		Item6                           int  `json:"item6"`
		MagicDamageDealtToChampions     int  `json:"magicDamageDealtToChampions"`
		PhysicalDamageDealtToChampions  int  `json:"physicalDamageDealtToChampions"`
		TotalDamageDealtToChampions     int  `json:"totalDamageDealtToChampions"`
		TrueDamageDealtPlayer           int  `json:"trueDamageDealtPlayer"`
		TrueDamageDealtToChampions      int  `json:"trueDamageDealtToChampions"`
		WardDilled                      int  `json:"wardKilled"`
		WardPlaced                      int  `json:"wardPlaced"`
		Neutralminionskilledyourjungle  int  `json:"neutralMinionsKilledYourJungle"`
		Barrackskilled                  int  `json:"barrackskilled"`
		Combatplayerscore               int  `json:"combatPlayerScore"`
		Consumablespurchased            int  `json:"consumablesPurchased"`
		Damagedealtplayer               int  `json:"damageDealtPlayer"`
		Doublekills                     int  `json:"doubleKills"`
		Firstblood                      int  `json:"firstBlood"`
		Gold                            int  `json:"gold"`
		Itemspurchased                  int  `json:"itemsPurchased"`
		Largestcriticalstrike           int  `json:"largestCriticalStrike"`
		Legendaryitemscreated           int  `json:"legendaryItemsCreated"`
		Minionsdenied                   int  `json:"minionsDenied"`
		Neutralminionskilledenemyjungle int  `json:"neutralMinionsKilledEnemyJungle"`
		Nexuskilled                     bool `json:"nexusKilled"`
		Nodecapture                     int  `json:"nodeCapture"`
		Nodecaptureassist               int  `json:"nodeCaptureAssist"`
		Nodeneutralize                  int  `json:"nodeNeutralize"`
		Nodeneutralizeassist            int  `json:"nodeNeutralizeAssist"`
		Numitemsbought                  int  `json:"numItemsBought"`
		Objectiveplayerscore            int  `json:"objectivePlayerScore"`
		Pentakills                      int  `json:"pentaKills"`
		PlayerPosition                  int  `json:"playerPosition"`
		PlayerRole                      int  `json:"playerRole"`
		Quadrakills                     int  `json:"quadraKills"`
		Sightwardsbought                int  `json:"sightWardsBought"`
		Spell1cast                      int  `json:"spell1Cast"`
		Spell2cast                      int  `json:"spell2Cast"`
		Spell3cast                      int  `json:"spell3Cast"`
		Spell4cast                      int  `json:"spell4Cast"`
		Summonspell1cast                int  `json:"summonSpell1Cast"`
		Summonspell2cast                int  `json:"summonSpell2Cast"`
		Supermonsterkilled              int  `json:"superMonsterKilled"`
		Teamobjective                   int  `json:"teamObjective"`
		Totalplayerscore                int  `json:"totalPlayerScore"`
		Totalscorerank                  int  `json:"totalScoreRank"`
		Totaltimecrowdcontroldealt      int  `json:"totalTimeCrowdControlDealt"`
		Triplekills                     int  `json:"tripleKills"`
		Truedamagetaken                 int  `json:"trueDamageTaken"`
		Unrealkills                     int  `json:"unrealKills"`
		Victorypointtotal               int  `json:"victoryPointTotal"`
		Visionwardsbought               int  `json:"visionWardsBought"`
	} `json:"stats"`
}

type RecentGames struct {
	ID    SummonerID `json:"summonerId"`
	Games []Game     `json:"games"`
}

func (a *APIRegionalEndpoint) GetSummonerRecentGames(id SummonerID) ([]Game, error) {
	resp := &RecentGames{}
	err := a.Get(fmt.Sprintf("/v1.3/game/by-summoner/%d/recent", id), nil, resp)
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
