package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type APIName string

const (
	Champion      = "champion"
	CurrentGame   = "current-game"
	FeaturedGame  = "featured-game"
	Game          = "game"
	League        = "league"
	LolStaticData = "lol-static-data"
	LolStatus     = "lol-status"
	Match         = "match"
	MatchHistory  = "match-history"
	Stats         = "stats"
	Summoner      = "summoner"
	Team          = "team"
)

var APIVersions = map[APIName]string{
	Champion:      "1.2",
	CurrentGame:   "1.0",
	FeaturedGame:  "1.0",
	Game:          "1.3",
	League:        "2.5",
	LolStatus:     "1.0",
	LolStaticData: "1.2",
	Match:         "2.2",
	MatchHistory:  "2.2",
	Stats:         "1.3",
	Summoner:      "1.4",
	Team:          "2.4",
}

var teamIds = []string{
	"TEAM-ae23cae0-32dd-11e4-9ce8-c81f66dba0e7",
	"TEAM-2d06b2b0-1ab3-11e5-885f-c81f66dd7106",
	"TEAM-85a717f0-1cef-11e5-85a3-c81f66dd7106",
	"TEAM-e4de4220-50a3-11e4-8eee-c81f66db96d8",
}

var summonerNames = []string{
	"YellowStar",
	"Papa Schultzz",
	"Froggen",
	"Arkanoum",
	"DominGod",
}

var summonerIds = []string{
	"20637495",
	"214132",
	"19531813",
	"50805989",
	"245111",
}

var championIds = []string{
	"102",
	"103",
}

var gameIds = []string{
	"2178836472",
}

var regionCode = "euw"
var regionPrefix = "euw.api.pvp.net"
var platformId = "EUW1"

type RequestGenerator func(args ...string) string
type RequestArgs []string
type Requests struct {
	gen  RequestGenerator
	args []RequestArgs
}

