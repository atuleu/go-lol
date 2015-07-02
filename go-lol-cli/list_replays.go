package main

import "fmt"

type ListReplaysCommand struct {
	RegionCode string `long:"region" short:"r" description:"region to use for looking up summoners" default:"euw"`
}

func (x *ListReplaysCommand) Execute(args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("list-replays does not take an Arguments")
	}

	i, err := NewInteractor(options)
	if err != nil {
		return err
	}

	byCode, err := i.manager.AvailableReplay()
	if err != nil {
		return err
	}

	gdata := byCode[x.RegionCode]

	if len(gdata) == 0 {
		return fmt.Errorf("No replay for region %s", x.RegionCode)
	}

	fmt.Printf("There are %d replay available for %s:\n", len(gdata), x.RegionCode)
	for i, gm := range gdata {
		fmt.Printf("  * %d : Game %d started at %s\n", i+1, gm.GameKey.ID, gm.CreateTime)
	}

	return nil

}

func init() {
	parser.AddCommand("list-replays",
		"List replays give the list of recorded replays for the given region",
		"List replays give the list of recorded replays for the given region",
		&ListReplaysCommand{})

}
