package xlol

import . "gopkg.in/check.v1"

type ExpandedFormatSuite struct {
}

var _ = Suite(&ExpandedFormatSuite{})

func (s *ExpandedFormatSuite) TestFilenameRegexp(c *C) {
	valid := []string{"chunk.0006.bin", "keyframe.9999.bin"}

	l := &ExpandedReplayFormatter{}

	for _, txt := range valid {
		c.Check(l.checkFileName(txt), IsNil)
	}

}
