package xlol

import (
	"strings"
	"time"

	lol ".."
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
func (t LolTime) UnmarshalText(text []byte) error {
	var err error
	t.Time, err = time.Parse(lolTimeFormat, string(text))
	return err
}

// UnmarshalJSON unmarshalize the time according to the lol format
func (t LolTime) UnmarshalJSON(text []byte) error {
	var err error
	t.Time, err = time.Parse(lolTimeFormat, strings.Trim(string(text), `"`))
	return err
}

// GameMetadata represents a game metadata for downloading / replaying it
type GameMetadata struct {
	GameKey struct {
		ID         lol.GameID `json:"gameId"`
		PlatformID string     `json:"platformId"`
	} `json:"gameKey"`

	GameServerAddress string     `json:"gameServerAddress"`
	Port              int        `json:"port"`
	EncryptionKey     string     `json:"encryptionKey"`
	ChunkTimeInterval DurationMs `json:"chunkTimeInterval"`
	StartTime         LolTime    `json:"startTime"`
	LastChunkID       int        `json:"lastChunkId"`
	LastKeyFrameID    int        `json:"lastKeyFrameId"`
	EndStartupChunkID int        `json:"endStartupChunkId"`
	DelayTime         DurationMs `json:"delayTime"`

	PendingAvailableChunkInfo []struct {
		ID           int        `json:"id"`
		Duration     DurationMs `json:"duration"`
		ReceivedTime LolTime    `json:"receivedTime"`
	} `json:"pendingAvailableChunkInfo"`

	PendingAvailableKeyFrameInfo []struct {
		ID           int     `json:"id"`
		ReceivedTime LolTime `json:"receivedTime"`
		NextChunkID  int     `json:"nextChunkId"`
	} `json:"pendingAvailableKeyFrameInfo"`

	KeyFrameInterval          DurationMs
	DecodedEncryptionKey      string     `json:"decodedEncryptionKey"`
	StartGameChunkID          int        `json:"startGameChunkId"`
	ClientAddedLag            DurationMs `json:"clientAddedLag"`
	ClientBackFetchingEnabled bool       `json:"clientBackFetchingEnabled"`
	ClientBackFetchingFreq    int        `json:"clientBackFetchingFreq"`
	InterestScore             int        `json:"interestScore"`
	FeaturedGame              bool       `json:"featuredGame"`
	CreateTime                LolTime    `json:"createTime"`
	EndGameChunkID            int        `json:"endGameChunkId"`
	EndGameKeyFrameID         int        `json:"endGameKeyFrameId"`
}
