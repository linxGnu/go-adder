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
	"time"
)

/**
 * One or more variables that together maintain an initially zero
 * {@code long} sum.  When updates (method {@link #add}) are contended
 * across threads, the set of variables may grow dynamically to reduce
 * contention. Method {@link #sum} (or, equivalently, {@link
 * #longValue}) returns the current total combined across the
 * variables maintaining the sum.
 *
 * <p>This class is usually preferable to {@link AtomicLong} when
 * multiple threads update a common sum that is used for purposes such
 * as collecting statistics, not for fine-grained synchronization
 * control.  Under low update contention, the two classes have similar
 * characteristics. But under high contention, expected throughput of
 * this class is significantly higher, at the expense of higher space
 * consumption.
 *
 * <p>LongAdders can be used with a {@link
 * java.util.concurrent.ConcurrentHashMap} to maintain a scalable
 * frequency map (a form of histogram or multiset). For example, to
 * add a count to a {@code ConcurrentHashMap<String,LongAdder> freqs},
 * initializing if not already present, you can use {@code
 * freqs.computeIfAbsent(key, k -> new LongAdder()).increment();}
 *
 * <p>This class extends {@link Number}, but does <em>not</em> define
 * methods such as {@code equals}, {@code hashCode} and {@code
 * compareTo} because instances are expected to be mutated, and so are
 * not useful as collection keys.
 *
 * @since 1.8
 * @author Doug Lea
 */

// JDKLongAdder ported from jdk
type JDKLongAdder struct {
	*Striped64
}

// NewJDKLongAdder create new JDKLongAdder
func NewJDKLongAdder() *JDKLongAdder {
	return &JDKLongAdder{&Striped64{}}
}

// Add the given value
func (u *JDKLongAdder) Add(x int64) {
	as, uncontended := u.cells, false
	if as != nil {
		uncontended = true
	} else if b := u.base; !u.casBase(b, b+x) {
		uncontended = true
	}

	if uncontended {
		m := len(as) - 1
		if m < 0 {
			u.accumulate(time.Now().Nanosecond(), x, nil, true)
			return
		}

		probe := time.Now().Nanosecond() & m
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
func (u *JDKLongAdder) Inc() {
	u.Add(1)
}

// Dec by 1
func (u *JDKLongAdder) Dec() {
	u.Add(-1)
}

// Sum return the current sum. The returned value is NOT an
// atomic snapshot; invocation in the absence of concurrent
// updates returns an accurate result, but concurrent updates that
// occur while the sum is being calculated might not be
// incorporated.
func (u *JDKLongAdder) Sum() int64 {
	sum, cells := u.base, u.cells
	for _, v := range cells {
		if v != nil {
			sum += v.val
		}
	}
	return sum
}
