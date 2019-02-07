package goadder

import (
	"math"
	"sync/atomic"
)

// AtomicF64Adder is simple atomic-based adder. Fastest at single routine but slow at multi routine when high-contention happens.
type AtomicF64Adder struct {
	value uint64
}

// NewAtomicF64Adder create new AtomicF64Adder
func NewAtomicF64Adder() *AtomicF64Adder {
	return &AtomicF64Adder{}
}

// Add the given value
func (a *AtomicF64Adder) Add(v float64) {
	for {
		old := a.Sum()
		new := old + v
		if atomic.CompareAndSwapUint64(&a.value, math.Float64bits(old), math.Float64bits(new)) {
			return
		}
	}
}

// Inc by 1
func (a *AtomicF64Adder) Inc() {
	a.Add(1)
}

// Dec by 1
func (a *AtomicF64Adder) Dec() {
	a.Add(-1)
}

// Sum return the current sum. The returned value is NOT an
// atomic snapshot because of concurrent update.
func (a *AtomicF64Adder) Sum() float64 {
	return math.Float64frombits(atomic.LoadUint64(&a.value))
}

// Reset variables maintaining the sum to zero. This method may be a useful alternative
// to creating a new adder, but is only effective if there are no concurrent updates.
// Because this method is intrinsically racy.
func (a *AtomicF64Adder) Reset() {
	atomic.StoreUint64(&a.value, 0)
}

// SumAndReset equivalent in effect to sum followed by reset. Like the nature of Sum and Reset,
// this function is only effective if there are no concurrent updates.
func (a *AtomicF64Adder) SumAndReset() (sum float64) {
	sum = a.Sum()
	a.Reset()
	return
}

// Store value. This function is only effective if there are no concurrent updates.
func (a *AtomicF64Adder) Store(v float64) {
	atomic.StoreUint64(&a.value, math.Float64bits(v))
}
