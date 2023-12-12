package exotic

import "github.com/dannybabbev/ord-api/ordinals"

type ExoticRange struct {
	Range      *ordinals.Range `json:"range"`
	Offset     int64           `json:"offset"`
	Satributes []Satribute     `json:"satributes"`
}

// FindExoticRangesUTXO finds and returns the special ranges
// Has to be called only for one UTXO, since it calculates the offset
func FindExoticRangesUTXO(ranges []*ordinals.Range) []*ExoticRange {
	res := []*ExoticRange{}

	offset := int64(0)
	for _, r := range ranges {
		// Check the first sat
		sat := Sat(r.Start)
		satributes := sat.Satributes()

		if len(satributes) > 0 && (IsRodarmorRare(satributes[0]) || sat.IsAlpha()) {
			sr := &ExoticRange{
				Range: &ordinals.Range{
					Start: int64(sat),
					Size:  1,
				},
				Offset:     offset,
				Satributes: satributes,
			}

			res = append(res, sr)
		} else if len(satributes) > 0 {
			sr := &ExoticRange{
				Range: &ordinals.Range{
					Start: int64(sat),
					Size:  r.Size,
				},
				Offset:     offset,
				Satributes: satributes,
			}

			res = append(res, sr)
		}

		if r.Size == 1 {
			offset += r.Size
			continue
		}

		o := r.Size - 1
		sat = Sat(r.Start + o)
		satributes = sat.Satributes()
		if len(satributes) > 0 && (sat.IsBlack() || sat.IsOmega()) {
			sr := &ExoticRange{
				Range: &ordinals.Range{
					Start: int64(sat),
					Size:  1,
				},
				Offset:     offset + o,
				Satributes: satributes,
			}

			res = append(res, sr)
		}

		offset += r.Size
	}

	return res
}
