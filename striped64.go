/*
 * DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.
 *
 * This code is free software; you can redistribute it and/or modify it
 * under the terms of the GNU General Public License version 2 only, as
 * published by the Free Software Foundation.  Oracle designates this
 * particular file as subject to the "Classpath" exception as provided
 * by Oracle in the LICENSE file that accompanied this code.
 *
 * This code is distributed in the hope that it will be useful, but WITHOUT
 * ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
 * FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License
 * version 2 for more details (a copy is included in the LICENSE file that
 * accompanied this code).
 *
 * You should have received a copy of the GNU General Public License version
 * 2 along with this work; if not, write to the Free Software Foundation,
 * Inc., 51 Franklin St, Fifth Floor, Boston, MA 02110-1301 USA.
 *
 * Please contact Oracle, 500 Oracle Parkway, Redwood Shores, CA 94065 USA
 * or visit www.oracle.com if you need additional information or have any
 * questions.
 */

/*
 * This file is available under and governed by the GNU General Public
 * License version 2 only, as published by the Free Software Foundation.
 * However, the following notice accompanied the original version of this
 * file:
 *
 * Written by Doug Lea with assistance from members of JCP JSR-166
 * Expert Group and released to the public domain, as explained at
 * http://creativecommons.org/publicdomain/zero/1.0/
 */

package goadder

import (
	"runtime"
	"sync/atomic"
	"time"
)

/**
 * A package-local class holding common representation and mechanics
 * for classes supporting dynamic striping on 64bit values. The class
 * extends Number so that concrete subclasses must publicly do so.
 */

/*
 * This class maintains a lazily-initialized table of atomically
 * updated variables, plus an extra "base" field. The table size
 * is a power of two. Indexing uses masked per-thread hash codes.
 * Nearly all declarations in this class are package-private,
 * accessed directly by subclasses.
 *
 * Table entries are of class Cell; a variant of AtomicLong padded
 * (via @Contended) to reduce cache contention. Padding is
 * overkill for most Atomics because they are usually irregularly
 * scattered in memory and thus don't interfere much with each
 * other. But Atomic objects residing in arrays will tend to be
 * placed adjacent to each other, and so will most often share
 * cache lines (with a huge negative performance impact) without
 * this precaution.
 *
 * In part because Cells are relatively large, we avoid creating
 * them until they are needed.  When there is no contention, all
 * updates are made to the base field.  Upon first contention (a
 * failed CAS on base update), the table is initialized to size 2.
 * The table size is doubled upon further contention until
 * reaching the nearest power of two greater than or equal to the
 * number of CPUS. Table slots remain empty (null) until they are
 * needed.
 *
 * A single spinlock ("cellsBusy") is used for initializing and
 * resizing the table, as well as populating slots with new Cells.
 * There is no need for a blocking lock; when the lock is not
 * available, threads try other slots (or the base).  During these
 * retries, there is increased contention and reduced locality,
 * which is still better than alternatives.
 *
 * The Thread probe fields maintained via ThreadLocalRandom serve
 * as per-thread hash codes. We let them remain uninitialized as
 * zero (if they come in this way) until they contend at slot
 * 0. They are then initialized to values that typically do not
 * often conflict with others.  Contention and/or table collisions
 * are indicated by failed CASes when performing an update
 * operation. Upon a collision, if the table size is less than
 * the capacity, it is doubled in size unless some other thread
 * holds the lock. If a hashed slot is empty, and lock is
 * available, a new Cell is created. Otherwise, if the slot
 * exists, a CAS is tried.  Retries proceed by "double hashing",
 * using a secondary hash (Marsaglia XorShift) to try to find a
 * free slot.
 *
 * The table size is capped because, when there are more threads
 * than CPUs, supposing that each thread were bound to a CPU,
 * there would exist a perfect hash function mapping threads to
 * slots that eliminates collisions. When we reach capacity, we
 * search for this mapping by randomly varying the hash codes of
 * colliding threads.  Because search is random, and collisions
 * only become known via CAS failures, convergence can be slow,
 * and because threads are typically not bound to CPUS forever,
 * may not occur at all. However, despite these limitations,
 * observed contention rates are typically low in these cases.
 *
 * It is possible for a Cell to become unused when threads that
 * once hashed to it terminate, as well as in the case where
 * doubling the table causes no thread to hash to it under
 * expanded mask.  We do not try to detect or remove such cells,
 * under the assumption that for long-running instances, observed
 * contention levels will recur, so the cells will eventually be
 * needed again; and for short-lived ones, it does not matter.
 */

var maxCells = runtime.NumCPU() << 2

func init() {
	if maxCells > (1 << 11) {
		maxCells = (1 << 11)
	}

	if maxCells < 128 {
		maxCells = 128
	}
}

// LongBinaryOperator represents an operation upon two int64-valued operands and producing an
// int64-valued result
type LongBinaryOperator interface {
	Apply(left, right int64) int64
}

type cell struct {
	val int64
}

func (c *cell) cas(old, new int64) bool {
	return atomic.CompareAndSwapInt64(&c.val, old, new)
}

// Striped64 porting from jdk
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
	var n int
	for {
		as = s.cells
		if n = len(as) - 1; n >= 0 {
			if a = as[probe&n]; a == nil {
				if s.cellsBusy == 0 {
					r := &cell{val: x}
					if s.cellsBusy == 0 && s.casCellsBusy() {
						if rs, m := s.cells, len(s.cells)-1; m >= 0 {
							if j := probe & m; rs[j] == nil {
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
		} else if as == nil {
			if s.cellsBusy == 0 && s.cells == nil && s.casCellsBusy() {
				if s.cells == nil { // Initialize table
					s.cells = make([]*cell, 2, 4)
					s.cells[0] = &cell{val: x}
					s.cellsBusy = 0
					break
				}
				s.cellsBusy = 0
			} else {
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
