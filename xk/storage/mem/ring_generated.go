// Code generated by gommon from xk/storage/mem/ring_generated.go.tmpl DO NOT EDIT.

package mem

import "sync"

// IntRing is a simple cache without compression
type IntRing struct {
	numPartitions uint64
	// TODO: benchmark performance of using pointer and struct https://segment.com/blog/allocation-efficiency-in-high-performance-go-services/
	partitions []*IntPartition
}

// TODO: default value for n?
func NewIntRing(n int) *IntRing {
	r := &IntRing{
		numPartitions: uint64(n),
		partitions:    make([]*IntPartition, n),
	}
	for i := 0; i < n; i++ {
		r.partitions[i] = &IntPartition{
			store: make(map[uint64]*IntStore),
		}
	}
	return r
}

func (r *IntRing) getPartition(hash uint64) *IntPartition {
	return r.partitions[int(hash%r.numPartitions)]
}

type IntPartition struct {
	mu    sync.RWMutex
	store map[uint64]*IntStore
}

type IntStore struct {
	// TODO: metrics tags
	mu     sync.RWMutex
	times  []int64
	values []int64
	size   int
}

func NewIntStore(len int) *IntStore {
	return &IntStore{
		times:  make([]int64, 0, len),
		values: make([]int64, 0, len),
	}
}

func (p *IntPartition) WritePoints(hash uint64, times []int64, values []int64) {
	p.mu.RLock()
	s := p.store[hash]
	if s != nil {
		s.mu.Lock()
		p.mu.RUnlock()
		s.size += len(times)
		s.times = append(s.times, times...)
		s.values = append(s.values, values...)
		s.mu.Unlock()
		return
	}
	p.mu.RUnlock()
	p.mu.Lock()
	s = NewIntStore(len(times))
	s.size += len(times)
	s.times = append(s.times, times...)
	s.values = append(s.values, values...)
	p.store[hash] = s
	p.mu.Unlock()
	return
}

// DoubleRing is a simple cache without compression
type DoubleRing struct {
	numPartitions uint64
	// TODO: benchmark performance of using pointer and struct https://segment.com/blog/allocation-efficiency-in-high-performance-go-services/
	partitions []*DoublePartition
}

// TODO: default value for n?
func NewDoubleRing(n int) *DoubleRing {
	r := &DoubleRing{
		numPartitions: uint64(n),
		partitions:    make([]*DoublePartition, n),
	}
	for i := 0; i < n; i++ {
		r.partitions[i] = &DoublePartition{
			store: make(map[uint64]*DoubleStore),
		}
	}
	return r
}

func (r *DoubleRing) getPartition(hash uint64) *DoublePartition {
	return r.partitions[int(hash%r.numPartitions)]
}

type DoublePartition struct {
	mu    sync.RWMutex
	store map[uint64]*DoubleStore
}

type DoubleStore struct {
	// TODO: metrics tags
	mu     sync.RWMutex
	times  []int64
	values []float64
	size   int
}

func NewDoubleStore(len int) *DoubleStore {
	return &DoubleStore{
		times:  make([]int64, 0, len),
		values: make([]float64, 0, len),
	}
}

func (p *DoublePartition) WritePoints(hash uint64, times []int64, values []float64) {
	p.mu.RLock()
	s := p.store[hash]
	if s != nil {
		s.mu.Lock()
		p.mu.RUnlock()
		s.size += len(times)
		s.times = append(s.times, times...)
		s.values = append(s.values, values...)
		s.mu.Unlock()
		return
	}
	p.mu.RUnlock()
	p.mu.Lock()
	s = NewDoubleStore(len(times))
	s.size += len(times)
	s.times = append(s.times, times...)
	s.values = append(s.values, values...)
	p.store[hash] = s
	p.mu.Unlock()
	return
}
