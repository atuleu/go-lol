package xlol

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"sort"

	"github.com/charlywazzap/go-lol"
)

// A Chunk represent data modification of a Replay between frames
type Chunk struct {
	ChunkInfo
	KeyFrame KeyFrameID
	data     []byte
}

func (c Chunk) isAssociated() bool {
	return c.KeyFrame > 0
}

// A ChunkList is slice of Chunk that implements sort.Interface
type ChunkList []Chunk

// Len returns the length of the ChunkList
func (l ChunkList) Len() int {
	return len(l)
}

// Less returns true if element in i should be sorted before element
// in j in the ChunkList
func (l ChunkList) Less(i, j int) bool {
	return l[i].ID < l[j].ID
}

// Swap swaps the element placed in i and j in the ChunkList
func (l ChunkList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

// A KeyFrame Represent a Replay state at a given point in time.
type KeyFrame struct {
	KeyFrameInfo
	Chunks []ChunkID
	data   []byte
}

// A KeyFrameList is slice of KeyFrame that implements sort.Interface
type KeyFrameList []KeyFrame

// Len returns the length of the KeyFrameList
func (l KeyFrameList) Len() int {
	return len(l)
}

// Less returns true if element in i should be sorted before element
// in j in the KeyFrameList
func (l KeyFrameList) Less(i, j int) bool {
	return l[i].ID < l[j].ID
}

// Swap swaps the element placed in i and j in the KeyFrameList
func (l KeyFrameList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

// A ReplayDataLoader can load all data of a Replay
type ReplayDataLoader interface {
	HasChunk(ChunkID) bool
	HasKeyFrame(KeyFrameID) bool
	HasEndOfGameStats() bool
	Open() (io.ReadCloser, error)
	OpenChunk(ChunkID) (io.ReadCloser, error)
	OpenKeyFrame(KeyFrameID) (io.ReadCloser, error)
	OpenEndOfGameStats() (io.ReadCloser, error)
}

// A ReplayDataWriter can be used to write data of a replay
type ReplayDataWriter interface {
	Create() (io.WriteCloser, error)
	CreateChunk(ChunkID) (io.WriteCloser, error)
	CreateKeyFrame(KeyFrameID) (io.WriteCloser, error)
	CreateEndOfGameStats() (io.WriteCloser, error)
}

//A ReplayDataFormatter is both a ReplayDataWriter and ReplayDataLoader
type ReplayDataFormatter interface {
	ReplayDataLoader
	ReplayDataWriter
}

// A Replay is a in memory structure that represent all data that is
// needed to spectate a LoL game.
type Replay struct {
	Chunks    []Chunk
	KeyFrames []KeyFrame

	endOfGameStats []byte

	chunksByID   map[ChunkID]int
	keyframeByID map[KeyFrameID]int

	MetaData GameMetadata
	//Version of the spectator API
	Version string
	//Encryption Key used by the data
	EncryptionKey string
	//GameInfo contains the GameInformation
	GameInfo lol.CurrentGameInfo

	// PlayerOfInterest
	PlayerOfInterest lol.SummonerID
}

// NewEmptyReplay creates a new empty replay
func NewEmptyReplay() *Replay {
	return &Replay{
		Chunks:       nil,
		KeyFrames:    nil,
		chunksByID:   make(map[ChunkID]int),
		keyframeByID: make(map[KeyFrameID]int),
		MetaData: GameMetadata{
			PendingAvailableChunkInfo:    []ChunkInfo{},
			PendingAvailableKeyFrameInfo: []KeyFrameInfo{},
		},
	}
}

func (r *Replay) addChunk(c Chunk) {
	if c.ID == 0 {
		return
	}
	if _, ok := r.chunksByID[c.ID]; ok == true {
		return
	}

	r.Chunks = append(r.Chunks, c)
	sort.Sort(ChunkList(r.Chunks))
	r.rebuildChunksMap()
}

func (r *Replay) addKeyFrame(kf KeyFrame) {
	if kf.ID == 0 {
		return
	}
	if _, ok := r.keyframeByID[kf.ID]; ok == true {
		return
	}

	r.KeyFrames = append(r.KeyFrames, kf)
	sort.Sort(KeyFrameList(r.KeyFrames))
	r.rebuildKeyFramesMap()
}

// MergeFromMetaData merge the internal Replay data from GameMetadata
// that can be fetch through the SpectateAPI
func (r *Replay) MergeFromMetaData(gm GameMetadata) {
	newValue := reflect.ValueOf(gm)
	currentValue := reflect.ValueOf(&(r.MetaData))

	for i := 0; i < newValue.NumField(); i++ {
		newField := newValue.Field(i)
		currentField := currentValue.Elem().Field(i)
		switch newField.Kind() {
		case reflect.String:
			if len(newField.String()) != 0 {
				currentField.SetString(newField.String())
			}
		case reflect.Int:
			fallthrough
		case reflect.Int64:
			if newField.Int() != 0 {
				currentField.SetInt(newField.Int())
			}
		}
	}

	if len(gm.GameKey.PlatformID) != 0 && len(r.MetaData.GameKey.PlatformID) == 0 {
		r.MetaData.GameKey = gm.GameKey
	}

	if gm.StartTime.IsZero() == false {
		r.MetaData.StartTime = gm.StartTime
	}

	if gm.CreateTime.IsZero() == false {
		r.MetaData.CreateTime = gm.CreateTime
	}

	for _, ci := range gm.PendingAvailableChunkInfo {
		if cIdx, ok := r.chunksByID[ci.ID]; ok == true {
			if r.Chunks[cIdx].Duration > 0 {
				continue
			}
			r.Chunks[cIdx].ChunkInfo = ci
			continue
		}
		newChunk := Chunk{
			ChunkInfo: ci,
		}

		r.addChunk(newChunk)
	}

	for _, kfi := range gm.PendingAvailableKeyFrameInfo {
		if kfIdx, ok := r.keyframeByID[kfi.ID]; ok == true {
			if r.KeyFrames[kfIdx].NextChunkID > 0 {
				continue
			}
			r.KeyFrames[kfIdx].KeyFrameInfo = kfi
			if idx, ok := r.chunksByID[kfi.NextChunkID]; ok == true {
				r.Chunks[idx].KeyFrame = kfi.ID
			}
			continue
		}
		newKF := KeyFrame{
			KeyFrameInfo: kfi,
		}

		newKF.Chunks = []ChunkID{kfi.NextChunkID}
		if idx, ok := r.chunksByID[kfi.NextChunkID]; ok == true {
			r.Chunks[idx].KeyFrame = kfi.ID
		}

		r.addKeyFrame(newKF)
	}

}

func (r *Replay) appendSortedIfUnique(slice []ChunkID, id ChunkID) []ChunkID {
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

// MergeFromLastChunkInfo merges Replay internal data with the one
// that could be fetch with a LastChunkInfo structure, obtained
// through the SpectateAPI
func (r *Replay) MergeFromLastChunkInfo(ci LastChunkInfo) {
	//avoids ChunkId 0
	if ci.ID == 0 {
		return
	}

	if cIdx, ok := r.chunksByID[ci.ID]; ok == false {
		//we create a new Chunk
		res := Chunk{
			ChunkInfo: ChunkInfo{
				ID:       ci.ID,
				Duration: ci.Duration,
			},
			KeyFrame: ci.AssociatedKeyFrameID,
		}

		if lastIdx, ok := r.chunksByID[ci.ID-1]; ok == true {
			res.ReceivedTime.Time = r.Chunks[lastIdx].ReceivedTime.Add(r.Chunks[lastIdx].Duration.Duration())
		}

		r.addChunk(res)
	} else {
		if r.Chunks[cIdx].Duration == 0 {
			r.Chunks[cIdx].Duration = ci.Duration
		}

		if r.Chunks[cIdx].isAssociated() == false {
			r.Chunks[cIdx].KeyFrame = ci.AssociatedKeyFrameID
		}

	}

	if ci.AssociatedKeyFrameID == 0 {
		return
	}

	kfIdx, ok := r.keyframeByID[ci.AssociatedKeyFrameID]
	if ok == false {
		res := KeyFrame{
			KeyFrameInfo: KeyFrameInfo{
				ID:          ci.AssociatedKeyFrameID,
				NextChunkID: ci.NextChunkID,
			},
			Chunks: []ChunkID{ci.ID},
		}

		res.Chunks = r.appendSortedIfUnique(res.Chunks, ci.NextChunkID)

		if cIdx, ok := r.chunksByID[ci.NextChunkID]; ok == true {
			res.ReceivedTime = r.Chunks[cIdx].ReceivedTime
		}

		r.addKeyFrame(res)
		kfIdx, ok = r.keyframeByID[ci.AssociatedKeyFrameID]
		if ok == false {
			return
		}
	} else {
		if len(r.KeyFrames[kfIdx].Chunks) == 0 {
			r.KeyFrames[kfIdx].Chunks = r.appendSortedIfUnique(r.KeyFrames[kfIdx].Chunks, ci.NextChunkID)
			r.KeyFrames[kfIdx].Chunks = r.appendSortedIfUnique(r.KeyFrames[kfIdx].Chunks, ci.ID)
		}
		if r.KeyFrames[kfIdx].NextChunkID == 0 {
			r.KeyFrames[kfIdx].NextChunkID = ci.NextChunkID
		}
	}

	r.KeyFrames[kfIdx].Chunks = r.appendSortedIfUnique(r.KeyFrames[kfIdx].Chunks, ci.ID)

}

// Consolidate is reconstructing missing internal data (KeyFrame and
// Chunk association) from the internal data we haev so far.
func (r *Replay) Consolidate() {
	if len(r.KeyFrames) == 0 {
		return
	}

	for cIdx, c := range r.Chunks {
		if c.isAssociated() == true {
			continue
		}

		if r.KeyFrames[0].NextChunkID > c.ID {
			// in that case, we could not determine the associated
			// KeyFrame with certainty
			continue
		}
		lastKFID := r.KeyFrames[0].ID
		for _, kf := range r.KeyFrames {
			if kf.NextChunkID > c.ID {
				r.Chunks[cIdx].KeyFrame = lastKFID
				break
			}
			lastKFID = kf.ID
		}
	}
}

// check is checking for integrity of all Replay data, i.e. that all
// data is loaded in memory or is accessible through the given
// loader. Passing a nil loader, will ensure that all required data is
// loaded in memory
func (r *Replay) check(loader ReplayDataLoader) error {
	if len(r.Chunks) == 0 {
		return nil
	}

	// checks that we do not miss a chunk, and all have an associated
	// keyFrame, and the keyframe is available
	noKeyFrameIsFailure := false
	for _, c := range r.Chunks {
		if c.isAssociated() {
			noKeyFrameIsFailure = true
		} else {
			if noKeyFrameIsFailure == true {
				return fmt.Errorf("Missing associated frame for chunk %d", c.ID)
			}
		}

		if len(c.data) == 0 {
			if loader == nil {
				return fmt.Errorf("Data for chunk %d is not loaded, and no loader defined", c.ID)
			}
			if loader.HasChunk(c.ID) == false {
				return fmt.Errorf("Missing data for Chunk %d", c.ID)
			}
		}

		if c.isAssociated() == false {
			continue
		}

		kfIdx, ok := r.keyframeByID[c.KeyFrame]
		if ok == false {
			return fmt.Errorf("Missing metadata for Keyframe %d (associated with chunk %d)", c.KeyFrame, c.ID)
		}

		kf := r.KeyFrames[kfIdx]
		if kf.NextChunkID <= 0 {
			return fmt.Errorf("KeyFrame %d does not know its next chunk id", kfIdx)
		}

		if len(kf.data) == 0 {
			if loader == nil {
				return fmt.Errorf("Data for KeyFrame %d is not loaded, and no loader defined", c.KeyFrame)
			}
			if loader.HasKeyFrame(c.KeyFrame) == false {
				return fmt.Errorf("Missing data for KeyFrame %d", c.KeyFrame)
			}
		}

	}

	if len(r.endOfGameStats) != 0 {
		return nil
	}
	if loader == nil {
		return fmt.Errorf("Missing end of game stat data, and no loader defined")
	}

	if loader.HasEndOfGameStats() == false {
		return fmt.Errorf("Missing end of game stat data")
	}
	return nil
}

func (r *Replay) loadChunk(loader ReplayDataLoader, id ChunkID) error {
	cIdx, ok := r.chunksByID[id]
	if ok == false {
		return fmt.Errorf("Unknown Chunk ID %d", id)
	}

	reader, err := loader.OpenChunk(id)
	if err != nil {
		return fmt.Errorf("Could not open Chunk %d: %s", id, err)
	}
	defer reader.Close()
	r.Chunks[cIdx].data, err = ioutil.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("Could not read Chunk %d data:  %s", id, err)
	}
	return nil
}

func (r *Replay) loadKeyFrame(loader ReplayDataLoader, id KeyFrameID) error {
	kfIdx, ok := r.keyframeByID[id]
	if ok == false {
		return fmt.Errorf("Unknown KeyFrame ID %d (Replay.check() should have reported it)", id)
	}

	reader, err := loader.OpenKeyFrame(id)
	if err != nil {
		return fmt.Errorf("Could not open KeyFrame %d: %s", id, err)
	}
	defer reader.Close()
	r.KeyFrames[kfIdx].data, err = ioutil.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("Could not read KeyFrame %d data:  %s", id, err)
	}

	return nil
}

// LoadData is loading in memory all binary data of Replay (KeyFrame,
// Chunk and EndOfGameStats) through a ReplayDataLoader
func (r *Replay) LoadData(loader ReplayDataLoader) error {
	if err := r.check(loader); err != nil {
		return err
	}
	for _, c := range r.Chunks {
		if err := r.loadChunk(loader, c.ID); err != nil {
			return err
		}

		if c.isAssociated() == false {
			continue
		}

		if err := r.loadKeyFrame(loader, c.KeyFrame); err != nil {
			return err
		}
	}

	reader, err := loader.OpenEndOfGameStats()
	if err != nil {
		return fmt.Errorf("Could not open End Of Game Stat: %s", err)
	}
	defer reader.Close()
	r.endOfGameStats, err = ioutil.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("Could not read end of game stat data: %s", err)
	}
	return nil
}

