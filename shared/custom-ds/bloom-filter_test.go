package customds

import (
	"fmt"
	"math/rand/v2"
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
		b := NewBloomFilter(1000, 0.01)

		for _, val := range tc.inputs {
			b.Add([]byte(val))
		}

		assert.Equal(t, b.Check([]byte(tc.check)), true)
		assert.Equal(t, b.Check([]byte(tc.negative)), false)
	}
}

func BenchmarkBloomFilterAdd(b *testing.B) {
	tests := []struct {
		filterSize int
	}{
		{
			filterSize: 1000,
		},
		{
			filterSize: 10000,
		},
		{
			filterSize: 100000,
		},
		{
			filterSize: 1000000,
		},
	}

	for _, tc := range tests {
		bf := NewBloomFilter(tc.filterSize, 0.01)
		b.Run(fmt.Sprintf("input_size_%d", tc.filterSize), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bf.Add([]byte(fmt.Sprintf("input_%d", rand.Int())))
			}
		})
	}
}

func BenchmarkBloomFilterCheck(b *testing.B) {
	tests := []struct {
		filterSize int
	}{
		{
			filterSize: 1000,
		},
		{
			filterSize: 10000,
		},
		{
			filterSize: 100000,
		},
		{
			filterSize: 1000000,
		},
	}

	for _, tc := range tests {
		bf := NewBloomFilter(tc.filterSize, 0.01)
		bf.Add([]byte("hello world"))
		b.Run(fmt.Sprintf("input_size_%d", tc.filterSize), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bf.Check([]byte(fmt.Sprintf("input_%d", rand.Int())))
			}
		})
	}
}
