package xlol

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	lol ".."
)

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

// AssociatedChunkInfo is ChunkINfo, but it keeps the ID of KeyFrame
// it is associated with. Main puprose is to deal with streaming of
// replay
type AssociatedChunkInfo struct {
	ChunkInfo
	KeyFrame KeyFrameID
}

// AssociatedKeyFrameInfo is a KeyFrameInfo, but with the list of
// ChunkID that follow the KeyFrame. MainPurpose is to deal with
// streaming of replay.
type AssociatedKeyFrameInfo struct {
	KeyFrameInfo
	Chunks []ChunkID
}

// ReplayMetadata keeps track of all metadata that is needed to stream
// replay to the LoL client, i.e. generate all LastChunkInfo, and
// GameMetadata dynamically
type ReplayMetadata struct {
	chunks    map[ChunkID]AssociatedChunkInfo
	keyframes map[KeyFrameID]AssociatedKeyFrameInfo

	//Version of the spectator API
	Version string
	//Encryption Key used by the data
	EncryptionKey string
	//ID of the First saved Chunk
	FirstChunk ChunkID
	// ID of the last saved Chunk
	MaxChunk ChunkID
	// ID of the ChunkID of end of Pick & Ban
	EndStartupChunkID ChunkID
	// ID of the CHunk for game time of 0:0
	StartGameChunkID ChunkID
	// ID of the last chunk of the game
	EndGameChunkID ChunkID
	// ID of the last keyframe
	EndGameKeyframeID KeyFrameID
}

func NewReplayMetadata() *ReplayMetadata {
	return &ReplayMetadata{
		chunks:    make(map[ChunkID]AssociatedChunkInfo),
		keyframes: make(map[KeyFrameID]AssociatedKeyFrameInfo),

		FirstChunk:        -1,
		MaxChunk:          -1,
		EndStartupChunkID: -1,
		StartGameChunkID:  -1,
		EndGameChunkID:    -1,
		EndGameKeyframeID: -1,
	}
}

func (d *ReplayMetadata) MergeFromMetaData(gm GameMetadata) {
	for _, c := range gm.PendingAvailableChunkInfo {
		if _, ok := d.chunks[c.ID]; ok == true {
			continue
		}
		d.chunks[c.ID] = AssociatedChunkInfo{ChunkInfo: c}
	}

	for _, kf := range gm.PendingAvailableKeyFrameInfo {
		if _, ok := d.keyframes[kf.ID]; ok == true {
			continue
		}
		res := AssociatedKeyFrameInfo{KeyFrameInfo: kf}
		res.Chunks = []ChunkID{kf.NextChunkID}
		if c, ok := d.chunks[kf.NextChunkID]; ok == true {
			c.KeyFrame = kf.ID
			d.chunks[kf.NextChunkID] = c
		}
		d.keyframes[res.ID] = res
	}

	if gm.EndStartupChunkID > 0 {
		d.EndStartupChunkID = ChunkID(gm.EndStartupChunkID)
	}

	if gm.StartGameChunkID > 0 {
		d.StartGameChunkID = ChunkID(gm.StartGameChunkID)
	}

	if gm.EndGameChunkID > 0 {
		d.EndGameChunkID = ChunkID(gm.EndGameChunkID)
	}

	if gm.EndGameKeyFrameID > 0 {
		d.EndGameKeyframeID = KeyFrameID(gm.EndGameKeyFrameID)
	}
}

func (d *ReplayMetadata) appendSortedIfUnique(slice []ChunkID, id ChunkID) []ChunkID {
	pos := -1
	for i, cid := range slice {
		if cid == id {
			return slice
		}
		if cid < id {
			pos = i + 1
		}
	}

	if pos == -1 {
		return append([]ChunkID{id}, slice...)
	} else if pos == len(slice) {
		return append(slice, id)
	}
	return append(slice[:pos], append([]ChunkID{id}, slice[pos:]...)...)
}

func (d *ReplayMetadata) MergeFromLastChunkInfo(ci LastChunkInfo) {
	if _, ok := d.chunks[ci.ID]; ok == false {
		//we create a new Chunk
		res := AssociatedChunkInfo{
			ChunkInfo: ChunkInfo{
				ID:       ci.ID,
				Duration: ci.Duration,
			},
			KeyFrame: ci.AssociatedKeyFrameID,
		}

		if last, ok := d.chunks[ci.ID-1]; ok == true {
			res.ReceivedTime.Time = last.ReceivedTime.Add(last.Duration.Duration())
		}
	}

	chunk := d.chunks[ci.ID]
	keyframe, ok := d.keyframes[ci.AssociatedKeyFrameID]
	if ok == false {

		res := AssociatedKeyFrameInfo{
			KeyFrameInfo: KeyFrameInfo{
				ID:          ci.AssociatedKeyFrameID,
				NextChunkID: ci.NextChunkID,
			},
			Chunks: []ChunkID{ci.ID},
		}
		if res.NextChunkID == ci.ID {
			res.ReceivedTime = chunk.ReceivedTime
		}

		d.keyframes[ci.AssociatedKeyFrameID] = res
		keyframe = d.keyframes[ci.AssociatedKeyFrameID]
	}

	keyframe.Chunks = d.appendSortedIfUnique(keyframe.Chunks, ci.ID)
	d.keyframes[ci.AssociatedKeyFrameID] = keyframe

}

