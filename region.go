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
	platformId   string
	code         string
	url          string
	spectatorUrl string
}

var regionByID = map[RegionID]*Region{
	BR: &Region{
		platformId:   "BR1",
		code:         "br",
		url:          "br.api.pvp.net",
		spectatorUrl: "spectator.br.lol.riotgames.com",
	},
	EUNE: &Region{
		platformId:   "EUN1",
		code:         "eune",
		url:          "eune.api.pvp.net",
		spectatorUrl: "spectator.eu.lol.riotgames.com:8088",
	},
	EUW: &Region{
		platformId:   "EUW1",
		code:         "euw",
		url:          "euw.api.pvp.net",
		spectatorUrl: "spectator.euw1.lol.riotgames.com",
	},
	KR: &Region{
		platformId:   "KR1",
		code:         "kr",
		url:          "kr.api.pvp.net",
		spectatorUrl: "spectator.kr.lol.riotgames.com",
	},
	LAN: &Region{
		platformId:   "LA1",
		code:         "lan",
		url:          "lan.api.pvp.net",
		spectatorUrl: "spectator.la1.lol.riotgames.com",
	},
	LAS: &Region{
		platformId:   "LA2",
		code:         "las",
		url:          "las.api.pvp.net",
		spectatorUrl: "spectator.la2.lol.riotgames.com",
	},
	NA: &Region{
		platformId:   "NA1",
		code:         "na",
		url:          "na.api.pvp.net",
		spectatorUrl: "spectator.na.lol.riotgames.com",
	},
	OCE: &Region{
		platformId:   "OC1",
		code:         "oce",
		url:          "oce.api.pvp.net",
		spectatorUrl: "spectator.oc1.lol.riotgames.com",
	},
	TR: &Region{
		platformId:   "TR1",
		code:         "tr",
		url:          "tr.api.pvp.net",
		spectatorUrl: "spectator.tr.lol.riotgames.com",
	},
	RU: &Region{
		platformId:   "RU1",
		code:         "ru",
		url:          "ru.api.pvp.net",
		spectatorUrl: "spectator.ru.lol.riotgames.com",
	},
	PBE: &Region{
		platformId:   "PBE1",
		code:         "pbe",
		url:          "pbe.api.pvp.net",
		spectatorUrl: "spectator.pbe1.lol.riotgames.com:8088",
	},
	GLOBAL: &Region{
		platformId:   "",
		code:         "global",
		url:          "global.api.pvp.net",
		spectatorUrl: "",
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
