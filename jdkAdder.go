package longadder

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
	*Striped64
}

// NewJDKAdder create new JDKAdder
func NewJDKAdder() *JDKAdder {
	return &JDKAdder{&Striped64{}}
}

// Add the given value
func (u *JDKAdder) Add(x int64) {
	as, uncontended := u.cells, false
	if as != nil {
		uncontended = true
	} else if b := u.base; !u.casBase(b, b+x) {
		uncontended = true
	}

	if uncontended {
		if as == nil {
			u.accumulate(getRandomInt(), x, nil, true)
			return
		}

		m := len(as) - 1
		if m < 0 {
			u.accumulate(getRandomInt(), x, nil, true)
			return
		}

		probe := getRandomInt() & m
		if a := as[probe]; a == nil {
			u.accumulate(probe, x, nil, uncontended)
		} else {
			v := a.val
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
	sum, as := u.base, u.cells
	for _, a := range as {
		if a != nil {
			sum += a.val
		}
	}
	return sum
}

// Reset variables maintaining the sum to zero. This method may be a useful alternative
// to creating a new adder, but is only effective if there are no concurrent updates.
// Because this method is intrinsically racy.
func (u *JDKAdder) Reset() {
	u.base = 0
	as := u.cells
	for _, a := range as {
		if a != nil {
			a.val = 0
		}
	}
}

// SumAndReset equivalent in effect to sum followed by reset. Like the nature of Sum and Reset,
// this function is only effective if there are no concurrent updates.
func (u *JDKAdder) SumAndReset() (sum int64) {
	sum = u.base
	u.base = 0
	as := u.cells
	for _, a := range as {
		if a != nil {
			sum += a.val
			a.val = 0
		}
	}
	return
}

// Store value. This function is only effective if there are no concurrent updates.
func (u *JDKAdder) Store(v int64) {
	as := u.cells
	for _, a := range as {
		if a != nil {
			a.val = 0
		}
	}
	u.base = v
}
