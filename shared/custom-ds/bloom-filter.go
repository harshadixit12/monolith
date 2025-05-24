package customds

import (
	"hash"
	"sync"

	"github.com/spaolacci/murmur3"
)

var defaultHashFns = 5
var defaultBitSliceSize = 500

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

func NewBloomFilter(size int) *BloomFilter {
	b := &BloomFilter{bitfield: make([]bool, size), hashers: getHashes(defaultHashFns)}

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
