# longadder

[![](https://github.com/linxGnu/go-adder/workflows/CI/badge.svg)]()
[![Go Report Card](https://goreportcard.com/badge/github.com/linxGnu/go-adder)](https://goreportcard.com/report/github.com/linxGnu/go-adder)
[![Coverage Status](https://coveralls.io/repos/github/linxGnu/go-adder/badge.svg?branch=master)](https://coveralls.io/github/linxGnu/go-adder?branch=master)
[![godoc](https://img.shields.io/badge/docs-GoDoc-green.svg)](https://godoc.org/github.com/linxGnu/go-adder)

Thread-safe, high performance, contention-aware `LongAdder` and `DoubleAdder` for Go, inspired by OpenJDK9.
Beside JDK-based `LongAdder` and `DoubleAdder`, library includes other adders for various use.

# Usage

## JDKAdder (recommended)

```go
package main

import (
	"fmt"
	"time"

	ga "github.com/linxGnu/go-adder"
)

func main() {
	adder := ga.NewLongAdder(ga.JDKAdderType)

	for i := 0; i < 100; i++ {
		go func() {
			adder.Add(123)
		}()
	}

	time.Sleep(3 * time.Second)

	// get total added value
	fmt.Println(adder.Sum()) 
}
```

## RandomCellAdder

* A `LongAdder` with simple strategy of preallocating atomic cell and select random cell for update.
* Slower than JDK LongAdder but 1.5-2x faster than atomic adder on contention.
* Consume ~1KB to store cells.

```
adder := ga.NewLongAdder(ga.RandomCellAdderType)
```

## AtomicAdder

* A `LongAdder` based on atomic variable. All routines share this variable.

```
adder := ga.NewLongAdder(ga.AtomicAdderType)
```

## MutexAdder

* A `LongAdder` based on mutex. All routines share same value and mutex.

```
adder := ga.NewLongAdder(ga.MutexAdderType)
```

# Benchmark

* System:         Dell PowerEdge R640
* CPU:            2 x Xeon Silver 4114 2.20GHz (40/40cores)
* Memory:         64GB 2400MHz DDR4
* OS:             CentOS 7.5, 64-bit
* Source code: [pkg_bench_test.go](https://github.com/linxGnu/go-adder/blob/master/pkg_bench_test.go)
* Go version: 1.12

```scala
goos: linux
goarch: amd64
BenchmarkAtomicF64AdderSingleRoutine-8               100          11558784 ns/op               0 B/op          0 allocs/op
BenchmarkJDKF64AdderSingleRoutine-8                  100          12663869 ns/op               0 B/op          0 allocs/op
BenchmarkAtomicF64AdderMultiRoutine-8                  3         414693263 ns/op            6672 B/op         18 allocs/op
BenchmarkJDKF64AdderMultiRoutine-8                    30          39951114 ns/op            1178 B/op          4 allocs/op
BenchmarkAtomicF64AdderMultiRoutineMix-8               3         410603215 ns/op              16 B/op          1 allocs/op
BenchmarkJDKF64AdderMultiRoutineMix-8                 30          49619434 ns/op              56 B/op          1 allocs/op
BenchmarkMutexAdderSingleRoutine-40                   30          42931025 ns/op               0 B/op          0 allocs/op
BenchmarkAtomicAdderSingleRoutine-40                 100          10022343 ns/op               0 B/op          0 allocs/op
BenchmarkRandomCellAdderSingleRoutine-40              50          38920149 ns/op             108 B/op          0 allocs/op
BenchmarkJDKAdderSingleRoutine-40                    100          14030302 ns/op               0 B/op          0 allocs/op
BenchmarkMutexAdderMultiRoutine-40                     2         576540605 ns/op            1456 B/op         16 allocs/op
BenchmarkAtomicAdderMultiRoutine-40                   20          88861041 ns/op             419 B/op          2 allocs/op
BenchmarkRandomCellAdderMultiRoutine-40               30          45493866 ns/op             239 B/op          3 allocs/op
BenchmarkJDKAdderMultiRoutine-40                      50          25724032 ns/op             140 B/op          2 allocs/op
BenchmarkMutexAdderMultiRoutineMix-40                  2         581924480 ns/op            1120 B/op         12 allocs/op
BenchmarkAtomicAdderMultiRoutineMix-40                20          93733789 ns/op              16 B/op          1 allocs/op
BenchmarkRandomCellAdderMultiRoutineMix-40            20          62700287 ns/op             331 B/op          4 allocs/op
BenchmarkJDKAdderMultiRoutineMix-40                   30          45089173 ns/op             230 B/op          3 allocs/op
```
