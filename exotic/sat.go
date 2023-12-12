package exotic

import "github.com/dannybabbev/ord-api/ordinals"

const DificultyAdjustmentInterval = 2016
const HalvingInterval = 210000
const CycleInterval = HalvingInterval * 6

var PizzaRanges = ReadRangesFromOrdResponse(PIZZA_RANGES)
var NakamotoBlocks = []int64{9, 286, 688, 877, 1760, 2459, 2485, 3479, 5326, 9443, 9925, 10645, 14450, 15625, 15817, 19093, 23014, 28593, 29097}
var FirstTransactionRanges = []*ordinals.Range{
	{
		Start: 45000000000,
		Size:  1000000000,
	},
}

var HitmanRanges = SatingRangesToOrdinalsRanges(hitmanSatingRanges)
var JpegRanges = SatingRangesToOrdinalsRanges(jpegSatingRanges)

type Sat int64

func (s Sat) Epoch() Epoch {
	return EpochFromSat(s)
}

func (s Sat) Cycle() int64 {
	return s.Height() / CycleInterval
}

func (s Sat) Period() int64 {
	return s.Height() / DificultyAdjustmentInterval
}

func (s Sat) EpochPosition() int64 {
	r := s - s.Epoch().GetStartingSat()
	return int64(r)
}

func (s Sat) Height() int64 {
	r := int64(s.Epoch()) * HalvingInterval
	sub := s.Epoch().GetSubsidy()
	p := s.EpochPosition() / sub
	return p + r
}

func (s Sat) IsFirstSatInBlock() bool {
	sub := s.Epoch().GetSubsidy()
	return int64(s)%sub == 0
}

func (s Sat) GetRodarmorRarity() Satribute {
	isFirstSatInBlock := s.IsFirstSatInBlock()
	h := s.Height()

	if s == 0 {
		return Mythic
	}

	if isFirstSatInBlock && h%CycleInterval == 0 {
		return Legendary
	}

	if isFirstSatInBlock && h%HalvingInterval == 0 {
		return Epic
	}

	if isFirstSatInBlock && h%DificultyAdjustmentInterval == 0 {
		return Rare
	}

	if isFirstSatInBlock {
		return Uncommon
	}

	return Common
}

func (s Sat) IsBlack() bool {
	return Sat(s + 1).IsFirstSatInBlock()
}

func (s Sat) IsAlpha() bool {
	return s%1e8 == 0
}

func (s Sat) IsOmega() bool {
	return Sat(s + 1).IsAlpha()
}

func (s Sat) IsFibonacci() bool {
	a := int64(0)
	b := int64(1)
	next := a + b
	for next < int64(s) {
		next = a + b
		a = b
		b = next
	}

	if s == 0 || s == 1 {
		return true
	}

	return int64(s) == next
}

func (s Sat) Satributes() []Satribute {
	exotics := make([]Satribute, 0)

	rarity := s.GetRodarmorRarity()
	if rarity != Common {
		exotics = append(exotics, rarity)
	}

	if s.IsBlack() {
		exotics = append(exotics, Black)
	}

	if s.IsAlpha() {
		exotics = append(exotics, Alpha)
	}

	if s.IsOmega() {
		exotics = append(exotics, Omega)
	}

	if s.IsFibonacci() {
		exotics = append(exotics, Fibonacci)
	}

	if IsSatInRange(PizzaRanges, s) {
		exotics = append(exotics, Pizza)
	}

	if IsSatInRange(FirstTransactionRanges, s) {
		exotics = append(exotics, FirstTransaction)
	}

	if IsSatInRange(HitmanRanges, s) {
		exotics = append(exotics, Hitman)
	}

	if IsSatInRange(JpegRanges, s) {
		exotics = append(exotics, Jpeg)
	}

	h := s.Height()
	if h == 9 {
		exotics = append(exotics, Block9)
	}

	if h == 78 {
		exotics = append(exotics, Block78)
	}

	if h <= 1000 {
		exotics = append(exotics, Vintage)
	}

	if IsInBlocks(NakamotoBlocks, h) {
		exotics = append(exotics, Nakamoto)
	}

	return exotics
}
