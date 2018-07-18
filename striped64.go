package longadder

import (
	"runtime"
	"sync/atomic"
	"time"
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

type cell struct {
	val int64
}

func (c *cell) cas(old, new int64) bool {
	return atomic.CompareAndSwapInt64(&c.val, old, new)
}

/*
Striped64 is ported version of OpenJDK9 Striped64.
It maintains a lazily-initialized table of atomically
updated variables, plus an extra "base" field. The table size
is a power of two. Indexing uses masked per-routine hash codes.
Nearly all declarations in this class are package-private,
accessed directly by subclasses.

In part because Cells are relatively large, we avoid creating
them until they are needed. When there is no contention, all
updates are made to the base field. Upon first contention (a
failed CAS on base update), the table is initialized to size 2 and cap 4.
The table size is doubled upon further contention until
reaching the nearest power of two greater than or equal to the
number of CPUS. Table slots remain empty (null) until they are
needed.

A single spinlock ("cellsBusy") is used for initializing and
resizing the table, as well as populating slots with new Cells.
There is no need for a blocking lock; when the lock is not
available, routines try other slots (or the base). During these
retries, there is increased contention and reduced locality,
which is still better than alternatives.

The routine probe maintain by SystemTime nanoseconds instead of OpenJDK ThreadLocalRandom.
Contention and/or table collisions are indicated by failed CASes when performing an update
operation. Upon a collision, if the table size is less than
the capacity, it is doubled in size unless some other routine
holds the lock. If a hashed slot is empty, and lock is
available, a new Cell is created. Otherwise, if the slot
exists, a CAS is tried. Retries proceed with reproducing probe.

The table size is capped because, when there are more routines
than CPUs, supposing that each routine were bound to a CPU,
there would exist a perfect hash function mapping routines to
slots that eliminates collisions. When we reach capacity, we
search for this mapping by randomly varying the hash codes of
colliding routines. Because search is random, and collisions
only become known via CAS failures, convergence can be slow,
and because routines are typically not bound to CPUS forever,
may not occur at all. However, despite these limitations,
observed contention rates are typically low in these cases.

It is possible for a Cell to become unused when routines that
once hashed to it terminate, as well as in the case where
doubling the table causes no routine to hash to it under
expanded mask. We do not try to detect or remove such cells,
under the assumption that for long-running instances, observed
contention levels will recur, so the cells will eventually be
needed again; and for short-lived ones, it does not matter.
*/
type Striped64 struct {
	cells     []*cell
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
		probe = time.Now().Nanosecond()
		wasUncontended = true
	}

	collide := false
	var v, newV int64
	var as []*cell
	var a *cell
	var n, j int
	for {
		as = s.cells
		if as != nil {
			n = len(as) - 1
			if n < 0 {
				goto checkCells
			}

			if a = as[probe&n]; a == nil {
				if s.cellsBusy == 0 { // Try to attach new Cell
					r := &cell{val: x} // Optimistically create
					if s.cellsBusy == 0 && s.casCellsBusy() {
						if rs, m := s.cells, len(s.cells)-1; rs != nil && m >= 0 { // Recheck under lock
							if j = probe & m; rs[j] == nil {
								rs[j] = r
								s.cellsBusy = 0
								break
							}
						}
						s.cellsBusy = 0
						continue
					}
				}
				collide = false
			} else if !wasUncontended { // CAS already known to fail
				wasUncontended = true // Continue after rehash
			} else {
				probe &= n
				if v = a.val; fn == nil {
					newV = v + x
				} else {
					newV = fn.Apply(v, x)
				}
				if a.cas(v, newV) {
					break
				} else if n >= maxCells || &as[0] != &s.cells[0] { // At max size or stale
					collide = false
				} else if !collide {
					collide = true
				} else if s.cellsBusy == 0 && s.casCellsBusy() {
					if &as[0] == &s.cells[0] { // double size of cells
						if n = cap(as); len(as) < n {
							s.cells = s.cells[:n]
						} else {
							// slice is full, n == len(as) then we just x4 size for buffering
							// Note: this trick is different from jdk source code
							s.cells = make([]*cell, n<<1, n<<2)
							copy(s.cells, as)
						}
					}
					s.cellsBusy = 0
					collide = false
					continue
				}
			}

			probe = time.Now().Nanosecond()
			continue
		}

	checkCells:
		if as == nil {
			if s.cellsBusy == 0 && s.cells == nil && s.casCellsBusy() {
				if s.cells == nil { // Initialize table
					s.cells = make([]*cell, 2, 4)
					s.cells[probe&1] = &cell{val: x}
					s.cellsBusy = 0
					break
				}
				s.cellsBusy = 0
			} else { // Fall back on using base
				if v = s.base; fn == nil {
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
