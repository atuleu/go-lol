package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/charlywazzap/go-lol"
)

type TestCommand struct {
}

func (x *TestCommand) Execute(args []string) error {
	log.Printf("Initial")

	region, err := lol.NewRegionByCode(options.RegionCode)
	if err != nil {
		return err
	}

	keyStorer, err := lol.NewXdgAPIKeyStorer()
	if err != nil {
		return err
	}
	key, ok := keyStorer.Get()
	if ok == false {
		return fmt.Errorf("No key")
	}
	err = key.Check()
	if err != nil {
		return err
	}

	sapi, err := lol.NewStaticAPIEndpoint(region, key)
	if err != nil {
		return err
	}

	log.Printf("done")
	championIDs := make([]lol.ChampionID, 0, len(args))

	for _, idStr := range args {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return err
		}
		championIDs = append(championIDs, lol.ChampionID(id))
	}

	log.Printf("Getting Realm")
	realm, err := sapi.GetRealm()
	if err != nil {
		return err
	}
	log.Printf("Fetching")
	for i, id := range championIDs {
		champ, err := sapi.GetChampion(id)
		if err != nil {
			log.Printf("%s", err)
			continue
		}
		log.Printf("%d : %s\n", i, champ.Name)

		img, err := realm.GetImage(champ.Image)
		if err != nil {
			log.Printf("%s", err)
			continue
		}
		log.Printf("Image is %s in size", img.Bounds())
	}

	return nil
}

func init() {
	parser.AddCommand("test",
		"some test command. DO NOT USE",
		"DO NOT USE",
		&TestCommand{})
}
