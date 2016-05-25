package main

import (
	"fmt"
	"log"

	"github.com/charlywazzap/go-lol"
	"github.com/jessevdk/go-flags"
)

// foo
type Options struct {
	//Damn foo
	Key string `short:"k" long:"key" description:"key to save locally" required:"true"`
}

func Execute() error {
	opts := &Options{}
	parser := flags.NewParser(opts, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		return fmt.Errorf("Flag parsing error")
	}
	k := lol.APIKey(opts.Key)

	keyStorer, err := lol.NewXdgAPIKeyStorer()
	if err != nil {
		return err
	}

	return keyStorer.Store(k)
}

func main() {
	if err := Execute(); err != nil {
		log.Fatalf("Got unhandled error: %s", err)
	}
	log.Printf("Stored key into $XDG_CONFIG directory")
}
