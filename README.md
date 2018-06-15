# go-adder

[![Build Status](https://travis-ci.org/linxGnu/go-adder.svg?branch=master)](https://travis-ci.org/linxGnu/go-adder)
[![Go Report Card](https://goreportcard.com/badge/github.com/linxGnu/go-adder)](https://goreportcard.com/report/github.com/linxGnu/go-adder)
[![Coverage Status](https://coveralls.io/repos/github/linxGnu/go-adder/badge.svg?branch=master)](https://coveralls.io/github/linxGnu/go-adder?branch=master)
[![godoc](https://img.shields.io/badge/docs-GoDoc-green.svg)](https://godoc.org/github.com/linxGnu/go-adder)

Thread-safe, high performance, contention-aware `LongAdder` for Go. Beside JDK LongAdder (ported from OpenJDK9), library includes other adders for various use.

# Usage

## JDK LongAdder (recommended)

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
* Faster than JDK LongAdder in multi-routine (race) benchmark but much slower in case of single routine (no race).
* Consume 2KB to store cells, which is often larger than JDK LongAdder which number of cells is dynamic.

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

* Hardware: MacBookPro14,3

```
Number of routine: 200
Number of increment each routine: 1000000
```
```
BenchmarkMutexAdderSingleRoutine-8              2000000000               0.09 ns/op
BenchmarkAtomicAdderSingleRoutine-8             2000000000               0.04 ns/op
BenchmarkRandomCellAdderSingleRoutine-8         1000000000               0.27 ns/op
BenchmarkJDKAdderSingleRoutine-8                2000000000               0.05 ns/op
BenchmarkMutexAdderMultiRoutine-8                      1        20125355749 ns/op
BenchmarkAtomicAdderMultiRoutine-8                     1        4456265607 ns/op
BenchmarkRandomCellAdderMultiRoutine-8                 1        1824514151 ns/op
BenchmarkJDKAdderMultiRoutine-8                        1        2235518096 ns/op
BenchmarkMutexAdderMultiRoutineMix-8                   1        19153137432 ns/op
BenchmarkAtomicAdderMultiRoutineMix-8                  1        4516106413 ns/op
BenchmarkRandomCellAdderMultiRoutineMix-8              1        2170390082 ns/op
BenchmarkJDKAdderMultiRoutineMix-8                     1        3146307410 ns/op
```