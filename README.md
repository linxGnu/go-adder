# longadder

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
* Slower than JDK LongAdder but 1.5-2x faster than atomic adder.
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

```
BenchmarkAtomicF64AdderSingleRoutine-201                2000000000               0.08 ns/op
BenchmarkJDKF64AdderSingleRoutine-201                   2000000000               0.08 ns/op
BenchmarkAtomicF64AdderMultiRoutine-201                        1        1394931835 ns/op
BenchmarkJDKF64AdderMultiRoutine-201                           1        1383563619 ns/op
BenchmarkAtomicF64AdderMultiRoutineMix-201                     1        2389037607 ns/op
BenchmarkJDKF64AdderMultiRoutineMix-201                        1        2187544688 ns/op
BenchmarkMutexAdderSingleRoutine-201                    2000000000               0.22 ns/op
BenchmarkAtomicAdderSingleRoutine-201                   2000000000               0.05 ns/op
BenchmarkRandomCellAdderSingleRoutine-201               2000000000               0.19 ns/op
BenchmarkJDKAdderSingleRoutine-201                      2000000000               0.07 ns/op
BenchmarkMutexAdderMultiRoutine-201                            1        27553411399 ns/op
BenchmarkAtomicAdderMultiRoutine-201                           1        5661739378 ns/op
BenchmarkRandomCellAdderMultiRoutine-201                       1        2784614208 ns/op
BenchmarkJDKAdderMultiRoutine-201                              1        1242928566 ns/op
BenchmarkMutexAdderMultiRoutineMix-201                         1        28100487108 ns/op
BenchmarkAtomicAdderMultiRoutineMix-201                        1        5627914921 ns/op
BenchmarkRandomCellAdderMultiRoutineMix-201                    1        3808765307 ns/op
BenchmarkJDKAdderMultiRoutineMix-201                           1        2165949485 ns/op
```
