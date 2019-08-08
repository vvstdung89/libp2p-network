package chunk

type SimpleChunk struct {
	simpleChunkConfig simpleChunkConfig
}

//pattern to create default value in struct
//not expose config, using function to get & set default at that function
//
type simpleChunkConfig struct {
	MAX_SIZE int
}

func NewSimpleChunkConfig() simpleChunkConfig {
	return simpleChunkConfig{
		MAX_SIZE: 100 * 1024,
	}
}

func NewSimpleChunk(cfg simpleChunkConfig) *SimpleChunk {
	return &SimpleChunk{
		simpleChunkConfig: cfg,
	}
}

func (s *SimpleChunk) Split(data []byte) [][]byte {
	return nil
}

func (s *SimpleChunk) Join(data [][]byte) []byte {
	return nil
}
