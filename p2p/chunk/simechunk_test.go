package chunk

import (
	"testing"
)

func Test_Chunk(t *testing.T) {
	engine := ChunkEngine{}
	var data [][]byte = engine.Chunk()
}
