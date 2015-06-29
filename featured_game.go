package lol

import "fmt"

type FeaturedGame struct {
	Games          []CurrentGame `json:"gameList"`
	RefrehInterval int64         `json:"clientRefreshInterval"`
}

func (a *APIRegionalEndpoint) GetFeaturedGames() (*FeaturedGame, error) {
	res := &FeaturedGame{}

	err := a.g.Get(fmt.Sprintf("https://%s/observer-mode/rest/featured?api_key=%s",
		a.region.url,
		a.key), res)
	if err != nil {
		return nil, fmt.Errorf("Could not fetch featured game on %s: %s", a.region.code, err)
	}
	return res, nil

}
