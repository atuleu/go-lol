package main

import (
	"fmt"
	"path"

	lol ".."
	xlol "../x-go-lol"
	"launchpad.net/go-xdg"
)

type Interactor struct {
	region  *lol.Region
	storer  lol.APIKeyStorer
	key     lol.APIKey
	manager *xlol.LocalManager
	api     *lol.APIEndpoint
}

func NewInteractor(options *Options) (*Interactor, error) {
	res := &Interactor{}
	var err error

	res.region, err = lol.NewRegionByCode(options.RegionCode)
	if err != nil {
		return nil, err
	}

	res.storer, err = lol.NewXdgAPIKeyStorer()
	if err != nil {
		return nil, err
	}
	var ok bool
	res.key, ok = res.storer.Get()
	if ok == false {
		return nil, fmt.Errorf("It seems that there are no API key store, did you use set-api-key? ")
	}

	err = res.key.Check()
	if err != nil {
		return nil, err
	}

	res.api = lol.NewAPIEndpoint(res.region, res.key)

	cachedir, err := xdg.Cache.Ensure("go-lol/versions")
	if err != nil {
		return nil, fmt.Errorf("Could not initialize cachedir %s: %s", cachedir, err)
	}

	cachedir = path.Dir(cachedir)

	res.manager, err = xlol.NewLocalManager(cachedir)

	if err != nil {
		return nil, err
	}

	return res, nil
}
