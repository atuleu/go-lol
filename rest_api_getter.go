package lol

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type RESTGetter interface {
	Get(url string, v interface{}) error
}

type SimpleRESTGetter struct{}

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

	if resp.StatusCode == 429 {
		return fmt.Errorf("Too many request to server")
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Got response %d", resp.StatusCode)
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

type APIEndpoint struct {
	g      *RESTGetter
	region Region
}
