package p2p

import "node/p2p/chunk"

type ChunkManager struct {
	engine     ChunkEngine
	chunkCache map[string]ChunkCache
}

type ChunkCache struct {
	packetHash   string
	time         int
	chunks       []*chunk.ChunkPacket
	chunkSize    int
	chunkReceive int
}

type ChunkEngine interface {
	Split(data []byte) [][]byte
	Join(data [][]byte) []byte
	VerifyData(data [][]byte) []byte
	ReceiveFull(chunks []*chunk.ChunkPacket, hash string, size int) bool
}

func NewChunkManager() *ChunkManager {
	res := &ChunkManager{
		chunkCache: make(map[string]ChunkCache),
	}
	return res
}

func (s *ChunkManager) SetEngine(engine ChunkEngine) *ChunkManager {
	s.engine = engine
	return s
}

func (s *ChunkManager) Split(data []byte) [][]byte {
	splitData := s.engine.Split(data)
	chunks := [][]byte{}
	for i, chunk := range splitData {
		//serialize using protobuf

	}
	return chunks
}

func (s *ChunkManager) join(data [][]byte) []byte {
	return nil
}

func (s *ChunkManager) Receive(data []byte) {
	//deserialize data using protobuf
	//add to ChunkCache
	//check if ChunkCache is full ->
}
