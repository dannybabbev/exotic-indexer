package server

import (
	"fmt"

	"github.com/bitgemtech/ord-api/esplora"
	"github.com/bitgemtech/ord-api/exotic"
	"github.com/bitgemtech/ord-api/ordinals"
	log "github.com/sirupsen/logrus"
)

type Indexer interface {
	GetOrdinalsForUTXO(utxo string) ([]*ordinals.Range, error)
}

type EsploraAPI interface {
	GetUTXOs(address string) ([]*esplora.UTXO, error)
}

type ServerModel struct {
	indexer Indexer
	esplora EsploraAPI
}

func NewServerModel(indexer Indexer, esplora EsploraAPI) *ServerModel {
	return &ServerModel{
		indexer: indexer,
		esplora: esplora,
	}
}

type UTXORangesResponse struct {
	Ranges       []*RangeResponse       `json:"ranges"`
	ExoticRanges []*ExoticRangeResponse `json:"exoticRanges"`
}

func (s *ServerModel) GetUTXORanges(utxos []string, excludeCommonRanges bool) (*UTXORangesResponse, error) {
	ranges := make([]*RangeResponse, 0)
	exoticRanges := make([]*ExoticRangeResponse, 0)
	for _, utxo := range utxos {
		res, err := s.indexer.GetOrdinalsForUTXO(utxo)
		if err != nil {
			log.WithFields(log.Fields{
				"utxo": utxo,
				"err":  err,
			}).Info("Error getting ordinals for utxo")

			return nil, err
		}

		// Caluclate the offset for each range
		offset := int64(0)
		for _, r := range res {
			rr := &RangeResponse{
				Start:  r.Start,
				Size:   r.Size,
				End:    r.Start + r.Size,
				Utxo:   utxo,
				Offset: offset,
			}
			ranges = append(ranges, rr)
			offset += r.Size
		}

		sr := exotic.FindExoticRangesUTXO(res)
		for _, r := range sr {
			srr := &ExoticRangeResponse{
				RangeResponse: RangeResponse{
					Start:  r.Range.Start,
					Size:   r.Range.Size,
					End:    r.Range.Start + r.Range.Size,
					Utxo:   utxo,
					Offset: r.Offset,
				},
				Satributes: r.Satributes,
			}
			exoticRanges = append(exoticRanges, srr)
		}
	}

	if excludeCommonRanges {
		return &UTXORangesResponse{
			ExoticRanges: exoticRanges,
		}, nil
	}

	return &UTXORangesResponse{
		Ranges:       ranges,
		ExoticRanges: exoticRanges,
	}, nil
}

func (s *ServerModel) GetAddressRanges(address string, excludeCommonRanges bool) (*UTXORangesResponse, error) {
	utxos, err := s.esplora.GetUTXOs(address)
	if err != nil {
		return nil, err
	}

	utxoStrings := make([]string, 0)
	for _, utxo := range utxos {
		// Unconfirmed UTXOs are not indexed
		if utxo.Status.Confirmed == false {
			continue
		}
		utxoStrings = append(utxoStrings, fmt.Sprintf("%s:%d", utxo.Txid, utxo.Vout))
	}

	return s.GetUTXORanges(utxoStrings, excludeCommonRanges)
}

type SatResponse struct {
	Sat        int64              `json:"sat"`
	Height     int64              `json:"height"`
	Cycle      int64              `json:"cycle"`
	Epoch      int64              `json:"epoch"`
	Period     int64              `json:"period"`
	Satributes []exotic.Satribute `json:"satributes"`
}

func (s *ServerModel) GetSat(sat int64) *SatResponse {
	sm := exotic.Sat(sat)

	return &SatResponse{
		Sat:        int64(sm),
		Height:     sm.Height(),
		Epoch:      int64(sm.Epoch()),
		Cycle:      int64(sm.Cycle()),
		Period:     int64(sm.Period()),
		Satributes: sm.Satributes(),
	}
}
