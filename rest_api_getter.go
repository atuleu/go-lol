package lol

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type RESTGetter interface {
	Get(url string, v interface{}) error
}

type SimpleRESTGetter struct{}

type RESTError struct {
	Code int
}

func (e RESTError) Error() string {
	if e.Code == 429 {
		return "Too Many request to server"
	}
	return fmt.Sprintf("Non 200 return code: %d", e.Code)
}

func NewSimpleRESTGetter() *SimpleRESTGetter {
	return &SimpleRESTGetter{}
}

func (g *SimpleRESTGetter) Get(url string, v interface{}) error {
	resp, err := http.Get(url)
	//we
	if err != nil {
		return err
	}
	// we are nice, we close the Body
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return RESTError{Code: resp.StatusCode}
	}

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(v)

	return err

}

type RateLimitedRESTGetter struct {
	getter *SimpleRESTGetter
	window time.Duration
	tokens chan bool
}

func NewRateLimitedRESTGetter(limit uint, window time.Duration) *RateLimitedRESTGetter {
	return &RateLimitedRESTGetter{
		getter: NewSimpleRESTGetter(),
		window: window,
		tokens: make(chan bool, limit),
	}
}

func (g *RateLimitedRESTGetter) Get(url string, v interface{}) error {
	//place a token
	g.tokens <- true
	defer func() {
		go func() {
			time.Sleep(g.window)
			<-g.tokens
		}()
	}()

	return g.getter.Get(url, v)

}

type RESTStaticData struct {
	TeamIDs           []string
	SummonerNames     []string
	SummonerIDs       []string
	ChampionIDs       []string
	GameIDs           []string
	RegionCode        string
	Key               APIKey
	ResponseByRequest map[string][]byte
}

type RESTStaticGetter struct {
	data RESTStaticData
}

func NewRESTStaticGetter(data []byte) (*RESTStaticGetter, error) {
	r := bytes.NewReader(data)
	dec := json.NewDecoder(r)
	res := &RESTStaticGetter{}
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

func (g *RESTStaticGetter) Get(url string, v interface{}) error {
	_, ok := g.data.ResponseByRequest[url]
	if ok == false {
		return RESTError{Code: 404}
	}
	r := bytes.NewReader(g.data.ResponseByRequest[url])
	dec := json.NewDecoder(r)

	return dec.Decode(v)
}

func (g *RESTStaticGetter) ATeamID() string {
	return g.data.TeamIDs[0]
}

func (g *RESTStaticGetter) SeveralTeamIDs() []string {
	return g.data.TeamIDs
}

func (g *RESTStaticGetter) ASummonerName() string {
	return g.data.SummonerNames[0]
}

func (g *RESTStaticGetter) SeveralSummonerNames() []string {
	return g.data.SummonerNames
}

func (g *RESTStaticGetter) ASummonerID() string {
	return g.data.SummonerIDs[0]
}

func (g *RESTStaticGetter) SeveralSummonerIDs() []string {
	return g.data.SummonerIDs
}

func (g *RESTStaticGetter) AChampionID() string {
	return g.data.ChampionIDs[0]
}

func (g *RESTStaticGetter) SeveralChampionIDs() []string {
	return g.data.SummonerIDs
}

func (g *RESTStaticGetter) AGameID() string {
	return g.data.GameIDs[0]
}

func (g *RESTStaticGetter) SeveralGameIDs() []string {
	return g.data.GameIDs
}

func (g *RESTStaticGetter) RegionCode() string {
	return g.data.RegionCode
}

func (g *RESTStaticGetter) Key() APIKey {
	return g.data.Key
}
