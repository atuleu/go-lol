package xlol

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"sort"

	lol ".."
	"launchpad.net/go-xdg"
)

// A ReplayManager stores and retrieve replays
type ReplayManager interface {
	Store(*Replay) error
	Get(*lol.Region, lol.GameID) (ReplayDataLoader, error)
	Create(*lol.Region, lol.GameID) (ReplayDataWriter, error)
	Replays() map[string]*Replay
}

// A XdgReplayManager is a ReplayManager that stores its data within
// the XDG_CACHE_HOME directory
type XdgReplayManager struct {
	basedir string
}

const (
	xdgReplayManagerFormatVersion = "1~dev1"
)

// NewXdgReplayManager creates a ReplayManager that stores its data in
// XDG_CACHE_HOME
func NewXdgReplayManager() (*XdgReplayManager, error) {
	res := &XdgReplayManager{}

	versionPath := path.Join(xdg.Cache.Home(), "go-lol", "replays", "version")
	res.basedir = path.Dir(versionPath)
	err := os.MkdirAll(res.basedir, 0755)
	if err != nil {
		return nil, fmt.Errorf("Could not create cache directory %s: %s",
			res.basedir, err)
	}

	_, err = os.Stat(versionPath)
	if err != nil {
		if os.IsNotExist(err) == false {
			return nil, err
		}
		f, err := os.Create(versionPath)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		fmt.Fprintf(f, "%s\n", xdgReplayManagerFormatVersion)
		return res, nil
	}

	f, err := os.Open(versionPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var localVersion string
	_, err = fmt.Fscanf(f, "%s\n", &localVersion)
	if err != nil {
		return nil, fmt.Errorf("Could not read cache directory version: %s", err)
	}

	err = res.checkCompatible(localVersion)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (m *XdgReplayManager) checkCompatible(version string) error {
	if version != xdgReplayManagerFormatVersion {
		return fmt.Errorf("Invalid version format %s, expected %s", version, expandedFormatVersion)
	}
	return nil
}

// Store saves in the cache directory all the Replay data.
func (m *XdgReplayManager) Store(r *Replay) error {
	path := path.Join(m.basedir, r.MetaData.GameKey.PlatformID,
		fmt.Sprintf("%d", r.MetaData.GameKey.ID))

	formatter, err := NewExpandedReplayFormatter(path)
	if err != nil {
		return err
	}

	return r.SaveWithData(formatter)
}

// Get returns a ReplayDataLoader for a given lol.Region and
// lol.GameID that is fully stored in the Manager. It will return an
// error if the replay data is missing or incomplete for LoL client to
// spectate.
func (m *XdgReplayManager) Get(region *lol.Region, id lol.GameID) (ReplayDataLoader, error) {
	basepath := path.Join(m.basedir, region.PlatformID(), fmt.Sprintf("%s", id))

	_, err := os.Stat(basepath)
	if err != nil {
		return nil, err
	}

	formatter, err := NewExpandedReplayFormatter(basepath)
	if err != nil {
		return nil, err
	}

	if formatter.HasEndOfGameStats() == false {
		return nil, fmt.Errorf("Requested game %s/%d is not finished, missing EndOfGameStats",
			region.PlatformID(), id)
	}

	return formatter, nil

}

// Create returns a ReplayDataWriter for a replay identifieud by its
// lol.Region and lol.GameID. It will fails if a replay already exists
// for that Game.
func (m *XdgReplayManager) Create(region *lol.Region, id lol.GameID) (ReplayDataWriter, error) {
	basepath := path.Join(m.basedir, region.PlatformID(), fmt.Sprintf("%s", id))

	formatter, err := NewExpandedReplayFormatter(basepath)
	if err != nil {
		return nil, err
	}

	f, err := formatter.Open()
	if err == nil {
		f.Close()
	}
	if err == nil || os.IsNotExist(err) == false {
		return nil, fmt.Errorf("Cannot create a replay for game %s/%d: some replay data already exists",
			region.PlatformID(), id)
	}

	return formatter, nil
}

var gameIDRx = regexp.MustCompile(`\A[0-9]+\z`)

type replayList []*Replay

func (l replayList) Len() int {
	return len(l)
}

func (l replayList) Less(i, j int) bool {
	if l[i].MetaData.GameKey.PlatformID == l[j].MetaData.GameKey.PlatformID {
		return l[i].MetaData.GameKey.ID < l[j].MetaData.GameKey.ID
	}
	return l[i].MetaData.GameKey.PlatformID < l[j].MetaData.GameKey.PlatformID
}

func (l replayList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (m *XdgReplayManager) replaysOfRegion(platformID string) []*Replay {
	platformBasePath := path.Join(m.basedir, platformID)
	finfos, err := ioutil.ReadDir(platformBasePath)
	if err != nil {
		return nil
	}

	res := make([]*Replay, 0, len(finfos))
	for _, inf := range finfos {
		if gameIDRx.MatchString(inf.Name()) == false {
			continue
		}
		if inf.IsDir() == false {
			continue
		}
		replayBasePath := path.Join(platformBasePath, inf.Name())
		formatter, err := NewExpandedReplayFormatter(replayBasePath)
		if err != nil {
			continue
		}
		if formatter.HasEndOfGameStats() == false {
			continue
		}

		r, err := LoadReplay(formatter)
		if err != nil {
			continue
		}
		res = append(res, r)
	}

	sort.Sort(sort.Reverse(replayList(res)))

	return res
}

// Replays return all Replay stored in the XdgReplayManager
func (m *XdgReplayManager) Replays() map[string][]*Replay {
	res := make(map[string][]*Replay)
	for _, r := range lol.AllDynamicRegion() {
		pinfo, err := os.Stat(path.Join(m.basedir, r.PlatformID()))
		if err != nil {
			continue
		}
		if pinfo.IsDir() == false {
			continue
		}

		res[r.Code()] = m.replaysOfRegion(r.PlatformID())
	}

	return res
}
