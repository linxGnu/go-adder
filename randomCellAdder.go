package goadder

import (
	"sync/atomic"
	"time"
)

const (
	randomCellSize      = 1 << 8
	randomCellSizeMinus = randomCellSize - 1
)

// RandomCellAdder long adder with simple strategy of preallocating atomic<cell>
// and select random cell to add. It's faster than JDKAdder but cause more memory in general case.
// RandomCellAdder consumes 2KB to store cells.
type RandomCellAdder struct {
	cells []int64
}

// NewRandomCellAdder create new RandomCellAdder
func NewRandomCellAdder() *RandomCellAdder {
	return &RandomCellAdder{
		cells: make([]int64, randomCellSize),
	}
}

// Add the given value
func (r *RandomCellAdder) Add(x int64) {
	atomic.AddInt64(&r.cells[time.Now().Nanosecond()&randomCellSizeMinus], x)
}

// Inc by 1
func (r *RandomCellAdder) Inc() {
	r.Add(1)
}

// Dec by 1
func (r *RandomCellAdder) Dec() {
	r.Add(-1)
}

// Sum return the current sum. The returned value is NOT an
// atomic snapshot; invocation in the absence of concurrent
// updates returns an accurate result, but concurrent updates that
// occur while the sum is being calculated might not be
// incorporated.
func (r *RandomCellAdder) Sum() (sum int64) {
	for _, v := range r.cells {
		sum += v
	}
	return
}
