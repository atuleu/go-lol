package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/charlywazzap/go-lol"
	"github.com/jessevdk/go-flags"
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

var data = lol.RESTStaticData{
	ResponseByRequest: make(map[string][]byte),
	TeamIDs: []string{
		"TEAM-ae23cae0-32dd-11e4-9ce8-c81f66dba0e7",
		"TEAM-2d06b2b0-1ab3-11e5-885f-c81f66dd7106",
		"TEAM-85a717f0-1cef-11e5-85a3-c81f66dd7106",
		"TEAM-e4de4220-50a3-11e4-8eee-c81f66db96d8",
	},
	SummonerNames: []string{
		"YellowStar",
		"Papa Schultzz",
		"Froggen",
		"Arkanoum",
		"DominGod",
	},
	SummonerIDs: []string{
		"20637495",
		"214132",
		"19531813",
		"50805989",
		"245111",
	},
	ChampionIDs: []string{
		"102",
		"103",
	},
	GameIDs: []string{
		"2178836472",
		"2178791728",
		"2177565027",
		"2177332439",
	},
	RegionCode: "euw",
	Key:        "00000000-0000-0000-0000-00000000000",
}

var regionPrefix = "https://euw.api.pvp.net"
var regionPlatformId = "EUW1"

type RequestGenerator func(args ...string) string
type RequestArgs []string
type Requests struct {
	gen  RequestGenerator
	args []RequestArgs
}

func ToIF(strs []string) []interface{} {
	res := make([]interface{}, 0, len(strs))
	for _, s := range strs {
		res = append(res, s)
	}
	return res
}

