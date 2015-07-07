package main

import "fmt"

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
		replays := i.manager.Replays()

		if len(replays[i.region.Code()]) == 0 {
			return fmt.Errorf("No replay available for platform %s", i.region.Code())
		}

		x.GameID = uint64(replays[i.region.Code()][0].MetaData.GameKey.ID)
	}

	//gid := lol.GameID(x.GameID)

	return fmt.Errorf("Not yet implemented")
}

func init() {
	parser.AddCommand("replay",
		"Launch replay",
		"Launches a replay",
		&ReplayCommand{})
}
