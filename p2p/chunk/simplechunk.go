package chunk

import "math"

type ChunkPacket struct {
	msghash string
	chunkID int
	data    []byte
}

type simpleChunk struct {
	maxSize int
}

func (s simpleChunk) MaxSize(i int) simpleChunk {
	s.maxSize = i
	return s
}

func NewSimpleChunk() *simpleChunk {
	return &simpleChunk{
		maxSize: 100 * 1024,
	}
}

func (s *simpleChunk) Split(data []byte) [][]byte {
	chunkSize := int(math.Ceil(float64(len(data)) / float64(s.maxSize)))
	res := make([][]byte, chunkSize)
	for i := range res {
		if i == chunkSize-1 { //last chunk
			res[i] = data[i*s.maxSize:]
		} else {
			res[i] = data[i*s.maxSize : (i+1)*s.maxSize]
		}
	}
	return res
}

func (s *simpleChunk) Join(data [][]byte) []byte {
	size := 0
	for i := range data {
		size += len(data[i])
	}
	res := []byte{}
	for i := range data {
		res = append(res, data[i]...)
	}
	return res
}

func (s *simpleChunk) VerifyData(data [][]byte) []byte {
	return nil
}

func (s *simpleChunk) ReceiveFull(chunks []*ChunkPacket, hash string, size int) bool {
	return true
}
