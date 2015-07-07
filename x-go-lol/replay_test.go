package xlol

import (
	"bytes"
	"encoding/json"
	"sort"

	. "gopkg.in/check.v1"
)

type ReplaySuite struct {
	stubGM    GameMetadata
	stubCInfo LastChunkInfo
}

var _ = Suite(&ReplaySuite{})

func (s *ReplaySuite) TestChunkListAreSortedByID(c *C) {
	ids := []int{20, 2, 103, 48, 50}
	data := make([]Chunk, 0, len(ids))
	for _, id := range ids {
		data = append(data, Chunk{
			ChunkInfo: ChunkInfo{
				ID: ChunkID(id),
			},
		})
	}

	sort.Ints(ids)
	sort.Sort(ChunkList(data))
	c.Assert(len(data), Equals, len(ids))

	for i, idInt := range ids {
		id := ChunkID(idInt)
		c.Check(data[i].ID, Equals, id)
	}

}

func (s *ReplaySuite) TestKeyFrameListAreSortedByID(c *C) {
	ids := []int{20, 2, 103, 48, 50}
	data := make([]KeyFrame, 0, len(ids))
	for _, id := range ids {
		data = append(data, KeyFrame{
			KeyFrameInfo: KeyFrameInfo{
				ID: KeyFrameID(id),
			},
		})
	}

	sort.Ints(ids)
	sort.Sort(KeyFrameList(data))
	c.Assert(len(data), Equals, len(ids))

	for i, idInt := range ids {
		id := KeyFrameID(idInt)
		c.Check(data[i].ID, Equals, id)
	}

}

func (s *ReplaySuite) SetUpSuite(c *C) {
	getMetaDataJSON, err := Asset("data/getGameMetaData.json")
	c.Assert(err, IsNil)
	dec := json.NewDecoder(bytes.NewBuffer(getMetaDataJSON))
	err = dec.Decode(&(s.stubGM))
	c.Assert(err, IsNil)

	getLastChunkInfoJSON, err := Asset("data/getLastChunkInfo.json")
	c.Assert(err, IsNil)
	dec = json.NewDecoder(bytes.NewBuffer(getLastChunkInfoJSON))
	err = dec.Decode(&(s.stubCInfo))
	c.Assert(err, IsNil)

}

func (s *ReplaySuite) TestMergeDataShouldNotPanic(c *C) {
	defer func() {
		if r := recover(); r != nil {
			c.Fatalf("Recovered a panic from merge: %s", r)
		}
	}()

	replay := NewEmptyReplay()

	replay.MergeFromMetaData(s.stubGM)
	replay.MergeFromLastChunkInfo(s.stubCInfo)
	replay.Consolidate()

	//test the data is correctly merged
	newGM := s.stubGM
	newGM.PendingAvailableChunkInfo = []ChunkInfo{}
	newGM.PendingAvailableKeyFrameInfo = []KeyFrameInfo{}

	expectedJSON, err := json.Marshal(newGM)
	c.Assert(err, IsNil)
	JSON, err := json.Marshal(replay.MetaData)
	c.Assert(err, IsNil)

	c.Check(string(JSON), Equals, string(expectedJSON))

}
