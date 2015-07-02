package lol

import "fmt"

// A RegionID is used to uniquely identifies a Region
type RegionID uint

const (
	// BR represents Brazil region
	BR RegionID = iota
	// EUNE represents European Union North - East region
	EUNE
	// EUW represents European Union West region
	EUW
	// KR represents Korea region
	KR
	// LAN represents Latin America North region
	LAN
	// LAS represents Latin America South  region
	LAS
	// NA represents North America region
	NA
	// OCE represents Oceania region
	OCE
	// TR represents Turkish region
	TR
	// RU represents Russia region
	RU
	// PBE represents the Public Beta Environment region
	PBE
	// GLOBAL represents global static data access point
	GLOBAL
)

// A Region defines a set of severs where people can play together
type Region struct {
	platformID   string
	code         string
	url          string
	spectatorURL string
}

var regionByID = map[RegionID]*Region{
	BR: &Region{
		platformID:   "BR1",
		code:         "br",
		url:          "br.api.pvp.net",
		spectatorURL: "spectator.br.lol.riotgames.com",
	},
	EUNE: &Region{
		platformID:   "EUN1",
		code:         "eune",
		url:          "eune.api.pvp.net",
		spectatorURL: "spectator.eu.lol.riotgames.com:8088",
	},
	EUW: &Region{
		platformID:   "EUW1",
		code:         "euw",
		url:          "euw.api.pvp.net",
		spectatorURL: "spectator.euw1.lol.riotgames.com",
	},
	KR: &Region{
		platformID:   "KR1",
		code:         "kr",
		url:          "kr.api.pvp.net",
		spectatorURL: "spectator.kr.lol.riotgames.com",
	},
	LAN: &Region{
		platformID:   "LA1",
		code:         "lan",
		url:          "lan.api.pvp.net",
		spectatorURL: "spectator.la1.lol.riotgames.com",
	},
	LAS: &Region{
		platformID:   "LA2",
		code:         "las",
		url:          "las.api.pvp.net",
		spectatorURL: "spectator.la2.lol.riotgames.com",
	},
	NA: &Region{
		platformID:   "NA1",
		code:         "na",
		url:          "na.api.pvp.net",
		spectatorURL: "spectator.na.lol.riotgames.com",
	},
	OCE: &Region{
		platformID:   "OC1",
		code:         "oce",
		url:          "oce.api.pvp.net",
		spectatorURL: "spectator.oc1.lol.riotgames.com",
	},
	TR: &Region{
		platformID:   "TR1",
		code:         "tr",
		url:          "tr.api.pvp.net",
		spectatorURL: "spectator.tr.lol.riotgames.com",
	},
	RU: &Region{
		platformID:   "RU1",
		code:         "ru",
		url:          "ru.api.pvp.net",
		spectatorURL: "spectator.ru.lol.riotgames.com",
	},
	PBE: &Region{
		platformID:   "PBE1",
		code:         "pbe",
		url:          "pbe.api.pvp.net",
		spectatorURL: "spectator.pbe1.lol.riotgames.com:8088",
	},
	GLOBAL: &Region{
		platformID:   "",
		code:         "global",
		url:          "global.api.pvp.net",
		spectatorURL: "",
	},
}

var regionByCode map[string]*Region

// NewRegion returns a region identified by its RegionID
func NewRegion(id RegionID) (*Region, error) {
	r, ok := regionByID[id]
	if ok == false {
		return nil, fmt.Errorf("unknown RegionID %d", id)
	}

	return r, nil
}

// NewRegionByCode returns a region identified by its code (i.e. "euw"
// for EUW)
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

// Code returns the code used to identifies a Region
func (r *Region) Code() string {
	return r.code
}

// PlatformID returns the ID used by observer mode to identifies the
// Region
func (r *Region) PlatformID() string {
	return r.platformID
}

// SpectatorURL returns the url that should be used to spectate game
// for the Region
func (r *Region) SpectatorURL() string {
	return r.spectatorURL
}
