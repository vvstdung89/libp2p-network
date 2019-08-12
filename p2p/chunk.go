package p2p

/*
TODO: check for old chunk message and delete it
*/

import (
	"errors"
	"github.com/gogo/protobuf/proto"
	p2p_chunk "node/p2p/chunk"
	"sync"
	"time"
)

type ChunkManager struct {
	engine       ChunkEngine
	chunkCache   map[string]*ChunkCache
	chunkCacheMu sync.Mutex
}

type ChunkCache struct {
	packetHash   []byte
	time         time.Time
	chunks       []*p2p_chunk.ChunkPacketPB
	chunkSize    int
	chunkReceive int
}

type ChunkEngine interface {
	Split(data []byte) ([][]byte, error)
	VerifyData(data *p2p_chunk.ChunkPacketPB) bool
	ReceiveFull(chunks []*p2p_chunk.ChunkPacketPB, hash []byte) (res []byte, err error)
}

func NewChunkManager() *ChunkManager {
	res := &ChunkManager{
		chunkCache: make(map[string]*ChunkCache),
	}
	return res
}

func (s *ChunkManager) SetEngine(engine ChunkEngine) *ChunkManager {
	s.engine = engine
	return s
}

func (s *ChunkManager) Split(data []byte) ([][]byte, error) {
	return s.engine.Split(data)
}

func (s *ChunkManager) Receive(data []byte) (fullData []byte, err error) {
	chunkPB := &p2p_chunk.ChunkPacketPB{}
	err = proto.Unmarshal(data, chunkPB)
	if err != nil {
		return nil, errors.New("Cannot deserialize data")
	}

	isValid := s.engine.VerifyData(chunkPB)
	if isValid {
		s.chunkCacheMu.Lock()
		defer s.chunkCacheMu.Unlock()

		msgChunkCache := s.chunkCache[string(chunkPB.BloomData)]
		if msgChunkCache == nil {
			msgChunkCache = &ChunkCache{
				time:         time.Now(),
				chunkReceive: 0,
				chunks:       make([]*p2p_chunk.ChunkPacketPB, chunkPB.ChunkSize),
				chunkSize:    int(chunkPB.ChunkSize),
				packetHash:   chunkPB.BloomData,
			}
		}
		if msgChunkCache.chunks[chunkPB.ChunkId] == nil {
			msgChunkCache.chunks[chunkPB.ChunkId] = chunkPB
			msgChunkCache.chunkReceive++
		}
		s.chunkCache[string(chunkPB.BloomData)] = msgChunkCache
		if msgChunkCache.chunkReceive == msgChunkCache.chunkSize {
			fullData, err = s.engine.ReceiveFull(msgChunkCache.chunks, chunkPB.BloomData)
			delete(s.chunkCache, string(chunkPB.BloomData))
		}
	} else {
		err = errors.New("Chunk is not valid")
	}
	return fullData, err
}
