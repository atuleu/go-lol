package xlol

import lol ".."

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
