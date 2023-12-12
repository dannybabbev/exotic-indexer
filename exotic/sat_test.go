package exotic

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSatEpoch(t *testing.T) {
	tableTest := []struct {
		sat       Sat
		wantEpoch Epoch
	}{
		{
			Sat(0),
			Epoch(0),
		},
		{
			Sat(1953651875000000),
			Epoch(3),
		},
	}

	for _, tt := range tableTest {
		assert.Equal(t, tt.wantEpoch, tt.sat.Epoch())
	}
}

func TestSatCycle(t *testing.T) {
	tableTest := []struct {
		sat       Sat
		wantCycle int64
	}{
		{
			Sat(0),
			0,
		},
		{
			Sat(1953651875000000),
			0,
		},
		{
			Sat(2099999997689999),
			5,
		},
	}

	for _, tt := range tableTest {
		assert.Equal(t, tt.wantCycle, tt.sat.Cycle())
	}
}

func TestSatPeriod(t *testing.T) {
	tableTest := []struct {
		sat        Sat
		wantPeriod int64
	}{
		{
			Sat(0),
			0,
		},
		{
			Sat(1953651875000000),
			404,
		},
	}

	for _, tt := range tableTest {
		assert.Equal(t, tt.wantPeriod, tt.sat.Period())
	}
}

func TestSatHeight(t *testing.T) {
	tableTest := []struct {
		sat        Sat
		wantHeight int64
	}{
		{
			Sat(0),
			0,
		},
		{
			Sat(1805619726571369),
			604495,
		},
		{
			Sat(2099999997689999),
			6929999,
		},
		{
			Sat(1050000000000000),
			210000,
		},
	}

	for _, tt := range tableTest {
		assert.Equal(t, tt.wantHeight, tt.sat.Height())
	}
}

func TestSatGetRodarmorRarity(t *testing.T) {
	tableTest := []struct {
		sat    Sat
		rarity Satribute
	}{
		{
			Sat(1953568750000000),
			Uncommon,
		},
		{
			Sat(45027647018),
			Common,
		},
		{
			Sat(0),
			Mythic,
		},
		{
			Sat(1),
			Common,
		},
		{
			Sat(100800000000000),
			Rare,
		},
		{
			Sat(1030915000000000),
			Uncommon,
		},
		{
			Sat(1109640000000000),
			Rare,
		},
		{
			Sat(1109640000000001),
			Common,
		},
		{
			Sat(1937670000000000),
			Rare,
		},
		{
			Sat(1948833750000000),
			Uncommon,
		},
		{
			Sat(1847158750000000),
			Uncommon,
		},
		{
			Sat(1050000000000000),
			Epic,
		},
	}

	for _, tt := range tableTest {
		assert.Equal(t, tt.rarity, tt.sat.GetRodarmorRarity())
	}
}

func TestSatIsBlack(t *testing.T) {
	tableTest := []struct {
		sat  Sat
		want bool
	}{
		{
			Sat(1306054999999999),
			true,
		},
		{
			Sat(1306055000000000),
			false,
		},
	}

	for _, tt := range tableTest {
		assert.Equal(t, tt.want, tt.sat.IsBlack())
	}
}

func TestSatIsAlpha(t *testing.T) {
	tableTest := []struct {
		sat  Sat
		want bool
	}{
		{
			Sat(1904260000000000),
			true,
		},
		{
			Sat(100000001),
			false,
		},
	}

	for _, tt := range tableTest {
		assert.Equal(t, tt.want, tt.sat.IsAlpha())
	}
}

func TestSatIsFibonacci(t *testing.T) {
	tt := []struct {
		sat  Sat
		want bool
	}{
		{
			Sat(0),
			true,
		},
		{
			Sat(1),
			true,
		},
		{
			Sat(2),
			true,
		},
		{
			Sat(3),
			true,
		},
		{
			Sat(5),
			true,
		},
		{
			Sat(8),
			true,
		},
		{
			Sat(9),
			false,
		},
		{
			Sat(233),
			true,
		},
		{
			Sat(234),
			false,
		},
	}

	for _, tt := range tt {
		assert.Equal(t, tt.want, tt.sat.IsFibonacci(), fmt.Sprintf("exepected sat %d to be %t", int64(tt.sat), tt.want))
	}
}

func TestSatIsOmega(t *testing.T) {
	tableTest := []struct {
		sat  Sat
		want bool
	}{
		{
			Sat(1904260000000001),
			false,
		},
		{
			Sat(1879153999999999),
			true,
		},
	}

	for _, tt := range tableTest {
		assert.Equal(t, tt.want, tt.sat.IsOmega())
	}
}

func TestFindSatributes(t *testing.T) {
	tableTest := []struct {
		sat            Sat
		wantSatributes []Satribute
	}{
		{
			283934082571183,
			[]Satribute{Pizza},
		},
		{
			45027647018,
			[]Satribute{FirstTransaction, Block9, Vintage, Nakamoto},
		},
		{
			1937670000000000,
			[]Satribute{Rare, Alpha},
		},
		{
			1584663749999999,
			[]Satribute{Black},
		},
		{
			1904260000000000,
			[]Satribute{Uncommon, Alpha},
		},
		{
			1847158750000000,
			[]Satribute{Uncommon},
		},
		{
			1765168199999999,
			[]Satribute{Omega},
		},
		{
			145808809963218,
			[]Satribute{Hitman},
		},
		{
			942682707984987,
			[]Satribute{Hitman},
		},
		{
			162060995322118,
			[]Satribute{Jpeg},
		},
		{
			162062962006935,
			[]Satribute{Jpeg},
		},
		{
			27777890035288,
			[]Satribute{Fibonacci},
		},
	}

	for _, tt := range tableTest {
		assert.Equal(t, tt.wantSatributes, tt.sat.Satributes())
	}
}
