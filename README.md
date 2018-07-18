# longadder

Thread-safe, high performance, contention-aware `LongAdder` for Go, inspired by OpenJDK9 LongAdder.
Beside JDKAdder (ported from OpenJDK9), library includes other adders for various use.

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
* Faster than JDK LongAdder in multi-routine (race) benchmark but slower in case of single routine (no race).
* Consume ~1KB to store cells, which is often larger than JDK LongAdder which number of cells is dynamic.

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

* Hardware: MacBookPro14,3 (2.8 GHz Intel Core i7, 16 GB 2133 MHz LPDDR3)
* OS: Mac OS 10.13.5
* Source code: [pkg_bench_test.go](https://git.linecorp.com/LINE-DevOps/go-utils/blob/master/longadder/pkg_bench_test.go)

```
Number of routine: 200
Number of inc operation each routine: 1,000,000
Total ops: 200 * 1,000,000 = 200,000,000
```
```
goos: darwin
goarch: amd64
pkg: git.linecorp.com/LINE-DevOps/go-utils/longadder
BenchmarkMutexAdderSingleRoutine-200                    2000000000               0.05 ns/op
BenchmarkAtomicAdderSingleRoutine-200                   2000000000               0.05 ns/op
BenchmarkRandomCellAdderSingleRoutine-200               2000000000               0.05 ns/op
BenchmarkJDKAdderSingleRoutine-200                      2000000000               0.05 ns/op
BenchmarkMutexAdderMultiRoutine-200                            1        15624821667 ns/op
BenchmarkAtomicAdderMultiRoutine-200                           1        4447586221 ns/op
BenchmarkRandomCellAdderMultiRoutine-200                       1        1780231561 ns/op
BenchmarkJDKAdderMultiRoutine-200                              1        2092690423 ns/op
BenchmarkMutexAdderMultiRoutineMix-200                         1        14963214270 ns/op
BenchmarkAtomicAdderMultiRoutineMix-200                        1        4479809960 ns/op
BenchmarkRandomCellAdderMultiRoutineMix-200                    1        1951267127 ns/op
BenchmarkJDKAdderMultiRoutineMix-200                           1        2506134752 ns/op
```
