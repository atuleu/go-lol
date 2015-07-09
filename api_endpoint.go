package lol

import (
	"fmt"
	"time"
)

// An APIEndpoint represents an endpoint that can fetch dynamic data
// about League of Legend
type APIEndpoint struct {
	g      RESTGetter
	region *Region
	key    APIKey
}

// NewAPIEndpoint creates a new APIEndpoint from a Region and an
// APIKey
func NewAPIEndpoint(region *Region, key APIKey) (*APIEndpoint, error) {
	if region.IsDynamic() == false {
		return nil, fmt.Errorf("APIEndpoint only works with dynamic regions")
	}
	return &APIEndpoint{
		g:      NewRateLimitedRESTGetter(10, 10*time.Second),
		region: region,
		key:    key,
	}, nil
}

// formats an url for that endpoint
func (a *APIEndpoint) formatURL(url string, options map[string]string) string {
	res := fmt.Sprintf("https://%s/api/lol/%s%s?api_key=%s", a.region.url, a.region.code, url, a.key)
	for k, v := range options {
		res = fmt.Sprintf("%s&%s=%s", res, k, v)
	}
	return res
}

// get data from that endpoint
func (a *APIEndpoint) get(url string, options map[string]string, v interface{}) error {
	fullURL := a.formatURL(url, options)
	err := a.g.Get(fullURL, v)
	if err != nil {
		return fmt.Errorf("Cannot access %s: %s", fullURL, err)
	}
	return nil
}
