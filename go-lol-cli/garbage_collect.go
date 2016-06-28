package main

import (
	"fmt"
	"log"
	"time"

	"github.com/charlywazzap/go-lol"
)

type GarbageCollectCommand struct {
	Limit     int    `long:"limit" short:"l" description:"Keep at most this number of replay, negative number disable the limit. ) erase all" default:"-1"`
	OlderThan string `long:"older-than" short:"o" description:"Removes replays that are older than this date, default is 840h (i.e. 5 weeks)" default:"840h"`
}

func (x *GarbageCollectCommand) Execute(args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("garbage-collect does not take any arguments")
	}

	maxAge, err := time.ParseDuration(x.OlderThan)
	if err != nil {
		return fmt.Errorf("Invalid duration %s: %s", x.OlderThan, err)
	}

	i, err := NewInteractor(options)
	if err != nil {
		return err
	}

	allReplays := i.manager.Replays()

	thresholdDate := time.Now().Add(-maxAge)
	for _, region := range lol.AllDynamicRegion() {
		replays := allReplays[region.Code()]
		for idx, r := range replays {
			shouldDelete := false
			if x.Limit >= 0 && idx >= x.Limit {
				log.Printf("Deleting %s/%d: outside the limit of %d", i.region.Code(), r.MetaData.GameKey.ID, x.Limit)
				shouldDelete = true
			}
			if r.MetaData.CreateTime.Before(thresholdDate) {
				log.Printf("Deleting %s/%d: older than %s", i.region.Code(), r.MetaData.GameKey.ID, thresholdDate)
				shouldDelete = true
			}

			if shouldDelete == false {
				log.Printf("Keeping %s/%d", i.region.Code(), r.MetaData.GameKey.ID)
				continue
			}

			if err := i.manager.Delete(i.region, r.MetaData.GameKey.ID); err != nil {
				return err
			}
		}
	}

	return i.manager.CleanUp()
}

func init() {
	parser.AddCommand("garbage-collect",
		"Removes old replay from local cache data",
		"It removes replays from local cache (data will be lost forever). It uses two criteria: maximum number of replay to keep, and the maximal age of the replays to keep",
		&GarbageCollectCommand{})
}
