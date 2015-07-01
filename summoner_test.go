package lol

import . "gopkg.in/check.v1"

type SummonerSuite struct {
	api *APIRegionalEndpoint
}

var _ = Suite(&SummonerSuite{})

func (s *SummonerSuite) SetUpSuite(c *C) {
	s.api = &APIRegionalEndpoint{
		g:      getter,
		region: regionTest,
		key:    getter.Key(),
	}
}

func (s *SummonerSuite) TestGetSummonerByName(c *C) {
}
