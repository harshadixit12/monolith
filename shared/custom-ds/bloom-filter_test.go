package customds

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBloomFilterHappyPath(t *testing.T) {
	testCases := []struct {
		inputs   []string
		check    string
		negative string
	}{
		{
			inputs:   []string{"dog", "cat", "lion", "wolf"},
			check:    "cat",
			negative: "dolphin",
		},
	}

	for _, tc := range testCases {
		b := NewBloomFilter(1000)

		for _, val := range tc.inputs {
			b.Add([]byte(val))
		}

		assert.Equal(t, b.Check([]byte(tc.check)), true)
		assert.Equal(t, b.Check([]byte(tc.negative)), false)
	}
}
