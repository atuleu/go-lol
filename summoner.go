package lol

import (
	"fmt"
	"log"
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

type SummonerAPI struct {
	g      RESTGetter
	region *Region
	key    APIKey
}

func NewSummonerAPI(g RESTGetter, region *Region, key APIKey) *SummonerAPI {
	return &SummonerAPI{
		g:      g,
		region: region,
		key:    key,
	}
}

func (a *SummonerAPI) GetSummonerByName(names []string) ([]Summoner, error) {
	if len(names) > 40 {
		return nil, fmt.Errorf("Cannot checkout more than 40 IDs")
	}

	res := make(map[string]Summoner, 40)

	url := fmt.Sprintf("https://%s/api/lol/%s/v1.4/summoner/by-name/%s?api_key=%s",
		a.region.url,
		a.region.code,
		strings.Join(names, ","),
		a.key)

	log.Printf("url is:%s\n", url)

	err := a.g.Get(url, &res)
	if err != nil {
		return nil, err
	}
	actualRes := make([]Summoner, 0, len(res))
	for _, v := range res {
		actualRes = append(actualRes, v)
	}

	return actualRes, err

}

func (a *SummonerAPI) GetSummonerNames(ids []SummonerID) (map[SummonerID]string, error) {
	if len(ids) > 40 {
		return nil, fmt.Errorf("cannot checkout more than 40 Summoner name at once")
	}
	if len(ids) == 0 {
		return nil, fmt.Errorf("You need to provide an ID")
	}

	res := make(map[string]string, len(ids))

	idsStr := fmt.Sprintf("%d", ids[0])
	for i, id := range ids {
		if i == 0 {
			continue
		}
		idsStr = fmt.Sprintf("%s,%d", idsStr, id)
	}

	url := fmt.Sprintf("https://%s/api/lol/%s/v1.4/summoner/%s/name?api_key=%s",
		a.region.url,
		a.region.code,
		idsStr,
		a.key)

	err := a.g.Get(url, &res)
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
