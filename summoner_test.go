package lol

import (
	"encoding/json"
	"strconv"
	"strings"

	. "gopkg.in/check.v1"
)

type SummonerSuite struct{}

var _ = Suite(&SummonerSuite{})

func (s *SummonerSuite) TestGetSummonerByName(c *C) {
	tooLargeSummonerNames := []string{}
	for len(tooLargeSummonerNames) <= 40 {
		tooLargeSummonerNames = append(tooLargeSummonerNames, getter.SeveralSummonerNames()...)
	}

	summoners, err := api.GetSummonerByName(tooLargeSummonerNames)
	c.Check(len(summoners), Equals, 0)
	c.Check(err, ErrorMatches, "Cannot checkout more than 40 IDs, .* requested")

	summoners, err = api.GetSummonerByName(nil)
	c.Check(len(summoners), Equals, 0)
	c.Check(err, ErrorMatches, "You need to provide at least one Summoner name")

	summoners, err = api.GetSummonerByName([]string{getter.ASummonerName()})
	c.Check(err, IsNil)
	if c.Check(len(summoners), Not(Equals), 0) == true {
		jsonData := string(getter.LastJSONData())
		//particularity, we get here a map from name

		mapped := make(map[string]Summoner)
		for _, sum := range summoners {
			mapped[strings.ToLower(sum.Name)] = sum
		}
		reEncoded, err := json.Marshal(mapped)
		c.Assert(err, IsNil)
		c.Check(string(reEncoded), Equals, jsonData)
	}
}

func (s *SummonerSuite) TestGetSummonerName(c *C) {
	tooManyIds := make([]SummonerID, 0, 40)
	for len(tooManyIds) <= 40 {
		for _, id := range getter.SeveralSummonerIDs() {
			idI, _ := strconv.ParseInt(id, 10, 64)
			tooManyIds = append(tooManyIds, SummonerID(idI))
		}
	}
	names, err := api.GetSummonerNames(tooManyIds)
	c.Check(len(names), Equals, 0)
	c.Check(err, ErrorMatches, "Cannot checkout more than 40 Summoner names, got .*")

	names, err = api.GetSummonerNames(nil)
	c.Check(len(names), Equals, 0)
	c.Check(err, ErrorMatches, "Need at least one Summoner ID")
	manyIds := tooManyIds[0:1]
	names, err = api.GetSummonerNames(manyIds)
	c.Check(err, IsNil)
	if c.Check(len(names), Equals, len(manyIds)) {
		jsonData := string(getter.LastJSONData())
		mapped := make(map[string]string)
		mapped[getter.SeveralSummonerIDs()[0]] = names[manyIds[0]]
		reEncoded, err := json.Marshal(mapped)
		c.Assert(err, IsNil)
		c.Check(string(reEncoded), Equals, jsonData)
	}

}
