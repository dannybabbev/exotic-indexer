package indexer

type MockBitcoind struct {
	BlockCountReturn uint64
	BlockHashReturn  map[uint64]string
	RawBlockReturn   map[string]string
}

func (m *MockBitcoind) GetBlockCount() (uint64, error) {
	return m.BlockCountReturn, nil
}

func (m *MockBitcoind) GetBlockHash(blockHeight uint64) (string, error) {
	return m.BlockHashReturn[blockHeight], nil
}

func (m *MockBitcoind) GetRawBlock(blockHash string) (string, error) {
	return m.RawBlockReturn[blockHash], nil
}
