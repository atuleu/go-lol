package xlol

// LastChunkInfo represents latest available Chunk informations
type LastChunkInfo struct {
	ID                   ChunkID    `json:"chunkId"`
	AvailableSince       DurationMs `json:"availableSince"`
	NextAvailableChunk   DurationMs `json:"nextAvailableChunk"`
	AssociatedKeyFrameID KeyFrameID `json:"keyFrameId"`
	NextChunkID          ChunkID    `json:"nextChunkId"`
	EndStartupChunkID    ChunkID    `json:"endStartupChunkId"`
	StartGameChunkID     ChunkID    `json:"startGameChunkId"`
	EndGameChunkID       ChunkID    `json:"endGameChunkId"`
	Duration             DurationMs `json:"duration"`
}
