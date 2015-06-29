package lol

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type EpochMillisecond uint64

type SummonerID uint64

func (s EpochMillisecond) ToTime() time.Time {
	secs := int64(s) / 1000
	return time.Unix(secs, int64(s)-secs)
}

type Summoner struct {
	Name         string           `json:"name"`
	Id           SummonerID       `json:"id"`
	RevisionDate EpochMillisecond `json:"revisionDate"`
	Level        uint32           `json:"summonerLevel"`
}

func (a *APIRegionalEndpoint) GetSummonerByName(names []string) ([]Summoner, error) {
	if len(names) > 40 {
		return nil, fmt.Errorf("Cannot checkout more than 40 IDs")
	}
	if len(names) == 0 {
		return nil, fmt.Errorf("You need to provide a Summoner name")
	}

	res := make(map[string]Summoner, len(names))

	err := a.Get(fmt.Sprintf("/v1.4/summoner/by-name/%s", strings.Join(names, ",")),
		nil, &res)
	if err != nil {
		return nil, err
	}
	actualRes := make([]Summoner, 0, len(res))
	for _, v := range res {
		actualRes = append(actualRes, v)
	}

	return actualRes, err

}

func (a *APIRegionalEndpoint) GetSummonerNames(ids []SummonerID) (map[SummonerID]string, error) {
	if len(ids) > 40 {
		return nil, fmt.Errorf("cannot checkout more than 40 Summoner name at once")
	}
	if len(ids) == 0 {
		return nil, fmt.Errorf("You need to provide an ID")
	}

	res := make(map[string]string, len(ids))

	idsStr := make([]string, 0, len(ids))
	for _, id := range ids {
		idsStr = append(idsStr, strconv.FormatInt(int64(id), 10))
	}

	err := a.g.Get(fmt.Sprintf("/v1.4/summoner/%s/name", strings.Join(idsStr, ",")), &res)
	if err != nil {
		return nil, err
	}

	actualRes := make(map[SummonerID]string, len(res))
	for idStr, name := range res {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return nil, err
		}
		actualRes[SummonerID(id)] = name
	}
	return actualRes, err

}
