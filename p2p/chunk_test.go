package p2p

import (
	"fmt"
	"node/p2p/chunk"
	"testing"
)

func TestChunkManager(t *testing.T) {
	manager := NewChunkManager().SetEngine(chunk.NewSimpleChunk().MaxSize(5))
	data := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	chunkPackets, err := manager.Split(data)
	if err != nil {
		fmt.Println(err)
	}

	manager.Receive(chunkPackets[0])
	manager.Receive(chunkPackets[2])

	isFull, err := manager.Receive(chunkPackets[1])

	if err != nil {
		fmt.Println("err", err)
		return
	}
	if isFull == nil {
		fmt.Println("Receive not full")
	} else {
		fmt.Println("Receive full", isFull)
	}
}
