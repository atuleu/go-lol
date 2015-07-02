package lol

import "fmt"

type GameID uint64 //this is a 64 bit, EUW reached limit of int32 EUW > NA !

func (g GameID) String() string {
	return fmt.Sprintf("%d", g)
}

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
		Level                          int  `json:"level,omitempty"`
		GoldEarned                     int  `json:"goldEarned,omitempty"`
		Death                          int  `json:"numDeaths,omitempty"`
		TurretsKilled                  int  `json:"turretsKilled,omitempty"`
		MinionsKilled                  int  `json:"minionsKilled,omitempty"`
		Kills                          int  `json:"championsKilled,omitempty"`
		GoldSpent                      int  `json:"goldSpent,omitempty"`
		TotalDamageDealt               int  `json:"totalDamageDealt,omitempty"`
		TotalDamageTaken               int  `json:"totalDamageTaken,omitempty"`
		KillingSprees                  int  `json:"killingSprees,omitempty"`
		LargestKillingSpree            int  `json:"largestKillingSpree,omitempty"`
		Team                           int  `json:"team,omitempty"`
		Win                            bool `json:"win"`
		NeutralMinionsKilled           int  `json:"neutralMinionsKilled,omitempty"`
		LargestMultiKill               int  `json:"largestMultiKill,omitempty"`
		PhysicalDamageDealtPlayer      int  `json:"physicalDamageDealtPlayer,omitempty"`
		MagicDamageDealtPlayer         int  `json:"magicDamageDealtPlayer,omitempty"`
		PhysicalDamageTaken            int  `json:"physicalDamageTaken,omitempty"`
		MagicDamageTaken               int  `json:"magicDamageTaken,omitempty"`
		TimePlayed                     int  `json:"timePlayed,omitempty"`
		TotalHeal                      int  `json:"totalHeal,omitempty"`
		TotalUnitsHealed               int  `json:"totalUnitsHealed,omitempty"`
		Assists                        int  `json:"assists,omitempty"`
		Item0                          int  `json:"item0,omitempty"`
		Item1                          int  `json:"item1,omitempty"`
		Item2                          int  `json:"item2,omitempty"`
		Item3                          int  `json:"item3,omitempty"`
		Item4                          int  `json:"item4,omitempty"`
		Item5                          int  `json:"item5,omitempty"`
		Item6                          int  `json:"item6,omitempty"`
		MagicDamageDealtToChampions    int  `json:"magicDamageDealtToChampions,omitempty"`
		PhysicalDamageDealtToChampions int  `json:"physicalDamageDealtToChampions,omitempty"`
		TotalDamageDealtToChampions    int  `json:"totalDamageDealtToChampions,omitempty"`
		TrueDamageDealtPlayer          int  `json:"trueDamageDealtPlayer,omitempty"`
		TrueDamageDealtToChampions     int  `json:"trueDamageDealtToChampions,omitempty"`
		TrueDamageTaken                int  `json:"trueDamageTaken,omitempty"`
		WardKilled                     int  `json:"wardKilled,omitempty"`
		WardPlaced                     int  `json:"wardPlaced,omitempty"`
		NeutralMinionsKilledYourJungle int  `json:"neutralMinionsKilledYourJungle,omitempty"`
		TotalTimeCrowdControlDealt     int  `json:"totalTimeCrowdControlDealt,omitempty"`
		PlayerPosition                 int  `json:"playerPosition,omitempty"`
		PlayerRole                     int  `json:"playerRole,omitempty"`

		Barrackskilled                  int  `json:"barrackskilled,omitempty"`
		Combatplayerscore               int  `json:"combatPlayerScore,omitempty"`
		Consumablespurchased            int  `json:"consumablesPurchased,omitempty"`
		Damagedealtplayer               int  `json:"damageDealtPlayer,omitempty"`
		Doublekills                     int  `json:"doubleKills,omitempty"`
		Firstblood                      int  `json:"firstBlood,omitempty"`
		Gold                            int  `json:"gold,omitempty"`
		Itemspurchased                  int  `json:"itemsPurchased,omitempty"`
		Largestcriticalstrike           int  `json:"largestCriticalStrike,omitempty"`
		Legendaryitemscreated           int  `json:"legendaryItemsCreated,omitempty"`
		Minionsdenied                   int  `json:"minionsDenied,omitempty"`
		Neutralminionskilledenemyjungle int  `json:"neutralMinionsKilledEnemyJungle,omitempty"`
		Nexuskilled                     bool `json:"nexusKilled,omitempty"`
		Nodecapture                     int  `json:"nodeCapture,omitempty"`
		Nodecaptureassist               int  `json:"nodeCaptureAssist,omitempty"`
		Nodeneutralize                  int  `json:"nodeNeutralize,omitempty"`
		Nodeneutralizeassist            int  `json:"nodeNeutralizeAssist,omitempty"`
		Numitemsbought                  int  `json:"numItemsBought,omitempty"`
		Objectiveplayerscore            int  `json:"objectivePlayerScore,omitempty"`
		Pentakills                      int  `json:"pentaKills,omitempty"`
		Quadrakills                     int  `json:"quadraKills,omitempty"`
		Sightwardsbought                int  `json:"sightWardsBought,omitempty"`
		Spell1cast                      int  `json:"spell1Cast,omitempty"`
		Spell2cast                      int  `json:"spell2Cast,omitempty"`
		Spell3cast                      int  `json:"spell3Cast,omitempty"`
		Spell4cast                      int  `json:"spell4Cast,omitempty"`
		Summonspell1cast                int  `json:"summonSpell1Cast,omitempty"`
		Summonspell2cast                int  `json:"summonSpell2Cast,omitempty"`
		Supermonsterkilled              int  `json:"superMonsterKilled,omitempty"`
		Teamobjective                   int  `json:"teamObjective,omitempty"`
		Totalplayerscore                int  `json:"totalPlayerScore,omitempty"`
		Totalscorerank                  int  `json:"totalScoreRank,omitempty"`
		Triplekills                     int  `json:"tripleKills,omitempty"`
		Unrealkills                     int  `json:"unrealKills,omitempty"`
		Victorypointtotal               int  `json:"victoryPointTotal,omitempty"`
		Visionwardsbought               int  `json:"visionWardsBought,omitempty"`
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
