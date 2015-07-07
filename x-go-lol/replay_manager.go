package xlol

import (
	lol ".."
)

// A ReplayManager stores and retrieve replays
type ReplayManager interface {
	Store(*Replay) error
	Get(*lol.Region, lol.GameID) (ReplayDataLoader, error)
	Replays() map[string]*Replay
}

// LocalManager is a stub
type XdgReplayManager struct{}

// NewLocalManager creates a new LocalManager, who data will
// be stored in basedir
func NewXdgReplayManager() (*XdgReplayManager, error) {
	res := &LocalManager{}
	return res, nil
}
