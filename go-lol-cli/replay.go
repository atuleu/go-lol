package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/charlywazzap/go-lol"
	"github.com/charlywazzap/go-lol/x-go-lol"
)

type ReplayCommand struct {
	GameID     uint64 `long:"game-id" short:"g" description:"ID of the game to replay, if none the most recent on the region is replayed"`
	Address    string `long:"address" short:"a" description:"Address of the replay server" default:"localhost:8088"`
	TimeFactor uint   `long:"time-factor" short:"t" description:"Time multiplication factor when streaming game, a too high value may hinder the performance of the client" default:"4"`
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
		replays := i.manager.Replays()
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
	server.TimeDivisor = xlol.DurationMs(x.TimeFactor)

	launcher, err := NewLolReplayLauncher("")
	if err != nil {
		log.Printf(`Could not find launcher replay: %s
You may want to launch manually the LoL client as described at https://developer.riotgames.com/docs/spectating-games, with arguments "spectator %s %s %d %s"`,
			err, x.Address, server.EncryptionKey(), gid, i.region.PlatformID())
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

	finish := make(chan struct{})

	if launcher != nil {
		go func() {
			errLauncher := launcher.Launch(x.Address, i.region, gid, server.EncryptionKey())
			if errLauncher != nil {
				log.Printf("Client error: %s", errLauncher)
			}
			close(finish)
		}()
	} else {
		log.Printf("Will stream until SIGINT")
	}

	select {
	case <-sigchan:
		log.Printf("Stopping the server after catching SIGINT")
	case <-finish:
		log.Printf("Stopping the server after client exit")
	}

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
