package customds

import (
	"hash"
	"math"
	"sync"

	"github.com/spaolacci/murmur3"
)

type BloomInterface interface {
	Add([]byte)
	Test([]byte) bool
}

type BloomFilter struct {
	bitfield []bool
	hashers  []hash.Hash64
	mu       sync.Mutex
}

type Hasher interface {
	GetHashes(n uint64) []hash.Hash64
}

func NewBloomFilter(size int, p float64) *BloomFilter {
	m, k := getOptimalTuning(size, p)
	b := &BloomFilter{bitfield: make([]bool, m), hashers: getHashes(k)}

	return b
}

func (b *BloomFilter) Add(val []byte) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, h := range b.hashers {
		h.Reset()
		h.Write(val)
		hashValue := h.Sum64() % uint64(len(b.bitfield))
		b.bitfield[hashValue] = true
	}
}

func (b *BloomFilter) Check(val []byte) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, h := range b.hashers {
		h.Reset()
		h.Write(val)
		hashValue := h.Sum64() % uint64(len(b.bitfield))

		if !b.bitfield[hashValue] {
			return false
		}
	}

	return true
}

func getHashes(n int) []hash.Hash64 {
	hashes := make([]hash.Hash64, n)

	for i := 0; i < n; i++ {
		hashes[i] = murmur3.New64WithSeed(uint32(i)) // Seeding with different integers reduces probability of collisions.
	}

	return hashes
}

// Given p - false positive rate, this function returns optimal value for m (bit field size) and k (number of hash functions)
// m = -n*ln(p)/(ln(2)^2)
// k = m*ln(2)/n
func getOptimalTuning(size int, p float64) (int, int) {
	mFloat := -1 * float64(size) * math.Log(p) / (math.Pow(math.Ln2, 2))
	m := int(math.Ceil(mFloat))

	kFloat := float64(m) * math.Ln2 / float64(size)
	k := int(math.Ceil(kFloat))

	return m, k
}
