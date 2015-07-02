package lol

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
)

type RESTStaticMock struct {
	data   RESTStaticData
	buffer bytes.Buffer
	sem    chan bool
}

func NewRESTStaticMock(data []byte) (*RESTStaticMock, error) {
	r := bytes.NewReader(data)
	dec := json.NewDecoder(r)
	res := &RESTStaticMock{
		sem: make(chan bool, 1),
	}
	err := dec.Decode(&(res.data))
	if err != nil {
		return nil, err
	}
	if len(res.data.TeamIDs) < 2 {
		return nil, fmt.Errorf("Incomplete team data")
	}

	if len(res.data.SummonerNames) < 2 || len(res.data.SummonerIDs) < 2 || len(res.data.SummonerIDs) != len(res.data.SummonerNames) {
		return nil, fmt.Errorf("Incompletes summoner data")
	}

	if len(res.data.ChampionIDs) < 2 {
		return nil, fmt.Errorf("Incomplete champions IDs")
	}

	if len(res.data.GameIDs) < 2 {
		return nil, fmt.Errorf("Incomplete Game IDs")
	}

	if len(res.data.RegionCode) == 0 {
		return nil, fmt.Errorf("missing region code")
	}

	if len(res.data.ResponseByRequest) == 0 {
		return nil, fmt.Errorf("missing static request")
	}

	return res, nil
}

func (g *RESTStaticMock) Get(url string, v interface{}) error {
	_, ok := g.data.ResponseByRequest[url]
	if ok == false {
		log.Printf("Non recognized url: %s", url)
		return RESTError{Code: 404}
	}

	g.sem <- true
	g.buffer.Reset()
	r := io.TeeReader(bytes.NewReader(g.data.ResponseByRequest[url]), &g.buffer)
	dec := json.NewDecoder(r)

	return dec.Decode(v)
}

func (g *RESTStaticMock) LastJSONData() string {
	defer func() { <-g.sem }()
	return g.buffer.String()
}

func (g *RESTStaticMock) ATeamID() string {
	return g.data.TeamIDs[0]
}

func (g *RESTStaticMock) SeveralTeamIDs() []string {
	return g.data.TeamIDs
}

func (g *RESTStaticMock) ASummonerName() string {
	return g.data.SummonerNames[0]
}

func (g *RESTStaticMock) SeveralSummonerNames() []string {
	return g.data.SummonerNames
}

func (g *RESTStaticMock) ASummonerID() string {
	return g.data.SummonerIDs[0]
}

func (g *RESTStaticMock) SeveralSummonerIDs() []string {
	return g.data.SummonerIDs
}

func (g *RESTStaticMock) AChampionID() string {
	return g.data.ChampionIDs[0]
}

func (g *RESTStaticMock) SeveralChampionIDs() []string {
	return g.data.SummonerIDs
}

func (g *RESTStaticMock) AGameID() string {
	return g.data.GameIDs[0]
}

func (g *RESTStaticMock) SeveralGameIDs() []string {
	return g.data.GameIDs
}

func (g *RESTStaticMock) RegionCode() string {
	return g.data.RegionCode
}

func (g *RESTStaticMock) Key() APIKey {
	return g.data.Key
}

var getter *RESTStaticMock
var regionTest *Region
var api *APIEndpoint

func init() {
	data, err := Asset("data/go-lol_testdata.json")
	if err != nil {
		panic(fmt.Sprintf("Could not load test data: %s", err))
	}

	getter, err = NewRESTStaticMock(data)
	if err != nil {
		panic(fmt.Sprintf("Could not parse static data: %s", err))
	}

	regionTest, err = NewRegionByCode(getter.RegionCode())
	if err != nil {
		panic(err)
	}

	api = &APIEndpoint{
		g:      getter,
		region: regionTest,
		key:    getter.Key(),
	}
}
