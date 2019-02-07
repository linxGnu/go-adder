package goadder

import (
	"math"
	"sync/atomic"
)

type cellf64 struct {
	_   [7]uint64
	val uint64
	_   [7]uint64
}

func (c *cellf64) load() float64 {
	return math.Float64frombits(atomic.LoadUint64(&c.val))
}

func (c *cellf64) store(v float64) {
	atomic.StoreUint64(&c.val, math.Float64bits(v))
}

func (c *cellf64) cas(old, new float64) bool {
	return atomic.CompareAndSwapUint64(&c.val, math.Float64bits(old), math.Float64bits(new))
}

func (c *cellf64) add(v float64) float64 {
	for {
		old := c.load()
		new := old + v
		if c.cas(old, new) {
			return new
		}
	}
}

// StripedF64 same as Striped64 but for float64
type StripedF64 struct {
	cells     atomic.Value
	cellsBusy int32
	base      cellf64
}

func (s *StripedF64) casCellsBusy() bool {
	return atomic.CompareAndSwapInt32(&s.cellsBusy, 0, 1)
}

func (s *StripedF64) accumulate(probe int, x float64, fn FloatBinaryOperator, wasUncontended bool) {
	if probe == 0 {
		probe = getRandomInt()
		wasUncontended = true
	}

	collide := false
	var v, newV float64
	var as, rs cells
	var a, r *cellf64
	var m, n, j int

	var _a, _as interface{}

	for {
		_as = s.cells.Load()
		if _as != nil {
			as = _as.(cells)

			n = len(as) - 1
			if n < 0 {
				goto checkCells
			}

			if _a = as[probe&n].Load(); _a != nil {
				a = _a.(*cellf64)
			} else {
				a = nil
			}

			if a == nil {
				if atomic.LoadInt32(&s.cellsBusy) == 0 { // Try to attach new Cell
					r = &cellf64{} // Optimistically create
					r.store(x)

					if atomic.LoadInt32(&s.cellsBusy) == 0 && s.casCellsBusy() {
						rs = s.cells.Load().(cells)
						if m = len(rs) - 1; rs != nil && m >= 0 { // Recheck under lock
							if j = probe & m; rs[j].Load() == nil {
								rs[j].Store(r)
								atomic.StoreInt32(&s.cellsBusy, 0)
								break
							}
						}
						atomic.StoreInt32(&s.cellsBusy, 0)
						continue
					}
				}
				collide = false
			} else if !wasUncontended { // CAS already known to fail
				wasUncontended = true // Continue after rehash
			} else {
				probe &= n
				if v = a.load(); fn == nil {
					newV = v + x
				} else {
					newV = fn.Apply(v, x)
				}
				if a.cas(v, newV) {
					break
				} else if n >= maxCells || &as[0] != &s.cells.Load().(cells)[0] { // At max size or stale
					collide = false
				} else if !collide {
					collide = true
				} else if atomic.LoadInt32(&s.cellsBusy) == 0 && s.casCellsBusy() {
					rs = s.cells.Load().(cells)
					if &as[0] == &rs[0] { // double size of cells
						if n = cap(as); len(as) < n {
							s.cells.Store(rs[:n])
						} else {
							// slice is full, n == len(as) then we just x4 size for buffering
							// Note: this trick is different from jdk source code
							rs = make(cells, n<<1, n<<2)
							copy(rs, as)
							s.cells.Store(rs)
						}
					}
					atomic.StoreInt32(&s.cellsBusy, 0)
					collide = false
					continue
				}
			}

			probe ^= probe << 13 // xorshift
			probe ^= probe >> 17
			probe ^= probe << 5
			continue
		}

	checkCells:
		if _as == nil {
			if atomic.LoadInt32(&s.cellsBusy) == 0 && s.cells.Load() == nil && s.casCellsBusy() {
				if s.cells.Load() == nil { // Initialize table
					rs = make(cells, 2, 4)

					r = &cellf64{}
					r.store(x)

					rs[probe&1].Store(r)
					s.cells.Store(rs)
					atomic.StoreInt32(&s.cellsBusy, 0)
					break
				}
				atomic.StoreInt32(&s.cellsBusy, 0)
			} else { // Fall back on using base
				if v = s.base.load(); fn == nil {
					newV = v + x
				} else {
					newV = fn.Apply(v, x)
				}
				if s.base.cas(v, newV) {
					break
				}
			}
		}
	}
}
