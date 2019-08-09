package chunk

import (
	"fmt"
	"testing"
)

func TestSimpleChunk_Join(t *testing.T) {
	chunk := NewSimpleChunk().MaxSize(5)

	data := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	res := chunk.Split(data)
	fmt.Println(res)
}

func TestSimpleChunk_Split(t *testing.T) {
	chunk := NewSimpleChunk().MaxSize(5)
	data := [][]byte{[]byte{1, 2, 3, 4, 5}, []byte{1, 2, 3, 4, 5}, []byte{1, 2, 3}}
	res := chunk.Join(data)
	fmt.Println(res)
}
