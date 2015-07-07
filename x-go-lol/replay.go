package xlol

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"sort"
)

// A Chunk represent data modification of a Replay between frames
type Chunk struct {
	ChunkInfo
	KeyFrame KeyFrameID
	data     []byte
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
}

// NewEmptyReplay creates a new empty replay
func NewEmptyReplay() *Replay {
	return &Replay{
		Chunks:       nil,
		KeyFrames:    nil,
		chunksByID:   make(map[ChunkID]int),
		keyframeByID: make(map[KeyFrameID]int),
	}
}

func (d *Replay) addChunk(c Chunk) {

	if _, ok := d.chunksByID[c.ID]; ok == true {
		return
	}

	d.Chunks = append(d.Chunks, c)
	sort.Sort(ChunkList(d.Chunks))
	d.chunksByID = make(map[ChunkID]int)
	for i, cc := range d.Chunks {
		d.chunksByID[cc.ID] = i
	}
}

func (d *Replay) addKeyFrame(kf KeyFrame) {
	if _, ok := d.keyframeByID[kf.ID]; ok == true {
		return
	}

	d.KeyFrames = append(d.KeyFrames, kf)
	sort.Sort(KeyFrameList(d.KeyFrames))
	d.keyframeByID = make(map[KeyFrameID]int)
	for i, kkf := range d.KeyFrames {
		d.keyframeByID[kkf.ID] = i
	}
}

// MergeFromMetaData merge the internal Replay data from GameMetadata
// that can be fetch through the SpectateAPI
func (d *Replay) MergeFromMetaData(gm GameMetadata) {
	newValue := reflect.ValueOf(gm)
	currentValue := reflect.ValueOf(d.MetaData)

	for i := 0; i < newValue.NumField(); i++ {
		newField := newValue.Field(i)
		currentField := currentValue.Field(i)
		switch newField.Kind() {
		case reflect.String:
			if len(newValue.String()) != 0 {
				currentField.SetString(newField.String())
			}
		case reflect.Int:
			fallthrough
		case reflect.Int64:
			if newValue.Int() > 0 {
				currentField.SetInt(newValue.Int())
			}
		}
	}

	if gm.StartTime.IsZero() == false {
		d.MetaData.StartTime = gm.StartTime
	}

	if gm.CreateTime.IsZero() == false {
		d.MetaData.CreateTime = gm.CreateTime
	}

	for _, ci := range gm.PendingAvailableChunkInfo {
		if _, ok := d.chunksByID[ci.ID]; ok == true {
			continue
		}
		newChunk := Chunk{
			ChunkInfo: ci,
		}

		d.addChunk(newChunk)
	}

	for _, kfi := range gm.PendingAvailableKeyFrameInfo {
		if _, ok := d.keyframeByID[kfi.ID]; ok == true {
			continue
		}
		newKF := KeyFrame{
			KeyFrameInfo: kfi,
		}

		newKF.Chunks = []ChunkID{kfi.NextChunkID}
		if idx, ok := d.chunksByID[kfi.NextChunkID]; ok == true {
			d.Chunks[idx].KeyFrame = kfi.ID
		}

		d.addKeyFrame(newKF)
	}

}

