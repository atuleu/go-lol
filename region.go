package lol

import "fmt"

type RegionID uint

const (
	BR RegionID = iota
	EUNE
	EUW
	KR
	LAN
	LAS
	NA
	OCE
	TR
	RU
	PBE
	GLOBAL
)

type Region struct {
	platformId string
	code       string
	url        string
}

var regionByID = map[RegionID]*Region{
	BR: &Region{
		platformId: "BR1",
		code:       "br",
		url:        "br.api.pvp.net",
	},
	EUNE: &Region{
		platformId: "EUN1",
		code:       "eune",
		url:        "eune.api.pvp.net",
	},
	EUW: &Region{
		platformId: "EUW1",
		code:       "euw",
		url:        "euw.api.pvp.net",
	},
	KR: &Region{
		platformId: "KR1",
		code:       "kr",
		url:        "kr.api.pvp.net",
	},
	LAN: &Region{
		platformId: "LA1",
		code:       "lan",
		url:        "lan.api.pvp.net",
	},
	LAS: &Region{
		platformId: "LA2",
		code:       "las",
		url:        "las.api.pvp.net",
	},
	NA: &Region{
		platformId: "NA1",
		code:       "na",
		url:        "na.api.pvp.net",
	},
	OCE: &Region{
		platformId: "OC1",
		code:       "oce",
		url:        "oce.api.pvp.net",
	},
	TR: &Region{
		platformId: "TR1",
		code:       "tr",
		url:        "tr.api.pvp.net",
	},
	RU: &Region{
		platformId: "RU1",
		code:       "ru",
		url:        "ru.api.pvp.net",
	},
	PBE: &Region{
		platformId: "PBE1",
		code:       "pbe",
		url:        "pbe.api.pvp.net",
	},
	GLOBAL: &Region{
		platformId: "",
		code:       "global",
		url:        "global.api.pvp.net",
	},
}

var regionByCode map[string]*Region

func NewRegion(id RegionID) (*Region, error) {
	r, ok := regionByID[id]
	if ok == false {
		return nil, fmt.Errorf("unknown RegionID %d", id)
	}

	return r, nil
}

func NewRegionByCode(code string) (*Region, error) {
	r, ok := regionByCode[code]
	if ok == false {
		return nil, fmt.Errorf("unknown RegionID %s", code)
	}

	return r, nil
}

func init() {
	regionByCode = make(map[string]*Region)
	for _, r := range regionByID {
		regionByCode[r.code] = r
	}
}
