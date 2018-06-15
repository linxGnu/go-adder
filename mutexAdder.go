package goadder

import (
	"sync"
)

// MutexAdder mutex base adder
type MutexAdder struct {
	value int64
	lock  sync.Mutex
}

// NewMutexAdder create new MutexAdder
func NewMutexAdder() *MutexAdder {
	return &MutexAdder{}
}

// Add the given value
func (m *MutexAdder) Add(x int64) {
	m.lock.Lock()
	m.value += x
	m.lock.Unlock()
}

// Inc by 1
func (m *MutexAdder) Inc() {
	m.Add(1)
}

// Dec by 1
func (m *MutexAdder) Dec() {
	m.Add(-1)
}

// Sum return the current sum. The returned value is NOT an
// atomic snapshot; invocation in the absence of concurrent
// updates returns an accurate result, but concurrent updates that
// occur while the sum is being calculated might not be
// incorporated.
func (m *MutexAdder) Sum() (sum int64) {
	return m.value
}

// Reset variables maintaining the sum to zero. This method may be a useful alternative
// to creating a new adder, but is only effective if there are no concurrent updates.
// Because this method is intrinsically racy
func (m *MutexAdder) Reset() {
	m.value = 0
}

// SumAndReset equivalent in effect to sum followed by reset.
// This method may apply for example during quiescent
// points between multithreaded computations. If there are
// updates concurrent with this method, the returned value is
// guaranteed to be the final value occurring before
// the reset.
func (m *MutexAdder) SumAndReset() (sum int64) {
	sum = m.value
	m.value = 0
	return
}
