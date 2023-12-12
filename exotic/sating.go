package exotic

import "github.com/dannybabbev/ord-api/ordinals"

// SatingRange is a range of sats as defined by sating.io
type SatingRange struct {
	start int64
	end   int64
}

func (r *SatingRange) ToOrdinalsRange() *ordinals.Range {
	return &ordinals.Range{
		Start: r.start,
		Size:  r.end - r.start,
	}
}

func SatingRangesToOrdinalsRanges(ranges []*SatingRange) []*ordinals.Range {
	var result []*ordinals.Range
	for _, r := range ranges {
		result = append(result, r.ToOrdinalsRange())
	}
	return result
}
