package exotic

import (
	"testing"

	"github.com/bitgemtech/exotic-indexer/ordinals"
	"github.com/stretchr/testify/assert"
)

func TestIsSatInRange(t *testing.T) {
	ranges := []*ordinals.Range{
		{
			Start: 100,
			Size:  100,
		},
		{
			Start: 250,
			Size:  50,
		},
	}

	res := IsSatInRange(ranges, 100)
	assert.True(t, res)

	res = IsSatInRange(ranges, 250)
	assert.True(t, res)

	res = IsSatInRange(ranges, 200)
	assert.False(t, res)

	res = IsSatInRange(ranges, 500)
	assert.False(t, res)
}
