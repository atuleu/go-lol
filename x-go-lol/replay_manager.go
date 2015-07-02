package xlol

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"sort"
	"time"

	lol ".."
)

// A ReplayDownloader is able to Download (spectate) a game currently
// played on a LoL server.
type ReplayDownloader interface {
	Download(region *lol.Region, id lol.GameID) error
}

// A ReplayHandler is able to server over HTTP a game, as if it would
// be on a LoL server.
type ReplayHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

// A ReplayManager is both a ReplayHandler and a ReplayDownloader
type ReplayManager interface {
	ReplayDownloader
	ReplayHandler
}

// A LocalManager is a ReplayManager that store its data in a
// perticular location on the FileSystem
type LocalManager struct {
	datadir *replaysDataDir
}

// NewLocalManager creates a new LocalManager, who data will
// be stored in basedir
func NewLocalManager(basedir string) (*LocalManager, error) {
	res := &LocalManager{}
	var err error
	res.datadir, err = newReplaysDataDir(basedir)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *LocalManager) saveJSON(path string, v interface{}) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	_, err = io.Copy(f, bytes.NewBuffer(data))
	return err
	// enc := json.NewEncoder(f)
	// return enc.Encode(v)
}

func (m *LocalManager) downloadBinary(api *SpectateAPI, fn SpectateFunction, id int, filepath string) error {
	err := os.MkdirAll(path.Dir(filepath), 0755)
	if err != nil {
		return err
	}
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	err = api.ReadAll(fn, id, f)
	if err == nil {
		log.Printf("Saved %s", filepath)
		return nil
	}

	if rerr, ok := err.(lol.RESTError); ok == true {
		if rerr.Code == 404 {
			log.Printf("Skipped %s", filepath)
			return nil
		}
	}
	return err
}

type associatedChunkInfo struct {
	ChunkInfo
	KeyFrame KeyFrameID
}

type associatedKeyFrameInfo struct {
	KeyFrameInfo
	Chunks []ChunkID
}

type managerLocalData struct {
	chunks    map[ChunkID]associatedChunkInfo
	keyframes map[KeyFrameID]associatedKeyFrameInfo

	EncryptionKey     string
	FirstChunk        ChunkID
	MaxChunk          ChunkID
	EndStartupChunkID ChunkID
	StartGameChunkID  ChunkID
	EndGameChunkID    ChunkID
	EndGameKeyframeID KeyFrameID
}

func newManagerLocalData() *managerLocalData {
	return &managerLocalData{
		chunks:    make(map[ChunkID]associatedChunkInfo),
		keyframes: make(map[KeyFrameID]associatedKeyFrameInfo),

		FirstChunk:        -1,
		MaxChunk:          -1,
		EndStartupChunkID: -1,
		StartGameChunkID:  -1,
		EndGameChunkID:    -1,
		EndGameKeyframeID: -1,
	}
}

