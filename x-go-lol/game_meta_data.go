package xlol

import "github.com/charlywazzap/go-lol"

// A ChunkID identifies a Chunk in a game stream
type ChunkID int

// A KeyFrameID identifies a KeyFrame in a game stream
type KeyFrameID int

// ChunkInfo are information about chunk
type ChunkInfo struct {
	ID           ChunkID    `json:"id"`
	Duration     DurationMs `json:"duration"`
	ReceivedTime LolTime    `json:"receivedTime"`
}

// KeyFrameInfo are information about KeyFrame
type KeyFrameInfo struct {
	ID           KeyFrameID `json:"id"`
	ReceivedTime LolTime    `json:"receivedTime"`
	NextChunkID  ChunkID    `json:"nextChunkId"`
}

// GameMetadata represents a game metadata for downloading / replaying
// it. This is the data send and receive by the unofficial LoL
// spectator API
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

	PendingAvailableChunkInfo    []ChunkInfo    `json:"pendingAvailableChunkInfo"`
	PendingAvailableKeyFrameInfo []KeyFrameInfo `json:"pendingAvailableKeyFrameInfo"`

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
