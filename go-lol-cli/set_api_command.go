package main

import (
	"fmt"

	"github.com/charlywazzap/go-lol"
)

type SetAPIKeyCommand struct{}

func (x *SetAPIKeyCommand) Execute(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("'set-api-key need one argument: the api-key to set")
	}

	s, err := lol.NewXdgAPIKeyStorer()
	if err != nil {
		return err
	}

	err = s.Store(lol.APIKey(args[0]))
	if err != nil {
		return err
	}

	return nil
}

func init() {
	parser.AddCommand("set-api-key",
		"Sets the Riot API key to use.",
		"It will install the API key in a safe location, so you don't have to remind it afterwards",
		&SetAPIKeyCommand{})

}