func (d *ReplayMetadata) Consolidate() {
	//So we go through all the Chunk, and we first determine the first and the last we have
	if len(d.chunks) == 0 {
		return
	}
	d.FirstChunk = ChunkID(int(^uint(0) >> 1)) //Maximal int value
	d.MaxChunk = -(d.FirstChunk - 1)           //Minimal int value

	kfIDs := make([]int, 0, len(d.keyframes))
	for id := range d.keyframes {
		kfIDs = append(kfIDs, int(id))
	}
	sort.Sort(sort.IntSlice(kfIDs))

	for _, c := range d.chunks {
		//computes min and max
		if d.FirstChunk > c.ID {
			d.FirstChunk = c.ID
		}
		if d.MaxChunk < c.ID {
			d.MaxChunk = c.ID
		}

		if c.KeyFrame != 0 {
			continue
		}

		var lastKf KeyFrameID = -1
		for i, kfi := range kfIDs {
			if d.keyframes[KeyFrameID(kfi)].NextChunkID > c.ID {
				if i > 0 {
					c.KeyFrame = lastKf
					d.chunks[c.ID] = c
					kf := d.keyframes[lastKf]
					kf.Chunks = d.appendSortedIfUnique(kf.Chunks, c.ID)
					d.keyframes[lastKf] = kf
				}
				break
			}
			lastKf = KeyFrameID(kfi)
		}
	}
}

func (d *ReplayMetadata) check(ddir *replayDataDir) error {
	// checks that we do not miss a chunk, and all have an associated
	// keyFrame, and the keyframe is available
	noKeyFrameIsFailure := false
	for i := d.FirstChunk; i <= d.MaxChunk; i++ {
		c, ok := d.chunks[ChunkID(i)]
		if ok == false {
			return fmt.Errorf("Missing chunk %d", i)
		}

		_, err := os.Stat(ddir.chunkPath(c.ID))
		if err != nil {
			return err
		}

		if c.KeyFrame > 0 {
			noKeyFrameIsFailure = true
		} else {
			if noKeyFrameIsFailure == true {
				return fmt.Errorf("Missing associated frame for chunk %d", c.ID)
			}
		}

		_, ok = d.keyframes[c.KeyFrame]
		if ok == false {
			return fmt.Errorf("Missing Keyframe %d", c.KeyFrame)
		}
		_, err = os.Stat(ddir.keyFramePath(c.KeyFrame))
		if err != nil {
			return err
		}
	}
	return nil
}

type ReplayMetadataForJSON struct {
	ReplayMetadata
	Chunks    []AssociatedChunkInfo
	KeyFrames []AssociatedKeyFrameInfo
}

type AssociatedChunkInfoList []AssociatedChunkInfo
type AssociatedKeyFrameInfoList []AssociatedKeyFrameInfo

func (l AssociatedChunkInfoList) Len() int {
	return len(l)
}

func (l AssociatedChunkInfoList) Less(i, j int) bool {
	return l[i].ID < l[j].ID
}

func (l AssociatedChunkInfoList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l AssociatedKeyFrameInfoList) Len() int {
	return len(l)
}

func (l AssociatedKeyFrameInfoList) Less(i, j int) bool {
	return l[i].ID < l[j].ID
}

func (l AssociatedKeyFrameInfoList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (d *ReplayMetadata) MarshalJSON() ([]byte, error) {
	temp := ReplayMetadataForJSON{
		ReplayMetadata: *d,
		Chunks:         make([]AssociatedChunkInfo, 0, len(d.chunks)),
		KeyFrames:      make([]AssociatedKeyFrameInfo, 0, len(d.keyframes)),
	}

	for _, c := range d.chunks {
		temp.Chunks = append(temp.Chunks, c)
	}

	for _, kf := range d.keyframes {
		temp.KeyFrames = append(temp.KeyFrames, kf)
	}
	sort.Sort(AssociatedChunkInfoList(temp.Chunks))
	sort.Sort(AssociatedKeyFrameInfoList(temp.KeyFrames))

	return json.Marshal(temp)
}

func (d *ReplayMetadata) UnmarshalJSON(text []byte) error {
	temp := &ReplayMetadataForJSON{}

	err := json.Unmarshal(text, &temp)
	if err != nil {
		return err
	}

	d.EncryptionKey = temp.EncryptionKey
	d.FirstChunk = temp.FirstChunk
	d.MaxChunk = temp.MaxChunk
	d.EndStartupChunkID = temp.EndStartupChunkID
	d.StartGameChunkID = temp.StartGameChunkID
	d.EndGameChunkID = temp.EndGameChunkID
	d.EndGameKeyframeID = temp.EndGameKeyframeID

	for _, c := range temp.Chunks {
		d.chunks[c.ID] = c
	}

	for _, kf := range temp.KeyFrames {
		d.keyframes[kf.ID] = kf
	}

	return nil
}
