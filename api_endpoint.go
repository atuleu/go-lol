package lol

import (
	"fmt"
	"time"
)

type APIRegionalEndpoint struct {
	g      RESTGetter
	region *Region
	key    APIKey
}

func NewAPIRegionalEndpoint(region *Region, key APIKey) *APIRegionalEndpoint {
	return &APIRegionalEndpoint{
		g:      NewRateLimitedRESTGetter(10, 10*time.Second),
		region: region,
		key:    key,
	}
}

func (a *APIRegionalEndpoint) FormatUrl(url string, options map[string]string) string {
	res := fmt.Sprintf("https://%s/api/lol/%s%s?api_key=%s", a.region.url, a.region.code, url, a.key)
	for k, v := range options {
		res = fmt.Sprintf("%s&%s=%s", res, k, v)
	}
	return res
}

func (a *APIRegionalEndpoint) Get(url string, options map[string]string, v interface{}) error {
	fullUrl := a.FormatUrl(url, options)
	err := a.g.Get(fullUrl, v)
	if err != nil {
		return fmt.Errorf("Cannot access %s: %s", fullUrl, err)
	}
	return nil
}