func (r *Replay) saveChunk(writer ReplayDataWriter, c Chunk) error {
	w, err := writer.CreateChunk(c.ID)
	if err != nil {
		return fmt.Errorf("Could not create Chunk %d: %s", c.ID, err)
	}
	defer w.Close()
	_, err = io.Copy(w, bytes.NewBuffer(c.data))
	if err != nil {
		return fmt.Errorf("Could not write Chunk %d data: %s", c.ID, err)
	}
	return nil
}

func (r *Replay) saveKeyFrame(writer ReplayDataWriter, kf KeyFrame) error {
	w, err := writer.CreateKeyFrame(kf.ID)
	if err != nil {
		return fmt.Errorf("Could not create KeyFrame %d: %s", kf.ID, err)
	}
	defer w.Close()
	_, err = io.Copy(w, bytes.NewBuffer(kf.data))
	if err != nil {
		return fmt.Errorf("Could not write KeyFrame %d data: %s", kf.ID, err)
	}
	return nil
}

// SaveData is saving all binary data of a Replay through a ReplayDataWriter
func (r *Replay) SaveData(writer ReplayDataWriter) error {
	if err := r.check(nil); err != nil {
		return err
	}

	for _, c := range r.Chunks {
		if err := r.saveChunk(writer, c); err != nil {
			return err
		}

		if c.isAssociated() == false {
			continue
		}

		kfIdx, ok := r.keyframeByID[c.KeyFrame]
		if ok == false {
			return fmt.Errorf("Internal consistency error, Replay.check should hae reported an error")
		}

		if err := r.saveKeyFrame(writer, r.KeyFrames[kfIdx]); err != nil {
			return err
		}
	}

	w, err := writer.CreateEndOfGameStats()
	if err != nil {
		return fmt.Errorf("Could not create end of game stat data: %s", err)
	}
	defer w.Close()
	_, err = io.Copy(w, bytes.NewBuffer(r.endOfGameStats))
	if err != nil {
		return fmt.Errorf("Could not write end of game stat data: %s", err)
	}
	return nil
}

