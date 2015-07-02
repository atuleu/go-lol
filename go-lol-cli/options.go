package main

import "github.com/jessevdk/go-flags"

type Options struct{}

var options = &Options{}

var parser = flags.NewParser(options, flags.Default)
