package saga

import (
	"math"
	"sync"
	"time"
)

const (
	tsLen     = 5 * 8
	cntLen    = 2 * 8
	suffixLen = tsLen + cntLen
)

type Generator struct {
	mu sync.Mutex
	// high order byte
	prefix uint64
	// low order 7 bytes
	suffix uint64
}

func NewGenerator(memberID uint8, now time.Time) *Generator {
	prefix := uint64(memberID) << suffixLen
	unixMilli := uint64(now.UnixNano()) / uint64(time.Millisecond/time.Nanosecond)
	suffix := lowbit(unixMilli, tsLen) << cntLen
	return &Generator{
		prefix: prefix,
		suffix: suffix,
	}
}

// Next generates a id that is unique.
func (g *Generator) Next() uint64 {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.suffix++
	id := g.prefix | lowbit(g.suffix, suffixLen)
	return id
}

func lowbit(x uint64, n uint) uint64 {
	return x & (math.MaxUint64 >> (64 - n))
}
