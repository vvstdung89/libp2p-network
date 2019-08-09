package chunk

import (
	"fmt"
	"testing"
)

func TestSimpleChunk(t *testing.T) {
	chunk := NewSimpleChunk().MaxSize(5)
	data := []byte{1, 2, 3, 4}
	res, err := chunk.Split(data)
	fmt.Println(res, err)
}
