package main

import (
	"fmt"

	"github.com/charlywazzap/go-lol/x-go-lol"
)

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

	printer, err := xlol.NewReplayPrinter(i.region, i.key)
	if err != nil {
		return err
	}

	byCode := i.manager.Replays()

	replays := byCode[x.RegionCode]

	if len(replays) == 0 {
		return fmt.Errorf("No replay for region %s", x.RegionCode)
	}

	fmt.Printf("There are %d replay available for %s:\n", len(replays), x.RegionCode)
	for _, r := range replays {
		printer.Display(r)
		fmt.Printf("\n")
	}

	return nil

}

func init() {
	parser.AddCommand("list-replays",
		"List replays give the list of recorded replays for the given region",
		"List replays give the list of recorded replays for the given region",
		&ListReplaysCommand{})

}
