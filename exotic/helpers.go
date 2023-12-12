package exotic

import (
	"encoding/json"

	"github.com/dannybabbev/ord-api/ordinals"
)

type Item struct {
	Output string `json:"output"`
	Start  int64  `json:"start"`
	End    int64  `json:"end"`
	Size   int64  `json:"size"`
	Offset int64  `json:"offset"`
	Rarity string `json:"rarity"`
	Name   string `json:"name"`
}

func ReadRangesFromOrdResponse(fileContent string) []*ordinals.Range {
	// Parse the JSON content into a slice of the struct
	var items []Item
	err := json.Unmarshal([]byte(fileContent), &items)
	if err != nil {
		panic(err)
	}

	res := []*ordinals.Range{}
	for _, item := range items {
		r := &ordinals.Range{
			Start: item.Start,
			Size:  item.Size,
		}
		res = append(res, r)
	}

	return res
}

func IsInBlocks(blocks []int64, height int64) bool {
	for _, b := range blocks {
		if b == height {
			return true
		}
	}
	return false
}

func IsSatInRange(ranges []*ordinals.Range, sat Sat) bool {
	for _, r := range ranges {
		if int64(sat) >= r.Start && int64(sat) < r.Start+r.Size {
			return true
		}
	}
	return false
}

func GetRangeForBlock(height int) []*ordinals.Range {
	start := ordinals.FirstOrdinal(height)
	size := ordinals.Subsidy(height)

	r := &ordinals.Range{
		Start: start,
		Size:  size,
	}

	return []*ordinals.Range{r}
}

func GetRangesForBlocks(blocks []int) []*ordinals.Range {
	ranges := []*ordinals.Range{}
	for _, b := range blocks {
		ranges = append(ranges, GetRangeForBlock(b)...)
	}
	return ranges
}

func IsRodarmorRare(s Satribute) bool {
	if s == Mythic || s == Legendary || s == Epic || s == Rare || s == Uncommon {
		return true
	}

	return false
}