func (d *Replay) appendSortedIfUnique(slice []ChunkID, id ChunkID) []ChunkID {
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
func (d *Replay) MergeFromLastChunkInfo(ci LastChunkInfo) {
	if _, ok := d.chunksByID[ci.ID]; ok == false {
		//we create a new Chunk
		res := Chunk{
			ChunkInfo: ChunkInfo{
				ID:       ci.ID,
				Duration: ci.Duration,
			},
			KeyFrame: ci.AssociatedKeyFrameID,
		}

		if lastIdx, ok := d.chunksByID[ci.ID-1]; ok == true {
			res.ReceivedTime.Time = d.Chunks[lastIdx].ReceivedTime.Add(d.Chunks[lastIdx].Duration.Duration())
		}

		d.addChunk(res)
	}

	kfIdx, ok := d.keyframeByID[ci.AssociatedKeyFrameID]
	if ok == false {

		res := KeyFrame{
			KeyFrameInfo: KeyFrameInfo{
				ID:          ci.AssociatedKeyFrameID,
				NextChunkID: ci.NextChunkID,
			},
			Chunks: []ChunkID{ci.ID},
		}

		if cIdx, ok := d.chunksByID[ci.NextChunkID]; ok == true {
			res.ReceivedTime = d.Chunks[cIdx].ReceivedTime
		}

		d.addKeyFrame(res)
		kfIdx = d.keyframeByID[ci.AssociatedKeyFrameID]
	}

	d.KeyFrames[kfIdx].Chunks = d.appendSortedIfUnique(d.KeyFrames[kfIdx].Chunks, ci.ID)
}

// Consolidate is reconstructing missing internal data (KeyFrame and
// Chunk association) from the internal data we haev so far.
func (d *Replay) Consolidate() {
	if len(d.KeyFrames) == 0 {
		return
	}

	for cIdx, c := range d.Chunks {
		if c.KeyFrame != 0 {
			continue
		}

		if d.KeyFrames[0].NextChunkID > c.ID {
			// in that case, we could not determine the associated
			// KeyFrame with certainty
			continue
		}
		lastKFID := d.KeyFrames[0].ID
		for _, kf := range d.KeyFrames {
			if kf.NextChunkID > c.ID {
				d.Chunks[cIdx].KeyFrame = lastKFID
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
func (d *Replay) check(loader ReplayDataLoader) error {
	if len(d.Chunks) == 0 {
		return nil
	}

	// checks that we do not miss a chunk, and all have an associated
	// keyFrame, and the keyframe is available
	noKeyFrameIsFailure := false
	for _, c := range d.Chunks {
		if c.KeyFrame > 0 {
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

		kfIdx, ok := d.keyframeByID[c.KeyFrame]
		if ok == false {
			return fmt.Errorf("Missing metadata for Keyframe %d (associated with chunk %d)", c.KeyFrame, c.ID)
		}

		if len(d.KeyFrames[kfIdx].data) == 0 {
			if loader == nil {
				return fmt.Errorf("Data for KeyFrame %d is not loaded, and no loader defined", c.KeyFrame)
			}
			if loader.HasKeyFrame(c.KeyFrame) == false {
				return fmt.Errorf("Missing data for KeyFrame %d", c.KeyFrame)
			}
		}

	}

	if len(d.endOfGameStats) != 0 {
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

type replayForJSON struct {
	Replay
}

// UnmarshalJSON unmarshals a JSON version of a Replay. Indeed some
// internal field are not exported to JSON and need to be
// reconstructed after a standard unmarshal operation
func (d *Replay) UnmarshalJSON(text []byte) error {
	tmp := replayForJSON{}
	if err := json.Unmarshal(text, tmp); err != nil {
		return err
	}
	d.Chunks = tmp.Chunks
	d.KeyFrames = tmp.KeyFrames
	d.MetaData = tmp.MetaData
	d.Version = tmp.Version
	d.EncryptionKey = tmp.EncryptionKey

	d.chunksByID = make(map[ChunkID]int)
	for i, c := range d.Chunks {
		d.chunksByID[c.ID] = i
	}
	d.keyframeByID = make(map[KeyFrameID]int)
	for i, kf := range d.KeyFrames {
		d.keyframeByID[kf.ID] = i
	}

	return nil
}

// LoadData is loading in memory all binary data of Replay (KeyFrame,
// Chunk and EndOfGameStats) through a ReplayDataLoader
func (d *Replay) LoadData(loader ReplayDataLoader) error {
	if err := d.check(loader); err != nil {
		return err
	}
	for i, c := range d.Chunks {
		r, err := loader.OpenChunk(c.ID)
		if err != nil {
			return fmt.Errorf("Could not open Chunk %d: %s", c.ID, err)
		}
		defer r.Close()
		d.Chunks[i].data, err = ioutil.ReadAll(r)
		if err != nil {
			return fmt.Errorf("COuld not read Chunk %d data:  %s", c.ID, err)
		}

		kfIdx, ok := d.keyframeByID[c.KeyFrame]
		if ok == false {
			return fmt.Errorf("Internal consistency error, Replay.check should hae reported an error")
		}

		rr, err := loader.OpenKeyFrame(c.KeyFrame)
		if err != nil {
			return fmt.Errorf("Could not open KeyFrame %d: %s", c.KeyFrame, err)
		}
		defer rr.Close()
		d.KeyFrames[kfIdx].data, err = ioutil.ReadAll(rr)
		if err != nil {
			return fmt.Errorf("Could not read KeyFrame %d data:  %s", c.KeyFrame, err)
		}
	}

	r, err := loader.OpenEndOfGameStats()
	if err != nil {
		return fmt.Errorf("Could not open End Of Game Stat: %s", err)
	}
	defer r.Close()
	d.endOfGameStats, err = ioutil.ReadAll(r)
	if err != nil {
		return fmt.Errorf("Could not read end of game stat data: %s", err)
	}
	return nil
}

// SaveData is saving all binary data of a Replay through a ReplayDataWriter
func (d *Replay) SaveData(writer ReplayDataWriter) error {
	if err := d.check(nil); err != nil {
		return err
	}

	for _, c := range d.Chunks {
		w, err := writer.CreateChunk(c.ID)
		if err != nil {
			return fmt.Errorf("Could not create Chunk %d: %s", c.ID, err)
		}
		defer w.Close()
		_, err = io.Copy(w, bytes.NewBuffer(c.data))
		if err != nil {
			return fmt.Errorf("Could not write Chunk %d data: %s", c.ID, err)
		}
		kfIdx, ok := d.keyframeByID[c.KeyFrame]
		if ok == false {
			return fmt.Errorf("Internal consistency error, Replay.check should hae reported an error")
		}

		ww, err := writer.CreateKeyFrame(c.KeyFrame)
		if err != nil {
			return fmt.Errorf("Could not create KeyFrame %d: %s", c.KeyFrame, err)
		}
		defer ww.Close()
		_, err = io.Copy(ww, bytes.NewBuffer(d.KeyFrames[kfIdx].data))
		if err != nil {
			return fmt.Errorf("Could not write Chunk %d data: %s", c.KeyFrame, err)
		}

	}
	w, err := writer.CreateEndOfGameStats()
	if err != nil {
		return fmt.Errorf("Could not create end of game stat data: %s", err)
	}
	defer w.Close()
	_, err = io.Copy(w, bytes.NewBuffer(d.endOfGameStats))
	if err != nil {
		return fmt.Errorf("Could not write end of game stat data: %s", err)
	}
	return nil
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
func (d *Replay) Save(writer ReplayDataWriter) error {
	if err := d.check(nil); err != nil {
		return fmt.Errorf("Could not save replay: %s", err)
	}

	w, err := writer.Create()
	if err != nil {
		return err
	}
	enc := json.NewEncoder(w)
	return enc.Encode(d)
}

// SaveWithData is writing all of Replay data through a
// ReplayDataWriter
func (d *Replay) SaveWithData(writer ReplayDataWriter) error {
	if err := d.Save(writer); err != nil {
		return err
	}
	return d.SaveData(writer)
}