func (r *Replay) rebuildChunksMap() {
	r.chunksByID = make(map[ChunkID]int, len(r.Chunks))
	for idx, c := range r.Chunks {
		r.chunksByID[c.ID] = idx
	}
}

func (r *Replay) rebuildKeyFramesMap() {
	r.keyframeByID = make(map[KeyFrameID]int, len(r.KeyFrames))
	for idx, kf := range r.KeyFrames {
		r.keyframeByID[kf.ID] = idx
	}
}

// LoadReplay is loading the Replay data (without loading binary data
// like KeyFrame, Chunk and EndOfGameStats) from a ReplayDataLoader
func LoadReplay(loader ReplayDataLoader) (*Replay, error) {
	if loader == nil {
		return nil, fmt.Errorf("Empty data loader")
	}
	r, err := loader.Open()
	if err != nil {
		return nil, err
	}
	defer r.Close()
	dec := json.NewDecoder(r)
	res := &Replay{}
	err = dec.Decode(res)
	if err != nil {
		return nil, err
	}
	res.rebuildChunksMap()
	res.rebuildKeyFramesMap()

	err = res.check(loader)
	if err != nil {
		return nil, fmt.Errorf("Incomplete replay: %s", err)
	}
	return res, nil
}

// LoadReplayWithData is loading all of a Replay data from a
// ReplayDataLoader
func LoadReplayWithData(loader ReplayDataLoader) (*Replay, error) {
	res, err := LoadReplay(loader)
	if err != nil {
		return nil, err
	}
	return res, res.LoadData(loader)
}

