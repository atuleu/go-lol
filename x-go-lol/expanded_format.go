package xlol

import (
	"fmt"
	"os"
)

// ExpandedReplayDataLoader is a ReplayDataLoader that expand all its data to several file on the hardrive. It is meant to
type ExpandedReplayDataLoader struct {
	basedir string
}

func NewExpandedReplayDataLoader(basedir string) (*ExpandedReplayDataLoader, error) {
	res := ExpandedReplayDataLoader{
		basedir: basedir,
	}

	//ensure the basedir is an user writable directory

	info, err := os.Stat(basedir)
	if err != nil {
		if os.IsNotExist(err) == false {
			return nil, fmt.Errorf("Could not check directory '%s': %s", err)
		}

		err = os.MkdirAll(basedir, 0755)
		if err != nil {
			return nil, fmt.Errorf("Could not create directory %s: %s", basedir, err)
		}

	}

}