func (d *managerLocalData) MergeFromMetaData(gm GameMetadata) {
	for _, c := range gm.PendingAvailableChunkInfo {
		if _, ok := d.chunks[c.ID]; ok == true {
			continue
		}
		d.chunks[c.ID] = associatedChunkInfo{ChunkInfo: c}
	}

	for _, kf := range gm.PendingAvailableKeyFrameInfo {
		if _, ok := d.keyframes[kf.ID]; ok == true {
			continue
		}
		res := associatedKeyFrameInfo{KeyFrameInfo: kf}
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

func (d *managerLocalData) appendSortedIfUnique(slice []ChunkID, id ChunkID) []ChunkID {
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

func (d *managerLocalData) MergeFromLastChunkInfo(ci LastChunkInfo) {
	if _, ok := d.chunks[ci.ID]; ok == false {
		//we create a new Chunk
		res := associatedChunkInfo{
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

		res := associatedKeyFrameInfo{
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

func (d *managerLocalData) Consolidate() {
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

type managerLocalDataForJSON struct {
	managerLocalData
	Chunks    []associatedChunkInfo
	KeyFrames []associatedKeyFrameInfo
}

type associatedChunkInfoList []associatedChunkInfo
type associatedKeyFrameInfoList []associatedKeyFrameInfo

func (l associatedChunkInfoList) Len() int {
	return len(l)
}

func (l associatedChunkInfoList) Less(i, j int) bool {
	return l[i].ID < l[j].ID
}

func (l associatedChunkInfoList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l associatedKeyFrameInfoList) Len() int {
	return len(l)
}

func (l associatedKeyFrameInfoList) Less(i, j int) bool {
	return l[i].ID < l[j].ID
}

func (l associatedKeyFrameInfoList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (d *managerLocalData) MarshalJSON() ([]byte, error) {
	temp := managerLocalDataForJSON{
		managerLocalData: *d,
		Chunks:           make([]associatedChunkInfo, 0, len(d.chunks)),
		KeyFrames:        make([]associatedKeyFrameInfo, 0, len(d.keyframes)),
	}

	for _, c := range d.chunks {
		temp.Chunks = append(temp.Chunks, c)
	}

	for _, kf := range d.keyframes {
		temp.KeyFrames = append(temp.KeyFrames, kf)
	}
	sort.Sort(associatedChunkInfoList(temp.Chunks))
	sort.Sort(associatedKeyFrameInfoList(temp.KeyFrames))

	return json.Marshal(temp)
}

func (d *managerLocalData) UnmarshalJSON(text []byte) error {
	temp := &managerLocalDataForJSON{}

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

// Download fetches from the lol spectator server data of a game
// identified by its region and ID, and save it on the local hardrive
func (m *LocalManager) Download(region *lol.Region, id lol.GameID) error {
	d, err := newReplayDataDir(m.datadir, region, id)
	if err != nil {
		return err
	}

	api, err := NewSpectateAPI(region, id)
	if err != nil {
		return err
	}

	//now we get the data, we should in a loop :
	//1. getLastChunkInfo
	//2. compute from it next time there will be chunk available
	//3. get All Chunk and Keyframe available
	//4. if last chunk is available, downlaod it and break to 6
	//5. wait until specified time, repeat from 1
	//6. Compute the metadata data to connect from game starting 0:00

	//to serve, we should :
	//serve getMetaData, making it believe that
	nextChunkToDownload := 0
	nextKeyframeToDownload := 0

	maData := newManagerLocalData()

	for {
		var metadata GameMetadata
		err = api.Get(GetGameMetaData, 1, &metadata)
		if err != nil {
			return err
		}

		var cInfo LastChunkInfo
		err := api.Get(GetLastChunkInfo, 1, &cInfo)
		if err != nil {
			return err
		}
		nextAvailableChunkDate := time.Now().Add(cInfo.NextAvailableChunk.Duration() + cInfo.Duration.Duration()/10)

		maData.MergeFromMetaData(metadata)
		maData.MergeFromLastChunkInfo(cInfo)
		maData.Consolidate()

		if err != nil {
			return err
		}

		for ; nextChunkToDownload <= int(cInfo.ID); nextChunkToDownload++ {
			chunkPath := d.chunkPath(ChunkID(nextChunkToDownload))
			if err := m.downloadBinary(api, GetGameDataChunk, nextChunkToDownload, chunkPath); err != nil {
				return err
			}
		}

		for ; nextKeyframeToDownload <= int(cInfo.AssociatedKeyFrameID); nextKeyframeToDownload++ {
			keyFramePath := d.keyFramePath(KeyFrameID(nextKeyframeToDownload))
			if err := m.downloadBinary(api, GetKeyFrame, nextKeyframeToDownload, keyFramePath); err != nil {
				return err
			}
		}

		//saves the metadata
		//erases the pending info, we recompute it at this end
		metadata.PendingAvailableChunkInfo = []ChunkInfo{}
		metadata.PendingAvailableKeyFrameInfo = []KeyFrameInfo{}
		err = m.saveJSON(d.metaDataPath(), metadata)
		if err != nil {
			return err
		}

		err = m.saveJSON(d.managerDataPath(), maData)
		if err != nil {
			return err
		}

		if cInfo.EndGameChunkID > 0 && nextChunkToDownload > int(cInfo.EndGameChunkID) {
			log.Printf("End of game detected and reached at %d", cInfo.EndGameChunkID)
			break
		}

		cTime := time.Now()
		if cTime.After(nextAvailableChunkDate) == true {
			continue
		}
		log.Printf("Waiting until %s", nextAvailableChunkDate)
		time.Sleep(nextAvailableChunkDate.Sub(cTime))

	}

	f, err := os.Create(d.endOfGameDataPath())
	if err != nil {
		return err
	}
	defer f.Close()
	resp, err := http.Get(api.Format(EndOfGameStats, NullParam))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return lol.RESTError{Code: resp.StatusCode}
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	err = json.Indent(&buf, data, "", "  ")
	if err != nil {
		return err
	}
	fmt.Fprintf(f, "%s", buf.String())
	return nil
}
