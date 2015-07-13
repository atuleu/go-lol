package lol

import "fmt"

// StaticAPIEndpoint is used to access Riot REST API for static data.
type StaticAPIEndpoint struct {
	g            RESTGetter
	region       *Region
	key          APIKey
	staticRegion *Region
	realm        Realm
}

// NewStaticAPIEndpoint is creeating a new Static endpoint. You have
// to pass the Dynamic (i.e. EUW, KR NA) region you are interested in
// fecthing data.
func NewStaticAPIEndpoint(region *Region, key APIKey) (*StaticAPIEndpoint, error) {
	if region.IsDynamic() == false {
		return nil, fmt.Errorf("We need a duynamic region for looking up data")
	}
	res := &StaticAPIEndpoint{
		g:      NewSimpleRESTGetter(),
		region: region,
		key:    key,
	}
	res.staticRegion, _ = NewRegion(GLOBAL)

	return res, nil
}

func (a *StaticAPIEndpoint) formatURL(url string, options map[string]string) string {
	res := fmt.Sprintf("https://%s/api/lol/static-data/%s/v1.2%s?api_key=%s",
		a.staticRegion.url, a.region.code, url, a.key)
	for k, v := range options {
		res = res + "&" + k + "=" + v
	}
	return res
}

// get data from that endpoint
func (a *StaticAPIEndpoint) get(url string, options map[string]string, v interface{}) error {
	fullURL := a.formatURL(url, options)
	err := a.g.Get(fullURL, v)
	if err != nil {
		return fmt.Errorf("Cannot access %s: %s", fullURL, err)
	}
	return nil
}
