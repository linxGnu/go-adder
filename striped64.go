package goadder

import (
	"runtime"
	"sync/atomic"
)

var maxCells = runtime.NumCPU() << 2

func init() {
	if maxCells > (1 << 11) {
		maxCells = (1 << 11)
	}

	if maxCells < 64 {
		maxCells = 64
	}
}

type cells []atomic.Value

type cell struct {
	_   [7]uint64
	val int64
	_   [7]uint64
}

func (c *cell) cas(old, new int64) bool {
	return atomic.CompareAndSwapInt64(&c.val, old, new)
}

// Striped64 is ported version of OpenJDK9 Striped64.
// It maintains a lazily-initialized table of atomically
// updated variables, plus an extra "base" field. The table size
// is a power of two. Indexing uses masked per-routine hash codes.
// Nearly all declarations in this class are package-private,
// accessed directly by subclasses.

// In part because Cells are relatively large, we avoid creating
// them until they are needed. When there is no contention, all
// updates are made to the base field. Upon first contention (a
// failed CAS on base update), the table is initialized to size 2 and cap 4.
// The table size is doubled upon further contention until
// reaching the nearest power of two greater than or equal to the
// number of CPUS. Table slots remain empty (null) until they are
// needed.

// A single spinlock ("cellsBusy") is used for initializing and
// resizing the table, as well as populating slots with new Cells.
// There is no need for a blocking lock; when the lock is not
// available, routines try other slots (or the base). During these
// retries, there is increased contention and reduced locality,
// which is still better than alternatives.

// The routine probe maintain by SystemTime nanoseconds instead of OpenJDK ThreadLocalRandom.
// Contention and/or table collisions are indicated by failed CASes when performing an update
// operation. Upon a collision, if the table size is less than
// the capacity, it is doubled in size unless some other routine
// holds the lock. If a hashed slot is empty, and lock is
// available, a new Cell is created. Otherwise, if the slot
// exists, a CAS is tried. Retries proceed with reproducing probe.

// The table size is capped because, when there are more routines
// than CPUs, supposing that each routine were bound to a CPU,
// there would exist a perfect hash function mapping routines to
// slots that eliminates collisions. When we reach capacity, we
// search for this mapping by randomly varying the hash codes of
// colliding routines. Because search is random, and collisions
// only become known via CAS failures, convergence can be slow,
// and because routines are typically not bound to CPUS forever,
// may not occur at all. However, despite these limitations,
// observed contention rates are typically low in these cases.

// It is possible for a Cell to become unused when routines that
// once hashed to it terminate, as well as in the case where
// doubling the table causes no routine to hash to it under
// expanded mask. We do not try to detect or remove such cells,
// under the assumption that for long-running instances, observed
// contention levels will recur, so the cells will eventually be
// needed again; and for short-lived ones, it does not matter.
type Striped64 struct {
	cells     atomic.Value
	cellsBusy int32
	base      int64
}

func (s *Striped64) casBase(old, new int64) bool {
	return atomic.CompareAndSwapInt64(&s.base, old, new)
}

func (s *Striped64) casCellsBusy() bool {
	return atomic.CompareAndSwapInt32(&s.cellsBusy, 0, 1)
}

func (s *Striped64) accumulate(probe int, x int64, fn LongBinaryOperator, wasUncontended bool) {
	if probe == 0 {
		probe = getRandomInt()
		wasUncontended = true
	}

	collide := false
	var v, newV int64
	var as, rs cells
	var a, r *cell
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
				a = _a.(*cell)
			} else {
				a = nil
			}

			if a == nil {
				if atomic.LoadInt32(&s.cellsBusy) == 0 { // Try to attach new Cell
					r = &cell{val: x} // Optimistically create
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
				if v = atomic.LoadInt64(&a.val); fn == nil {
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
					rs[probe&1].Store(&cell{val: x})
					s.cells.Store(rs)
					atomic.StoreInt32(&s.cellsBusy, 0)
					break
				}
				atomic.StoreInt32(&s.cellsBusy, 0)
			} else { // Fall back on using base
				if v = atomic.LoadInt64(&s.base); fn == nil {
					newV = v + x
				} else {
					newV = fn.Apply(v, x)
				}
				if s.casBase(v, newV) {
					break
				}
			}
		}
	}
}
