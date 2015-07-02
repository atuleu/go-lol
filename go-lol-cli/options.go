package main

import "github.com/jessevdk/go-flags"

type Options struct {
	RegionCode string `long:"region" short:"r" description:"region to use for looking up summoners" default:"euw"`
}

var options = &Options{}

var parser = flags.NewParser(options, flags.Default)
