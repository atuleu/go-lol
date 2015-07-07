package xlol

import (
	"sort"

	. "gopkg.in/check.v1"
)

type ReplaySuite struct{}

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
