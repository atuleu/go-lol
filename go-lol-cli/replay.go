package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	lol ".."
	xlol "../x-go-lol"
)

type ReplayCommand struct {
	GameID  uint64 `long:"game-id" short:"g" description:"ID of the game to replay, if none the most recent on the region is replayed"`
	Address string `long:"address" short:"a" description:"Address of the replay server" default:":8088"`
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
		log.Printf("Loading all replays")
		replays := i.manager.Replays()
		log.Printf("done")
		if len(replays[i.region.Code()]) == 0 {
			return fmt.Errorf("No replay available for platform %s", i.region.Code())
		}

		x.GameID = uint64(replays[i.region.Code()][0].MetaData.GameKey.ID)
	}

	gid := lol.GameID(x.GameID)

	replayLoader, err := i.manager.Get(i.region, gid)
	if err != nil {
		return err
	}

	server, err := xlol.NewReplayServer(replayLoader)
	if err != nil {
		return err
	}

	sigchan := make(chan os.Signal)
	errchan := make(chan error)

	signal.Notify(sigchan, os.Interrupt)

	go func() {
		log.Printf("Starting replay serve on %s", x.Address)
		err := server.ListenAndServe(x.Address)
		log.Printf("Server finished")
		errchan <- err
	}()

	//TODO : starts the launcher

	<-sigchan
	log.Printf("Stopping the server after catching SIGINT")
	err = server.Close()
	if err != nil {
		return err
	}

	return <-errchan
}

func init() {
	parser.AddCommand("replay",
		"Launch replay",
		"Launches a replay",
		&ReplayCommand{})
}
