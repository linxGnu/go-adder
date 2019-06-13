package goadder

import (
	"sync/atomic"
)

// JDKAdder is ported version of OpenJDK9 LongAdder.
//
// When multiple routines update a common sum that is used for purposes such as collecting statistics,
// not for fine-grained synchronization control, contention overhead could be a pain.
//
// JDKAdder is preferable to atomic, delivers significantly higher throughput under high contention,
// at the expense of higher space consumption, while keeping same characteristics under low contention.
//
// One or more variables, called Cells, together maintain an initially zero sum. When updates are contended across routines,
// the set of variables may grow dynamically to reduce contention. In other words, updates are distributed over Cells.
// The value is lazy, only aggregated (sum) over Cells when needed.
//
// JDKAdder is high performance, non-blocking and safe for concurrent use.
type JDKAdder struct {
	Striped64
}

// NewJDKAdder create new JDKAdder
func NewJDKAdder() *JDKAdder {
	return &JDKAdder{}
}

// Add the given value
func (u *JDKAdder) Add(x int64) {
	_as, uncontended := u.cells.Load(), false
	if _as != nil {
		uncontended = true
	} else if b := atomic.LoadInt64(&u.base); !u.casBase(b, b+x) {
		uncontended = true
	}

	if uncontended {
		if _as == nil {
			u.accumulate(getRandomInt(), x, nil, true)
			return
		}

		as := _as.(cells)
		m := len(as) - 1
		if m < 0 {
			u.accumulate(getRandomInt(), x, nil, true)
			return
		}

		probe := getRandomInt() & m
		if _a := as[probe].Load(); _a == nil {
			u.accumulate(probe, x, nil, uncontended)
		} else {
			a := _a.(*cell)

			v := atomic.LoadInt64(&a.val)
			if uncontended = a.cas(v, v+x); !uncontended {
				u.accumulate(probe, x, nil, uncontended)
			}
		}
	}
}

// Inc by 1
func (u *JDKAdder) Inc() {
	u.Add(1)
}

// Dec by 1
func (u *JDKAdder) Dec() {
	u.Add(-1)
}

// Sum return the current sum. The returned value is NOT an
// atomic snapshot because of concurrent update.
func (u *JDKAdder) Sum() int64 {
	sum, _as := atomic.LoadInt64(&u.base), u.cells.Load()
	if _as != nil {
		as := _as.(cells)
		var a interface{}
		for i := range as {
			if a = as[i].Load(); a != nil {
				sum += atomic.LoadInt64(&a.(*cell).val)
			}
		}
	}
	return sum
}

// Reset variables maintaining the sum to zero. This method may be a useful alternative
// to creating a new adder, but is only effective if there are no concurrent updates.
// Because this method is intrinsically racy.
func (u *JDKAdder) Reset() {
	u.store(0)
}

// SumAndReset equivalent in effect to sum followed by reset. Like the nature of Sum and Reset,
// this function is only effective if there are no concurrent updates.
func (u *JDKAdder) SumAndReset() (sum int64) {
	sum = u.Sum()
	u.Reset()
	return
}

// Store value. This function is only effective if there are no concurrent updates.
func (u *JDKAdder) Store(v int64) {
	u.store(v)
}

func (u *JDKAdder) store(v int64) {
	atomic.StoreInt64(&u.base, v)
	if _as := u.cells.Load(); _as != nil {
		cells := make(cells, len(_as.(cells)))
		for i := range cells {
			cells[i].Store(&cell{})
		}
		u.cells.Store(cells)
	}
}
