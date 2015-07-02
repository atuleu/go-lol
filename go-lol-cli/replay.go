package main

import (
	"fmt"
	"log"
	"net/http"

	lol ".."
)

type ReplayCommand struct {
	GameID uint64 `long:"game-id" short:"g" description:"ID of the game to replay, if none the most recent on the region is replayed"`
}

func (x *ReplayCommand) Execute(args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("replay command does not take any arguments")
	}

	i, err := NewInteractor(options)
	if err != nil {
		return err
	}

	if x.GameID == 0 {
		gdata, err := i.manager.AvailableReplay()
		if err != nil {
			return err
		}
		if len(gdata[i.region.Code()]) == 0 {
			return fmt.Errorf("No replay available for platform %s", i.region.Code())
		}
		x.GameID = uint64(gdata[i.region.Code()][0].GameKey.ID)
	}

	gid := lol.GameID(x.GameID)

	handler, eKey, err := i.manager.GetHandler(i.region, gid)
	if err != nil {
		return fmt.Errorf("Could not get HTTP server handler: %s", err)
	}

	http.Handle("/observer-mode/rest/consumer/", handler)

	log.Printf("Started replay server on localhost:4000")
	go log.Fatal(http.ListenAndServe(":4080", nil))

	launcher, err := NewLolReplayLauncher("")
	if err != nil {
		fmt.Printf("Cannot launch automatically the replay, please launch it manually using this encryption key : %s", eKey)
		return nil
	}
	launcher.Launch(i.region, gid, eKey)
	return nil
}

func init() {
	parser.AddCommand("replay",
		"Launch replay",
		"Launches a replay",
		&ReplayCommand{})
}
