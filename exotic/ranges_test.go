package exotic

import (
	"testing"

	"github.com/bitgemtech/ord-api/ordinals"
	"github.com/stretchr/testify/assert"
)

func TestFindSpecialRangesUTXO(t *testing.T) {
	tabletest := []struct {
		ranges []*ordinals.Range
		want   []*ExoticRange
	}{
		{
			[]*ordinals.Range{
				{
					Start: 1953760625000000,
					Size:  625000000,
				},
			},
			[]*ExoticRange{
				{
					Range: &ordinals.Range{
						Start: 1953760625000000,
						Size:  1,
					},
					Offset:     0,
					Satributes: []Satribute{Uncommon},
				},
				{
					Range: &ordinals.Range{
						Start: 1953761249999999,
						Size:  1,
					},
					Offset:     624999999,
					Satributes: []Satribute{Black},
				},
			},
		},
		{
			[]*ordinals.Range{
				{
					Start: 1096735000000000,
					Size:  20000,
				},
				{
					Start: 912312093123000,
					Size:  100,
				},
				{
					Start: 390660767841,
					Size:  1000,
				},
			},
			[]*ExoticRange{
				{
					Range: &ordinals.Range{
						Start: 1096735000000000,
						Size:  1,
					},
					Offset:     0,
					Satributes: []Satribute{Uncommon, Alpha},
				},
				{
					Range: &ordinals.Range{
						Start: 390660767841,
						Size:  1000,
					},
					Offset:     20100,
					Satributes: []Satribute{Block78, Vintage},
				},
			},
		},
		{
			[]*ordinals.Range{
				{
					Start: 162062961756935,
					Size:  250000,
				},
			},
			[]*ExoticRange{
				{
					Range: &ordinals.Range{
						Start: 162062961756935,
						Size:  250000,
					},
					Offset:     0,
					Satributes: []Satribute{Jpeg},
				},
			},
		},
	}

	for _, tt := range tabletest {
		r := FindExoticRangesUTXO(tt.ranges)
		assert.Equal(t, tt.want, r)
	}
}
