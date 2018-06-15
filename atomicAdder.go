package goadder

import (
	"sync/atomic"
)

// AtomicAdder simple atomic adder
type AtomicAdder struct {
	value int64
}

// NewAtomicAdder create new AtomicAdder
func NewAtomicAdder() *AtomicAdder {
	return &AtomicAdder{}
}

// Add the given value
func (a *AtomicAdder) Add(x int64) {
	atomic.AddInt64(&a.value, x)
}

// Inc by 1
func (a *AtomicAdder) Inc() {
	a.Add(1)
}

// Dec by 1
func (a *AtomicAdder) Dec() {
	a.Add(-1)
}

// Sum return the current sum. The returned value is NOT an
// atomic snapshot; invocation in the absence of concurrent
// updates returns an accurate result, but concurrent updates that
// occur while the sum is being calculated might not be
// incorporated.
func (a *AtomicAdder) Sum() int64 {
	return a.value
}

// Reset variables maintaining the sum to zero. This method may be a useful alternative
// to creating a new adder, but is only effective if there are no concurrent updates.
// Because this method is intrinsically racy
func (a *AtomicAdder) Reset() {
	a.value = 0
}

// SumAndReset equivalent in effect to sum followed by reset.
// This method may apply for example during quiescent
// points between multithreaded computations. If there are
// updates concurrent with this method, the returned value is
// guaranteed to be the final value occurring before
// the reset.
func (a *AtomicAdder) SumAndReset() (sum int64) {
	sum = a.Sum()
	a.value = 0
	return
}
