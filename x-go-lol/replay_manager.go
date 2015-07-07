package xlol

// A ReplqyManager stores and retrieve replays
type ReplayManager interface {
	Store(*Replay) error
	Get(*lol.Region, lol.GameID) (*Replay, error)
	Replays() map[string]*Replay
}

// NewLocalManager creates a new LocalManager, who data will
// be stored in basedir
func NewLocalManager(basedir string) (*LocalManager, error) {
	res := &LocalManager{}
	var err error
	res.datadir, err = newReplaysDataDir(basedir)
	if err != nil {
		return nil, err
	}
	return res, nil
}