var APIRequests = map[APIName][]Requests{
	Champion: []Requests{
		Requests{
			gen: func(args ...string) string {
				return fmt.Sprintf("/api/lol/%s/v%s/champion/%s", ToIF(args)...)
			},
			args: []RequestArgs{
				RequestArgs{data.RegionCode, APIVersions[Champion], data.ChampionIDs[0]},
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
				return fmt.Sprintf("/api/lol/%s/v%s/game/by-summoner/%s/recent", ToIF(args)...)
			},
			args: []RequestArgs{
				RequestArgs{data.RegionCode, APIVersions[Game], data.SummonerIDs[0]},
			},
		},
	},
	League: []Requests{
		Requests{
			gen: func(args ...string) string {
				return fmt.Sprintf("/api/lol/%s/v%s/league/master", ToIF(args)...)
			},
			args: []RequestArgs{
				RequestArgs{data.RegionCode, APIVersions[League]},
			},
		},
		Requests{
			gen: func(args ...string) string {
				return fmt.Sprintf("/api/lol/%s/v%s/league/by-team/%s/entry", ToIF(args)...)
			},
			args: []RequestArgs{
				RequestArgs{data.RegionCode, APIVersions[League], data.TeamIDs[0]},
				RequestArgs{data.RegionCode, APIVersions[League], strings.Join(data.TeamIDs, ",")},
			},
		},
		Requests{
			gen: func(args ...string) string {
				return fmt.Sprintf("/api/lol/%s/v%s/league/challenger", ToIF(args)...)
			},
			args: []RequestArgs{
				RequestArgs{data.RegionCode, APIVersions[League]},
			},
		},
		Requests{
			gen: func(args ...string) string {
				return fmt.Sprintf("/api/lol/%s/v%s/league/by-team/%s", ToIF(args)...)
			},
			args: []RequestArgs{
				RequestArgs{data.RegionCode, APIVersions[League], data.TeamIDs[0]},
				RequestArgs{data.RegionCode, APIVersions[League], strings.Join(data.TeamIDs, ",")},
			},
		},
		Requests{
			gen: func(args ...string) string {
				return fmt.Sprintf("/api/lol/%s/v%s/league/by-summoner/%s", ToIF(args)...)
			},
			args: []RequestArgs{
				RequestArgs{data.RegionCode, APIVersions[League], data.SummonerIDs[0]},
				RequestArgs{data.RegionCode, APIVersions[League], strings.Join(data.SummonerIDs, ",")},
			},
		},
		Requests{
			gen: func(args ...string) string {
				return fmt.Sprintf("/api/lol/%s/v%s/league/by-summoner/%s/entry", ToIF(args)...)
			},
			args: []RequestArgs{
				RequestArgs{data.RegionCode, APIVersions[League], data.SummonerIDs[0]},
				RequestArgs{data.RegionCode, APIVersions[League], strings.Join(data.SummonerIDs, ",")},
			},
		},
	},
	Match: []Requests{
		Requests{
			gen: func(args ...string) string {
				return fmt.Sprintf("/api/lol/%s/v%s/match/%s%s", ToIF(args)...)
			},
			args: []RequestArgs{
				RequestArgs{data.RegionCode, APIVersions[Match], data.GameIDs[0], ""},
				RequestArgs{data.RegionCode, APIVersions[Match], data.GameIDs[0], "?includeTimeline=true"},
			},
		},
	},
	MatchHistory: []Requests{
		Requests{
			gen: func(args ...string) string {
				return fmt.Sprintf("/api/lol/%s/v%s/matchhistory/%s%s", ToIF(args)...)
			},
			args: []RequestArgs{
				RequestArgs{data.RegionCode, APIVersions[MatchHistory], data.SummonerIDs[0], ""},
				RequestArgs{data.RegionCode, APIVersions[MatchHistory], data.SummonerIDs[0], "?championIds=102,103&rankedQueues=RANKED_SOLO_5x5,RANKED_TEAM_5x5,RANKED_TEAM_3x3&beginIndex=0&endIndex=0"},
			},
		},
	},
	Stats: []Requests{
		Requests{
			gen: func(args ...string) string {
				return fmt.Sprintf("/api/lol/%s/v%s/stats/by-summoner/%s/ranked%s", ToIF(args)...)
			},
			args: []RequestArgs{
				RequestArgs{data.RegionCode, APIVersions[Stats], data.SummonerIDs[0], ""},
				RequestArgs{data.RegionCode, APIVersions[Stats], data.SummonerIDs[0], "?season=SEASON2015"},
			},
		},
		Requests{
			gen: func(args ...string) string {
				return fmt.Sprintf("/api/lol/%s/v%s/stats/by-summoner/%s/summary%s", ToIF(args)...)
			},
			args: []RequestArgs{
				RequestArgs{data.RegionCode, APIVersions[Stats], data.SummonerIDs[0], ""},
				RequestArgs{data.RegionCode, APIVersions[Stats], data.SummonerIDs[0], "?season=SEASON2015"},
			},
		},
	},
	Summoner: []Requests{
		Requests{
			gen: func(args ...string) string {
				return fmt.Sprintf("/api/lol/%s/v%s/summoner/by-name/%s", ToIF(args)...)

			},
			args: []RequestArgs{
				RequestArgs{data.RegionCode, APIVersions[Summoner], data.SummonerNames[0]},
				RequestArgs{data.RegionCode, APIVersions[Summoner], strings.Join(data.SummonerNames, ",")},
			},
		},
		Requests{
			gen: func(args ...string) string {
				return fmt.Sprintf("/api/lol/%s/v%s/summoner/%s/name", ToIF(args)...)
			},
			args: []RequestArgs{
				RequestArgs{data.RegionCode, APIVersions[Summoner], data.SummonerIDs[0]},
				RequestArgs{data.RegionCode, APIVersions[Summoner], strings.Join(data.SummonerIDs, ",")},
			},
		},
		Requests{
			gen: func(args ...string) string {
				return fmt.Sprintf("/api/lol/%s/v%s/summoner/%s/runes", ToIF(args)...)
			},
			args: []RequestArgs{
				RequestArgs{data.RegionCode, APIVersions[Summoner], data.SummonerIDs[0]},
				RequestArgs{data.RegionCode, APIVersions[Summoner], strings.Join(data.SummonerIDs, ",")},
			},
		},
		Requests{
			gen: func(args ...string) string {
				return fmt.Sprintf("/api/lol/%s/v%s/summoner/%s/masteries", ToIF(args)...)
			},
			args: []RequestArgs{
				RequestArgs{data.RegionCode, APIVersions[Summoner], data.SummonerIDs[0]},
				RequestArgs{data.RegionCode, APIVersions[Summoner], strings.Join(data.SummonerIDs, ",")},
			},
		},
		Requests{
			gen: func(args ...string) string {
				return fmt.Sprintf("/api/lol/%s/v%s/summoner/%s", ToIF(args)...)
			},
			args: []RequestArgs{
				RequestArgs{data.RegionCode, APIVersions[Summoner], data.SummonerIDs[0]},
				RequestArgs{data.RegionCode, APIVersions[Summoner], strings.Join(data.SummonerIDs, ",")},
			},
		},
	},
	Team: []Requests{
		Requests{
			gen: func(args ...string) string {
				return fmt.Sprintf("/api/lol/%s/v%s/team/%s", ToIF(args)...)
			},
			args: []RequestArgs{
				RequestArgs{data.RegionCode, APIVersions[Team], data.TeamIDs[0]},
				RequestArgs{data.RegionCode, APIVersions[Team], strings.Join(data.TeamIDs, ",")},
			},
		},
		Requests{
			gen: func(args ...string) string {
				return fmt.Sprintf("/api/lol/%s/v%s/team/by-summoner/%s", ToIF(args)...)
			},
			args: []RequestArgs{
				RequestArgs{data.RegionCode, APIVersions[Team], data.SummonerIDs[0]},
				RequestArgs{data.RegionCode, APIVersions[Team], strings.Join(data.SummonerIDs, ",")},
			},
		},
	},
}

