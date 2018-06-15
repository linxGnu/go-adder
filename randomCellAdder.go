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

// Reset variables maintaining the sum to zero. This method may be a useful alternative
// to creating a new adder, but is only effective if there are no concurrent updates.
// Because this method is intrinsically racy
func (r *RandomCellAdder) Reset() {
	for i := range r.cells {
		r.cells[i] = 0
	}
}

// SumAndReset equivalent in effect to sum followed by reset.
// This method may apply for example during quiescent
// points between multithreaded computations. If there are
// updates concurrent with this method, the returned value is
// guaranteed to be the final value occurring before
// the reset.
func (r *RandomCellAdder) SumAndReset() (sum int64) {
	for i := range r.cells {
		sum += r.cells[i]
		r.cells[i] = 0
	}
	return
}
