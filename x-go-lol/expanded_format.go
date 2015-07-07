package xlol

import (
	"fmt"
	"io"
	"os"
	"path"
)

// ExpandedReplayFormatter is a ReplayDataLoader that expand all its
// data to several file on the hardrive. It is meant to
type ExpandedReplayFormatter struct {
	basedir string
}

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
	return res, nil
}

func (l *ExpandedReplayFormatter) chunkPath(id ChunkID) string {
	return path.Join(l.basedir, fmt.Sprintf("chunk.%04d.bin", id))
}

func (l *ExpandedReplayFormatter) keyFramePath(id KeyFrameID) string {
	return path.Join(l.basedir, fmt.Sprintf("keyframe.%04d.bin", id))
}

func (l *ExpandedReplayFormatter) endOfGamePath() string {
	return path.Join(l.basedir, "endOfGameStats.bin")
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
