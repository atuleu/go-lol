package main

import "os"

func main() {
	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}
}
