package goadder

import (
	"sync/atomic"
)

// AtomicAdder is simple atomic-based adder. Fastest at single routine but slow at multi routine when high-contention happens.
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
// atomic snapshot because of concurrent update.
func (a *AtomicAdder) Sum() int64 {
	return atomic.LoadInt64(&a.value)
}

// Reset variables maintaining the sum to zero. This method may be a useful alternative
// to creating a new adder, but is only effective if there are no concurrent updates.
// Because this method is intrinsically racy.
func (a *AtomicAdder) Reset() {
	atomic.StoreInt64(&a.value, 0)
}

// SumAndReset equivalent in effect to sum followed by reset. Like the nature of Sum and Reset,
// this function is only effective if there are no concurrent updates.
func (a *AtomicAdder) SumAndReset() (sum int64) {
	sum = atomic.LoadInt64(&a.value)
	atomic.StoreInt64(&a.value, 0)
	return
}

// Store value. This function is only effective if there are no concurrent updates.
func (a *AtomicAdder) Store(v int64) {
	atomic.StoreInt64(&a.value, v)
}
