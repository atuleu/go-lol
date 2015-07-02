package lol

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// EpochMillisecond represents a point in time by the number of
// milliseconds since EPOCH
type EpochMillisecond uint64

// SummonerID uniquely identifies a Summoner
type SummonerID uint64

// Time converts EpochMillisecond to time.Time
func (s EpochMillisecond) Time() time.Time {
	secs := int64(s) / 1000
	return time.Unix(secs, int64(s)-secs)
}

// A Summoner is a representation of a player on LoL servers
type Summoner struct {
	ID            SummonerID       `json:"id"`
	Name          string           `json:"name"`
	ProfileIconID int              `json:"profileIconId"`
	Level         uint32           `json:"summonerLevel"`
	RevisionDate  EpochMillisecond `json:"revisionDate"`
}

// GetSummonerByName returns Summoner data identified by their names
func (a *APIEndpoint) GetSummonerByName(names []string) ([]Summoner, error) {
	if len(names) > 40 {
		return nil, fmt.Errorf("Cannot checkout more than 40 IDs, %d requested", len(names))
	}
	if len(names) == 0 {
		return nil, fmt.Errorf("You need to provide at least one Summoner name")
	}

	res := make(map[string]Summoner, len(names))

	err := a.get(fmt.Sprintf("/v1.4/summoner/by-name/%s", strings.Join(names, ",")),
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

// GetSummonerNames returns the name of Summoner identified by their
// IDs
func (a *APIEndpoint) GetSummonerNames(ids []SummonerID) (map[SummonerID]string, error) {
	if len(ids) > 40 {
		return nil, fmt.Errorf("Cannot checkout more than 40 Summoner names, got %d", len(ids))
	}
	if len(ids) == 0 {
		return nil, fmt.Errorf("Need at least one Summoner ID")
	}

	res := make(map[string]string, len(ids))

	idsStr := make([]string, 0, len(ids))
	for _, id := range ids {
		idsStr = append(idsStr, strconv.FormatInt(int64(id), 10))
	}

	err := a.get(fmt.Sprintf("/v1.4/summoner/%s/name", strings.Join(idsStr, ",")), nil, &res)
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
