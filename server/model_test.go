package server

import (
	"testing"

	"github.com/bitgemtech/ord-api/esplora"
	"github.com/bitgemtech/ord-api/exotic"
	"github.com/bitgemtech/ord-api/ordinals"
	"github.com/stretchr/testify/assert"
)

type MockIndexer struct {
	Ordinals map[string][]*ordinals.Range
}

func NewMockIndexer(ordinals map[string][]*ordinals.Range) *MockIndexer {
	return &MockIndexer{
		Ordinals: ordinals,
	}
}

func (m *MockIndexer) GetOrdinalsForUTXO(utxo string) ([]*ordinals.Range, error) {
	return m.Ordinals[utxo], nil
}

type MockEsploraAPI struct {
	UTXOs []*esplora.UTXO
}

func NewMockEsploraAPI(UTXOs []*esplora.UTXO) *MockEsploraAPI {
	return &MockEsploraAPI{
		UTXOs: UTXOs,
	}
}

func (m *MockEsploraAPI) GetUTXOs(address string) ([]*esplora.UTXO, error) {
	return m.UTXOs, nil
}

func GetMockIndexerWithUTXOs() *MockIndexer {
	mockResp := make(map[string][]*ordinals.Range)
	mockResp["utxo1:0"] = []*ordinals.Range{
		{
			Start: 1953760625000000,
			Size:  625000000,
		},
		{
			Start: 923804960941491,
			Size:  198741,
		},
	}

	mockResp["utxo2:0"] = []*ordinals.Range{
		{
			Start: 1096735000000000,
			Size:  20000,
		},
	}
	return NewMockIndexer(mockResp)
}

func TestGetUTXORanges(t *testing.T) {
	mockIndexer := GetMockIndexerWithUTXOs()
	model := NewServerModel(mockIndexer, nil)

	res, _ := model.GetUTXORanges([]string{"utxo1:0", "utxo2:0"}, false)

	wantRanges := []*RangeResponse{
		{
			Start:  1953760625000000,
			Size:   625000000,
			End:    1953761250000000,
			Utxo:   "utxo1:0",
			Offset: 0,
		},
		{
			Start:  923804960941491,
			Size:   198741,
			End:    923804961140232,
			Utxo:   "utxo1:0",
			Offset: 625000000,
		},
		{
			Start:  1096735000000000,
			Size:   20000,
			End:    1096735000020000,
			Utxo:   "utxo2:0",
			Offset: 0,
		},
	}

	wantExoticRanges := []*ExoticRangeResponse{
		{
			RangeResponse{
				Start:  1953760625000000,
				Size:   1,
				End:    1953760625000001,
				Utxo:   "utxo1:0",
				Offset: 0,
			},
			[]exotic.Satribute{exotic.Uncommon},
		},
		{
			RangeResponse{
				Start:  1953761249999999,
				Size:   1,
				End:    1953761250000000,
				Utxo:   "utxo1:0",
				Offset: 624999999,
			},
			[]exotic.Satribute{exotic.Black},
		},
		{
			RangeResponse{
				Start:  1096735000000000,
				Size:   1,
				End:    1096735000000001,
				Offset: 0,
				Utxo:   "utxo2:0",
			},
			[]exotic.Satribute{exotic.Uncommon, exotic.Alpha},
		},
	}

	assert.Equal(t, wantRanges, res.Ranges)
	assert.Equal(t, wantExoticRanges, res.ExoticRanges)
}

func TestGetAddressRanges(t *testing.T) {
	mockEsplora := NewMockEsploraAPI([]*esplora.UTXO{
		{
			Txid: "utxo1",
			Vout: 0,
			Status: esplora.UTXOStatus{
				Confirmed:   true,
				BlockHeight: 816017,
				BlockHash:   "0000000000000000000b1b1f7d0a588ac",
				BlockTime:   1624296000,
			},
		},
		{
			Txid: "utxo2",
			Vout: 0,
			Status: esplora.UTXOStatus{
				Confirmed: false,
			},
		},
	})

	mockIndexer := GetMockIndexerWithUTXOs()

	model := NewServerModel(mockIndexer, mockEsplora)

	res, _ := model.GetAddressRanges("address1", false)

	wantExoticRanges := []*ExoticRangeResponse{
		{
			RangeResponse{
				Start:  1953760625000000,
				Size:   1,
				End:    1953760625000001,
				Utxo:   "utxo1:0",
				Offset: 0,
			},
			[]exotic.Satribute{exotic.Uncommon},
		},
		{
			RangeResponse{
				Start:  1953761249999999,
				Size:   1,
				End:    1953761250000000,
				Utxo:   "utxo1:0",
				Offset: 624999999,
			},
			[]exotic.Satribute{exotic.Black},
		},
	}

	assert.Equal(t, wantExoticRanges, res.ExoticRanges)
}

func TestGetSat(t *testing.T) {
	model := NewServerModel(nil, nil)

	model.GetSat(1953760625000000)

	want := &SatResponse{
		Sat:    1953760625000000,
		Height: 816017,
		Cycle:  0,
		Epoch:  3,
		Period: 404,
		Satributes: []exotic.Satribute{
			exotic.Uncommon,
		},
	}

	assert.Equal(t, want, model.GetSat(1953760625000000))
}