// Save is writing all Replay data (without binary data like KeyFrame,
// Chunk and EndOfGameStats) through a ReplayDataWriter
func (r *Replay) Save(writer ReplayDataWriter) error {
	if err := r.check(nil); err != nil {
		return fmt.Errorf("Could not save replay: %s", err)
	}
	return r.unsafeSave(writer)
}

// Save, but allows incomplete data
func (r *Replay) unsafeSave(writer ReplayDataWriter) error {
	w, err := writer.Create()
	if err != nil {
		return err
	}
	defer w.Close()
	enc := json.NewEncoder(w)
	return enc.Encode(r)
}

// SaveWithData is writing all of Replay data through a
// ReplayDataWriter
func (r *Replay) SaveWithData(writer ReplayDataWriter) error {
	if err := r.Save(writer); err != nil {
		return err
	}
	return r.SaveData(writer)
}

// ChunkByID returns a chunk from its ChunkID
func (r *Replay) ChunkByID(id ChunkID) (*Chunk, bool) {
	if cidx, ok := r.chunksByID[id]; ok == true {
		return &(r.Chunks[cidx]), true
	}
	return nil, false
}

// KeyFrameByID returns a KeyFrame from its ChunkID
func (r *Replay) KeyFrameByID(id KeyFrameID) (*KeyFrame, bool) {
	if kfidx, ok := r.keyframeByID[id]; ok == true {
		return &(r.KeyFrames[kfidx]), true
	}
	return nil, false
}

// AddGameInfo adds game information for a
func (r *Replay) AddGameInfo(info lol.CurrentGameInfo) {
	r.GameInfo = info
	// test if we miss champion name
}

// HighlightSummoner highlights a summonner name
func (r *Replay) HighlightSummoner(id lol.SummonerID) error {
	ok := false
	for _, part := range r.GameInfo.Participants {
		if part.ID == id {
			ok = true
			break
		}
	}

	if ok == false {
		return fmt.Errorf("Summoner (ID:%d) is not in the game", id)
	}

	r.PlayerOfInterest = id
	return nil
}
