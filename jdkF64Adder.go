package goadder

// JDKF64Adder is ported version of OpenJDK9 DoubleAdder.
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
// JDKF64Adder is high performance, non-blocking and safe for concurrent use.
type JDKF64Adder struct {
	*StripedF64
}

// NewJDKF64Adder create new JDKF64Adder
func NewJDKF64Adder() *JDKF64Adder {
	return &JDKF64Adder{&StripedF64{}}
}

// Add the given value
func (f *JDKF64Adder) Add(x float64) {
	_as, uncontended := f.cells.Load(), false
	if _as != nil {
		uncontended = true
	} else if b := f.base.load(); !f.base.cas(b, b+x) {
		uncontended = true
	}

	if uncontended {
		if _as == nil {
			f.accumulate(getRandomInt(), x, nil, true)
			return
		}

		as := _as.(cells)
		m := len(as) - 1
		if m < 0 {
			f.accumulate(getRandomInt(), x, nil, true)
			return
		}

		probe := getRandomInt() & m
		if _a := as[probe].Load(); _a == nil {
			f.accumulate(probe, x, nil, uncontended)
		} else {
			a := _a.(*cellf64)

			v := a.load()
			if uncontended = a.cas(v, v+x); !uncontended {
				f.accumulate(probe, x, nil, uncontended)
			}
		}
	}
}

// Inc by 1
func (f *JDKF64Adder) Inc() {
	f.Add(1)
}

// Dec by 1
func (f *JDKF64Adder) Dec() {
	f.Add(-1)
}

// Sum return the current sum. The returned value is NOT an
// atomic snapshot because of concurrent update.
func (f *JDKF64Adder) Sum() float64 {
	sum, _as := f.base.load(), f.cells.Load()
	if _as != nil {
		as := _as.(cells)
		var a interface{}
		for i := range as {
			if a = as[i].Load(); a != nil {
				sum += a.(*cellf64).load()
			}
		}
	}
	return sum
}

// Reset variables maintaining the sum to zero. This method may be a useful alternative
// to creating a new adder, but is only effective if there are no concurrent updates.
// Because this method is intrinsically racy.
func (f *JDKF64Adder) Reset() {
	f.store(0)
}

// SumAndReset equivalent in effect to sum followed by reset. Like the nature of Sum and Reset,
// this function is only effective if there are no concurrent updates.
func (f *JDKF64Adder) SumAndReset() (sum float64) {
	sum = f.Sum()
	f.Reset()
	return
}

// Store value. This function is only effective if there are no concurrent updates.
func (f *JDKF64Adder) Store(v float64) {
	f.store(v)
}

func (f *JDKF64Adder) store(v float64) {
	f.base.store(v)
	if _as := f.cells.Load(); _as != nil {
		cells := make(cells, len(_as.(cells)))
		for i := range cells {
			cells[i].Store(&cellf64{})
		}
		f.cells.Store(cells)
	}
}
