package main

import (
	lol ".."
)

type LolReplayLauncher interface {
	Launch(region *lol.Region, id lol.GameID, encryptionKey string) error
}
