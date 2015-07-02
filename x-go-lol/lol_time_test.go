package xlol

import (
	"testing"
	"time"

	. "gopkg.in/check.v1"
)

type LolTimeSuite struct{}

var _ = Suite(&LolTimeSuite{})

func Test(t *testing.T) { TestingT(t) }

func (s *LolTimeSuite) TestCanBeParsed(c *C) {
	validData := map[string]time.Time{
		"Jul 2, 2015 10:47:51 AM": time.Date(2015, 7, 2, 10, 47, 51, 0, time.UTC),
		"Jul 2, 2015 10:48:21 PM": time.Date(2015, 7, 2, 22, 48, 21, 0, time.UTC),
	}

	for text, value := range validData {
		var t LolTime
		if c.Check(t.UnmarshalText([]byte(text)), IsNil) == true {
			c.Check(t.Time, Equals, value)
		}

		if c.Check(t.UnmarshalJSON([]byte(`"`+text+`"`)), IsNil) == true {
			c.Check(t.Time, Equals, value)
		}

	}
}
