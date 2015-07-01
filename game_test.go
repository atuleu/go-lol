package lol

import (
	"strconv"

	. "gopkg.in/check.v1"
)

type GameSuite struct{}

var _ = Suite(&GameSuite{})

func (s *SummonerSuite) TestGetSummonerRecentGames(c *C) {

	id, err := strconv.ParseInt(getter.ASummonerID(), 10, 64)
	c.Assert(err, IsNil)

	games, err := api.GetSummonerRecentGames(SummonerID(id))
	c.Assert(err, IsNil)
	getter.LastJSONData()
	c.Check(len(games), Not(Equals), 0)
	// API is not really clear and does not ship always all data
	// so we should not compare

	// reEncoded, err := json.Marshal(games)
	// c.Assert(err, IsNil)

	// var reEncodedIndented, jsonDataIndented bytes.Buffer
	// c.Assert(json.Indent(&reEncodedIndented, reEncoded, "", "  "), IsNil)
	// c.Assert(json.Indent(&jsonDataIndented, []byte(jsonData), "", "  "), IsNil)

	// c.Check(reEncodedIndented.String(), Equals, jsonDataIndented.String())

}
