package xlol

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"
)

// ExpandedReplayFormatter is a ReplayDataLoader that expand all its
// data to several file on the hardrive. It is meant to
type ExpandedReplayFormatter struct {
	basedir string
}

const (
	expandedFormatVersion string = "1~dev1"
)

var binRx = regexp.MustCompile(`\A([a-z]+)\.([0-9]{4})\.bin\z`)

// NewExpandedReplayFormatter returns a ReplayDataLoader that will
// save and load data from directory basedir
func NewExpandedReplayFormatter(basedir string) (*ExpandedReplayFormatter, error) {
	res := &ExpandedReplayFormatter{
		basedir: basedir,
	}

	//ensure the basedir is an user writable directory

	info, err := os.Stat(basedir)
	if err != nil {
		if os.IsNotExist(err) == false {
			return nil, fmt.Errorf("Could not check directory '%s': %s", basedir, err)
		}

		err = os.MkdirAll(basedir, 0755)
		if err != nil {
			return nil, fmt.Errorf("Could not create directory %s: %s", basedir, err)
		}

	} else {
		if info.IsDir() == false {
			return nil, fmt.Errorf("%s is not a directory", basedir)
		}
		if info.Mode()&0700 != 0700 {
			return nil, fmt.Errorf("%s is not user writable", basedir)
		}
	}
	//check for version of the format
	err = res.check()
	if err != nil {
		return nil, err
	}
	return res, nil
}

const (
	version    = "version"
	eogStat    = "endOfGameStats.bin"
	replayData = "replayData.json"
	chunk      = "chunk"
	keyframe   = "keyframe"
)

type invalidFileError struct {
	dir      string
	fileName string
}

func (e invalidFileError) Error() string {
	return fmt.Sprintf("expanded go-lol format: Invalid file %s in %s", e.fileName, e.dir)
}

func (l *ExpandedReplayFormatter) checkCompatible(version string) error {
	if version != expandedFormatVersion {
		return fmt.Errorf("Mismatched Replay data format version %s, expected %s",
			version, expandedFormatVersion)
	}
	return nil
}

func (l *ExpandedReplayFormatter) check() error {
	//check if the directory is empty
	finfos, err := ioutil.ReadDir(l.basedir)
	if err != nil {
		return err
	}

	if len(finfos) == 0 {
		// This is a new replay, we just create version information
		// and report no error
		f, err := os.Create(l.versionPath())
		if err != nil {
			return err
		}
		defer f.Close()
		fmt.Fprintf(f, "%s\n", expandedFormatVersion)
		return nil
	}

	// check for version match
	f, err := os.Open(l.versionPath())
	if err != nil {
		return fmt.Errorf("Could not check replay data format version: %s", err)
	}
	defer f.Close()
	var localVersion string
	_, err = fmt.Fscanf(f, "%s\n", &localVersion)
	if err != nil {
		return fmt.Errorf("Could not extract version from %s: %s", l.versionPath(), err)
	}

	err = l.checkCompatible(localVersion)
	if err != nil {
		return err
	}

	for _, inf := range finfos {
		m := binRx.FindStringSubmatch(inf.Name())
		if len(m) != 0 {
			if m[1] != chunk && m[1] != keyframe {
				return invalidFileError{fileName: inf.Name(), dir: l.basedir}
			}
		}

		switch inf.Name() {
		case version:
			fallthrough
		case eogStat:
			fallthrough
		case replayData:
		default:
			return invalidFileError{fileName: inf.Name(), dir: l.basedir}
		}
	}

	return nil
}

func (l *ExpandedReplayFormatter) versionPath() string {
	return path.Join(l.basedir, version)
}

func (l *ExpandedReplayFormatter) chunkPath(id ChunkID) string {
	return path.Join(l.basedir, fmt.Sprintf("%s.%04d.bin", chunk, id))
}

func (l *ExpandedReplayFormatter) keyFramePath(id KeyFrameID) string {
	return path.Join(l.basedir, fmt.Sprintf("%s.%04d.bin", keyframe, id))
}

func (l *ExpandedReplayFormatter) endOfGamePath() string {
	return path.Join(l.basedir, eogStat)
}

func (l *ExpandedReplayFormatter) dataPath() string {
	return path.Join(l.basedir, replayData)
}

func (l *ExpandedReplayFormatter) fileExists(path string) bool {
	info, err := os.Stat(path)
	if err == nil {
		return info.Size() > 0
	}
	if os.IsNotExist(err) == true {
		return false
	}
	panic(err)
}

// HasChunk returns true if data is available for a given Chunk
func (l *ExpandedReplayFormatter) HasChunk(id ChunkID) bool {
	return l.fileExists(l.chunkPath(id))
}

// HasKeyFrame returns true if data is available for a given KeyFrame
func (l *ExpandedReplayFormatter) HasKeyFrame(id KeyFrameID) bool {
	return l.fileExists(l.keyFramePath(id))
}

// HasEndOfGameStats returns true if data is available for End of Game
// statistics
func (l *ExpandedReplayFormatter) HasEndOfGameStats() bool {
	return l.fileExists(l.endOfGamePath())
}

// OpenChunk returns a io.ReadCloser for reading data for a given
// Chunk
func (l *ExpandedReplayFormatter) OpenChunk(id ChunkID) (io.ReadCloser, error) {
	return os.Open(l.chunkPath(id))
}

// OpenKeyFrame returns a io.ReadCloser for reading data for a given
// KeyFrame
func (l *ExpandedReplayFormatter) OpenKeyFrame(id KeyFrameID) (io.ReadCloser, error) {
	return os.Open(l.keyFramePath(id))
}

// OpenEndOfGameStats returns a io.ReadCloser for reading data for the end
// of game statistics
func (l *ExpandedReplayFormatter) OpenEndOfGameStats() (io.ReadCloser, error) {
	return os.Open(l.endOfGamePath())
}

// CreateChunk returns a truncated io.WriteCloser to write data for a
// given Chunk
func (l *ExpandedReplayFormatter) CreateChunk(id ChunkID) (io.WriteCloser, error) {
	return os.Create(l.chunkPath(id))
}

// CreateKeyFrame returns a truncated io.WriteCloser to write data for a
// given KeyFrame
func (l *ExpandedReplayFormatter) CreateKeyFrame(id KeyFrameID) (io.WriteCloser, error) {
	return os.Create(l.keyFramePath(id))
}

// CreateEndOfGameStats returns a truncated io.WriteCloser to write
// data for the end of game statistics
func (l *ExpandedReplayFormatter) CreateEndOfGameStats() (io.WriteCloser, error) {
	return os.Create(l.endOfGamePath())
}

// Create returns a truncatde io.WriteCloser to write replay data
func (l *ExpandedReplayFormatter) Create() (io.WriteCloser, error) {
	return os.Create(l.dataPath())
}

// Open returns a io.ReadCloser for reading replay data
func (l *ExpandedReplayFormatter) Open() (io.ReadCloser, error) {
	return os.Open(l.dataPath())
}
