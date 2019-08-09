package chunk

import (
	"crypto/sha256"
	"errors"
	"github.com/gogo/protobuf/proto"
	"math"
	"reflect"
)

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

func (s simpleChunk) Split(data []byte) ([][]byte, error) {
	hash := sha256.New()

	chunkSize := int(math.Ceil(float64(len(data)) / float64(s.maxSize)))
	chunkData := make([][]byte, chunkSize)
	for i := range chunkData {
		if i == chunkSize-1 { //last chunk
			chunkData[i] = data[i*s.maxSize:]
		} else {
			chunkData[i] = data[i*s.maxSize : (i+1)*s.maxSize]
		}
	}

	chunkPackets := make([][]byte, chunkSize)
	for i, chunk := range chunkData {
		chunkPacket := &ChunkPacketPB{
			Data:      chunk,
			ChunkHash: hash.Sum(chunk),
			ChunkId:   int32(i),
			ChunkSize: int32(len(chunkData)),
			MsgHash:   hash.Sum(data),
		}
		data, err := proto.Marshal(chunkPacket)
		if err != nil {
			return nil, err
		}
		chunkPackets[i] = data
	}

	return chunkPackets, nil
}

func (s simpleChunk) VerifyData(data *ChunkPacketPB) bool {
	hash := sha256.New()
	if reflect.DeepEqual(hash.Sum(data.Data), data.ChunkHash) {
		return true
	} else {
		return false
	}
}

func (s simpleChunk) ReceiveFull(chunks []*ChunkPacketPB, hash []byte) (res []byte, err error) {
	for i := range chunks {
		res = append(res, chunks[i].Data...)
	}
	hashRes := sha256.New()
	if reflect.DeepEqual(hashRes.Sum(res), hash) {
		return res, nil
	} else {
		return nil, errors.New("Check sum full nessage error")
	}

}
