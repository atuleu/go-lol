package xlol

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"sort"

	lol ".."
)

type replaysDataDir struct {
	basedir string
}

// A ChunkID identfies a chunk in a game stream
type ChunkID int

// A KeyFrameID identfies a keyframe in a game stream
type KeyFrameID int

func newReplaysDataDir(basedir string) (*replaysDataDir, error) {
	//check basedir exists, and is user writable
	res := &replaysDataDir{
		basedir: basedir,
	}
	if err := res.ensureUserWritableDirectory(basedir); err != nil {
		return nil, err
	}
	return res, nil
}

func (d *replaysDataDir) ensureUserWritableDirectory(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) == false {
			return err
		}

		err = os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
	} else {
		if info.IsDir() == false {
			return fmt.Errorf("'%s' is not a directory", path)
		}

		if (info.Mode() & 0700) != 0700 {
			return fmt.Errorf("Wrong permission (%s) on base directory '%s'", info.Mode(), path)
		}
	}

	return nil
}

type gameMetadataList []GameMetadata

func (l gameMetadataList) Len() int {
	return len(l)
}

func (l gameMetadataList) Less(i, j int) bool {
	return l[i].GameKey.ID < l[j].GameKey.ID
}

func (l gameMetadataList) Swap(i, j int) {
	l[j], l[i] = l[i], l[j]
}

func (d *replaysDataDir) allFinishedReplays() (map[string][]GameMetadata, error) {
	rx := regexp.MustCompile(`\A[0-9]+\z`)
	res := make(map[string][]GameMetadata)
	for _, r := range lol.AllDynamicRegion() {

		regionPath := path.Join(d.basedir, r.PlatformID())
		info, err := os.Stat(regionPath)
		if err != nil {
			if os.IsNotExist(err) == true {
				continue
			}
			return nil, err
		}
		if info.IsDir() == false {
			return nil, fmt.Errorf("%s is not a directory", regionPath)
		}

		infos, err := ioutil.ReadDir(regionPath)
		resRegion := make([]GameMetadata, 0, len(infos))
		for _, finfo := range infos {
			if finfo.IsDir() == false {
				continue
			}
			if rx.MatchString(finfo.Name()) == false {
				continue
			}

			_, err := os.Stat(path.Join(regionPath, finfo.Name(), endOfGameFile))
			if err != nil {
				if os.IsNotExist(err) == true {
					continue
				}
				return nil, err
			}

			f, err := os.Open(path.Join(regionPath, finfo.Name(), metaDataFile))
			if err != nil {
				return nil, err
			}
			defer f.Close()
			dec := json.NewDecoder(f)
			var gm GameMetadata
			err = dec.Decode(&gm)
			if err != nil {
				return nil, err
			}

			resRegion = append(resRegion, gm)

		}

		//sort resRegion
		sort.Sort(sort.Reverse(gameMetadataList(resRegion)))

		res[r.Code()] = resRegion
	}
	return res, nil
}

const (
	metaDataFile  string = "metadata.json"
	endOfGameFile string = "endofgame.json"
)

type replayDataDir struct {
	parent  *replaysDataDir
	region  *lol.Region
	game    lol.GameID
	gamedir string
}

func newReplayDataDir(parent *replaysDataDir, region *lol.Region, game lol.GameID) (*replayDataDir, error) {
	res := &replayDataDir{
		parent: parent,
		region: region,
		game:   game,
	}

	if err := res.checkRegion(); err != nil {
		return nil, err
	}

	res.gamedir = path.Join(parent.basedir, region.PlatformID(), game.String())

	if err := parent.ensureUserWritableDirectory(res.gamedir); err != nil {
		return nil, err
	}

	return res, nil
}

func (r *replayDataDir) checkRegion() error {
	if r.region == nil {
		return fmt.Errorf("Unitialized region")
	}
	if len(r.region.PlatformID()) == 0 || len(r.region.SpectatorURL()) == 0 {
		return fmt.Errorf("Region %s does not have a spectator mode (static endpoint)", r.region.Code())
	}
	return nil
}

func (r *replayDataDir) metaDataPath() string {
	return path.Join(r.gamedir, metaDataFile)
}

func (r *replayDataDir) endOfGameDataPath() string {
	return path.Join(r.gamedir, endOfGameFile)
}

func (r *replayDataDir) chunkPath(id ChunkID) string {
	return path.Join(r.gamedir, "chunks", fmt.Sprintf("%d", id))
}

func (r *replayDataDir) keyFramePath(id KeyFrameID) string {
	return path.Join(r.gamedir, "keyframes", fmt.Sprintf("%d", id))
}

func (r *replayDataDir) managerDataPath() string {
	return path.Join(r.gamedir, "manager.json")
}
