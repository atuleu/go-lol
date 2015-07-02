package xlol

import (
	"strings"
	"time"
)

// DurationMs represent the number of milliseconds between two point in time
type DurationMs int64

// Duration casts a DurationMs to time.Duration
func (d DurationMs) Duration() time.Duration {
	return time.Duration(d) * time.Millisecond
}

//LolTime is a time.Time that (Un)Marshal -izes itself according to
//Lol date format
type LolTime struct {
	time.Time
}

const (
	lolTimeFormat string = "Jan 2, 2006 3:04:05 PM"
)

func (t LolTime) format() string {
	return t.Time.Format(lolTimeFormat)
}

// MarshalText marshalize the time to the lol format
func (t LolTime) MarshalText() ([]byte, error) {
	return []byte(t.format()), nil
}

// MarshalJSON marshalize the time to the lol format
func (t LolTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.format() + `"`), nil
}

// UnmarshalText unmarshalize the time according to the lol format
func (t *LolTime) UnmarshalText(text []byte) error {
	var err error
	t.Time, err = time.Parse(lolTimeFormat, string(text))
	return err
}

// UnmarshalJSON unmarshalize the time according to the lol format
func (t *LolTime) UnmarshalJSON(text []byte) error {
	var err error
	t.Time, err = time.Parse(lolTimeFormat, strings.Trim(string(text), `"`))
	return err
}