var sem = make(chan bool, 10)

func PopulateAPI(api APIName, requests []Requests) error {
	log.Printf("Will fetch %d request type from %s-%s", len(requests), api, APIVersions[api])
	for _, req := range requests {
		for _, a := range req.args {
			uri := req.gen([]string(a)...)
			var url string
			if strings.ContainsRune(uri, '?') == true {
				url = fmt.Sprintf("%s%s&api_key=", regionPrefix, uri)
			} else {
				url = fmt.Sprintf("%s%s?api_key=", regionPrefix, uri)
			}

			sem <- true // pushing a token
			log.Printf("Will try to get %s", url)
			resp, err := http.Get(url + string(opts.ApiKey))
			defer func() {
				go func() {
					time.Sleep(10 * time.Second)
					<-sem
				}()
			}()

			if err != nil {
				return err
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				return fmt.Errorf("Got non 200 code: %s", resp.StatusCode)
			}

			data.ResponseByRequest[url+string(data.Key)], err = ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}

		}
	}
	return nil
}

func Execute() error {
	opts.OutputFile = "go-lol_testdata.json"

	if _, err := flags.Parse(&opts); err != nil {
		return err
	}

	if err := opts.ApiKey.Check(); err != nil {
		return err
	}

	//steps
	// 1. fetch all from the api endpoint

	for api, requests := range APIRequests {
		PopulateAPI(api, requests)
	}

	// 2. save it to the output file

	f, err := os.Create(opts.OutputFile)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)

	return enc.Encode(data)
}

type Options struct {
	ApiKey     lol.APIKey `short:"k" long:"key" description:"API Key to access riot API" required:"true"`
	OutputFile string     `short:"o" long:"output" description:"output file"`
}

var opts Options

func main() {
	if err := Execute(); err != nil {
		log.Printf("Got unhandled error: %s", err)
		os.Exit(1)
	}
}
