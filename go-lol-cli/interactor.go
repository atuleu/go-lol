package main

import (
	"fmt"

	"github.com/charlywazzap/go-lol"
	"github.com/charlywazzap/go-lol/x-go-lol"
)

type Interactor struct {
	region  *lol.Region
	storer  lol.APIKeyStorer
	key     lol.APIKey
	manager *xlol.XdgReplayManager
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

	res.api, err = lol.NewAPIEndpoint(res.region, res.key)
	if err != nil {
		return nil, err
	}

	res.manager, err = xlol.NewXdgReplayManager()

	if err != nil {
		return nil, err
	}

	return res, nil
}
