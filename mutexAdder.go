package longadder

import (
	"sync"
)

// MutexAdder is mutex-based LongAdder. Slowest compared to other alternatives.
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
// atomic snapshot because of concurrent update.
func (m *MutexAdder) Sum() (sum int64) {
	return m.value
}

// Reset variables maintaining the sum to zero. This method may be a useful alternative
// to creating a new adder, but is only effective if there are no concurrent updates.
// Because this method is intrinsically racy.
func (m *MutexAdder) Reset() {
	m.value = 0
}

// SumAndReset equivalent in effect to sum followed by reset. Like the nature of Sum and Reset,
// this function is only effective if there are no concurrent updates.
func (m *MutexAdder) SumAndReset() (sum int64) {
	sum = m.value
	m.value = 0
	return
}

// Store value. This function is only effective if there are no concurrent updates.
func (m *MutexAdder) Store(v int64) {
	m.lock.Lock()
	m.value = v
	m.lock.Unlock()
}