var APIRequests = map[APIName][]Requests{
	Champion: []Requests{
		Requests{
			gen: func(args ...string) string {
				return fmt.Sprintf("/api/lol/%s/v%s/champion/%s", args)
			},
			args: []RequestArgs{
				RequestArgs{regionCode, APIVersions[Champion], championIds[0]},
				RequestArgs{regionCode, APIVersions[Champion], championIds[1]},
			},
		},
	},
	// CurrentGame: []RequestGenerator{
	// 	func(args ...string) string {
	// 		return fmt.Sprintf("/observer-mode/rest/consumer/getSpectatorGameInfo/%s/%d", args)
	// 	},
	// },
	// FeaturedGame: []RequestGenerator{
	// 	func(args ...string) string {
	// 		return fmt.Sprintf("/observer-mode/rest/featured", args)
	// 	},
	// },
	Game: []Requests{
		Requests{
			gen: func(args ...string) string {
				return fmt.Sprintf("/api/lol/%s/v%s/game/by-summoner/%d/recent", args)
			},
			args: []RequestArgs{
				RequestArgs{regionCode, APIVersions[Game], summonerIds[0]},
			},
		},
	},
	League: []Requests{
		Requests{
			gen: func(args ...string) string {
				return fmt.Sprintf("/api/lol/%s/v%s/league/master", args)
			},
			args: []RequestArgs{
				RequestArgs{regionCode, APIVersions[League]},
			},
		},
		Requests{
			gen: func(args ...string) string {
				return fmt.Sprintf("/api/lol/%s/v%s/league/by-team/%s/entry", args)
			},
			args: []RequestArgs{
				RequestArgs{regionCode, APIVersions[League], teamIds[0]},
				RequestArgs{regionCode, APIVersions[League], strings.Join(teamIds, ",")},
			},
		},
		Requests{
			gen: func(args ...string) string {
				return fmt.Sprintf("/api/lol/%s/v%s/league/challenger", args)
			},
			args: []RequestArgs{
				RequestArgs{regionCode, APIVersions[League]},
			},
		},
		Requests{
			gen: func(args ...string) string {
				return fmt.Sprintf("/api/lol/%s/v%s/league/by-team/%s", args)
			},
			args: []RequestArgs{
				RequestArgs{regionCode, APIVersions[League], teamIds[0]},
				RequestArgs{regionCode, APIVersions[League], strings.Join(teamIds, ",")},
			},
		},
		Requests{
			gen: func(args ...string) string {
				return fmt.Sprintf("/api/lol/%s/v%s/league/by-summoner/%s", args)
			},
			args: []RequestArgs{
				RequestArgs{regionCode, APIVersions[League], summonerIds[0]},
				RequestArgs{regionCode, APIVersions[League], strings.Join(summonerIds, ",")},
			},
		},
		Requests{
			gen: func(args ...string) string {
				return fmt.Sprintf("/api/lol/%s/v%s/league/by-summoner/%s/entry", args)
			},
			args: []RequestArgs{
				RequestArgs{regionCode, APIVersions[League], summonerIds[0]},
				RequestArgs{regionCode, APIVersions[League], strings.Join(summonerIds, ",")},
			},
		},
	},
	Match: []Requests{
		Requests{
			gen: func(args ...string) string {
				return fmt.Sprintf("/api/lol/%s/v%s/match/%s%s")
			},
			args: []RequestArgs{
				RequestArgs{regionCode, APIVersions[Match], gameIds[0], ""},
				RequestArgs{regionCode, APIVersions[Match], gameIds[0], "?includeTimeline=true"},
			},
		},
	},
	MatchHistory: []Requests{
		Requests{
			gen: func(args ...string) string {
				return fmt.Sprintf("/api/lol/%s/v%s/matchhistory/%s%s")
			},
			args: []RequestArgs{
				RequestArgs{regionCode, APIVersions[MatchHistory], summonerIds[0], ""},
				RequestArgs{regionCode, APIVersions[MatchHistory], summonerIds[0], "?championIds=102,103&rankedQueues=RANKED_SOLO_5x5,RANKED_TEAM_5x5,RANKED_TEAM_3x3&beginIndex=0&endIndex=0"},
			},
		},
	},
	Stats: []Requests{
		Requests{
			gen: func(args ...string) string {
				return fmt.Sprintf("/api/lol/%s/v%s/stats/by-summoner/%s/ranked%s")
			},
			args: []RequestArgs{
				RequestArgs{regionCode, APIVersions[Stats], summonerIds[0], ""},
				RequestArgs{regionCode, APIVersions[Stats], summonerIds[0], "?season=SEASON2015"},
			},
		},
		Requests{
			gen: func(args ...string) string {
				return fmt.Sprintf("/api/lol/%s/v%s/stats/by-summoner/%s/summary%s")
			},
			args: []RequestArgs{
				RequestArgs{regionCode, APIVersions[Stats], summonerIds[0], ""},
				RequestArgs{regionCode, APIVersions[Stats], summonerIds[0], "?season=SEASON2015"},
			},
		},
	},
	Summoner: []Requests{
		Requests{
			gen: func(args ...string) string {
				return fmt.Sprintf("/api/lol/%s/v%s/summoner/by-name/%s")
			},
			args: []RequestArgs{
				RequestArgs{regionCode, APIVersions[Summoner], summonerNames[0]},
				RequestArgs{regionCode, APIVersions[Summoner], strings.Join(summonerNames, ",")},
			},
		},
		Requests{
			gen: func(args ...string) string {
				return fmt.Sprintf("/api/lol/%s/v%s/summoner/%s/name")
			},
			args: []RequestArgs{
				RequestArgs{regionCode, APIVersions[Summoner], summonerIds[0]},
				RequestArgs{regionCode, APIVersions[Summoner], strings.Join(summonerIds, ",")},
			},
		},
		Requests{
			gen: func(args ...string) string {
				return fmt.Sprintf("/api/lol/%s/v%s/summoner/%s/runes")
			},
			args: []RequestArgs{
				RequestArgs{regionCode, APIVersions[Summoner], summonerIds[0]},
				RequestArgs{regionCode, APIVersions[Summoner], strings.Join(summonerIds, ",")},
			},
		},
		Requests{
			gen: func(args ...string) string {
				return fmt.Sprintf("/api/lol/%s/v%s/summoner/%s/masteries")
			},
			args: []RequestArgs{
				RequestArgs{regionCode, APIVersions[Summoner], summonerIds[0]},
				RequestArgs{regionCode, APIVersions[Summoner], strings.Join(summonerIds, ",")},
			},
		},
		Requests{
			gen: func(args ...string) string {
				return fmt.Sprintf("/api/lol/%s/v%s/summoner/%s")
			},
			args: []RequestArgs{
				RequestArgs{regionCode, APIVersions[Summoner], summonerIds[0]},
				RequestArgs{regionCode, APIVersions[Summoner], strings.Join(summonerIds, ",")},
			},
		},
	},
	Team: []Requests{
		Requests{
			gen: func(args ...string) string {
				return fmt.Sprintf("/api/lol/%s/v%s/team/%s", args)
			},
			args: []RequestArgs{
				RequestArgs{regionCode, APIVersions[Team], teamIds[0]},
				RequestArgs{regionCode, APIVersions[Team], strings.Join(teamIds, ",")},
			},
		},
		Requests{
			gen: func(args ...string) string {
				return fmt.Sprintf("/api/lol/%s/v%s/team/by-summoner/%s", args)
			},
			args: []RequestArgs{
				RequestArgs{regionCode, APIVersions[Team], summonerIds[0]},
				RequestArgs{regionCode, APIVersions[Team], strings.Join(summonerIds, ",")},
			},
		},
	},
}

func Execute() error {

}

func main() {
	if err := Execute(); err != nil {
		log.Printf("Got unhabdled error: %s", err)
		os.Exit(1)
	}
}
