package chunk

import (
	"crypto/sha256"
	"errors"
	"github.com/gogo/protobuf/proto"
	"github.com/willf/bloom"
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
	filter := bloom.New(1280, 10)
	for _, chunk := range chunkData {
		filter.Add(hash.Sum(chunk))
	}

	for i, chunk := range chunkData {
		chunkPacket := &ChunkPacketPB{
			Data:      chunk,
			ChunkHash: hash.Sum(chunk),
			ChunkId:   int32(i),
			ChunkSize: int32(len(chunkData)),
		}
		bloomData, _ := filter.GobEncode()
		chunkPacket.BloomData = bloomData
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
	filter := bloom.New(1280, 10)
	filter.GobDecode(data.BloomData)
	if reflect.DeepEqual(hash.Sum(data.Data), data.ChunkHash) && filter.Test(data.ChunkHash) {
		return true
	} else {
		return false
	}
}

func (s simpleChunk) ReceiveFull(chunks []*ChunkPacketPB, hash []byte) (res []byte, err error) {
	if len(chunks) == 0 {
		return nil, errors.New("Chunk empty")
	}
	filter := bloom.New(1280, 10)
	filter.GobDecode(chunks[0].BloomData)
	for i := range chunks {
		if chunks[i].Data == nil {
			return nil, errors.New("Not enough data")
		}
		res = append(res, chunks[i].Data...)
	}
	return res, nil
}
