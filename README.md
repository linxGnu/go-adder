# go-adder

[![Build Status](https://travis-ci.org/linxGnu/go-adder.svg?branch=master)](https://travis-ci.org/linxGnu/go-adder)
[![Go Report Card](https://goreportcard.com/badge/github.com/linxGnu/go-adder)](https://goreportcard.com/report/github.com/linxGnu/go-adder)
[![Coverage Status](https://coveralls.io/repos/github/linxGnu/go-adder/badge.svg?branch=master)](https://coveralls.io/github/linxGnu/go-adder?branch=master)
[![godoc](https://img.shields.io/badge/docs-GoDoc-green.svg)](https://godoc.org/github.com/linxGnu/go-adder)

Thread-safe, high performance, contention-aware LongAdder for Go. Beside LongAdder ported from OpenJDK9, library includes other adders for various use case.

# Benchmark

* Hardware: MacbookPro 2017 (MacBookPro14,3)

```
Number of routine: 200
Number of increment each routine: 1000000
```
```
BenchmarkAtomicAdderSingleRoutine-8             2000000000               0.04 ns/op
BenchmarkJDKAdderSingleRoutine-8                2000000000               0.05 ns/op
BenchmarkRandomCellAdderSingleRoutine-8         1000000000               0.27 ns/op
BenchmarkAtomicAdderMultiRoutine-8                     1        4454227441 ns/op
BenchmarkJDKAdderMultiRoutine-8                        1        2081959141 ns/op
BenchmarkRandomCellAdderMultiRoutine-8                 1        1845090920 ns/op
BenchmarkAtomicAdderMultiRoutineMix-8                  1        4524862037 ns/op
BenchmarkJDKAdderMultiRoutineMix-8                     1        2907979237 ns/op
BenchmarkRandomCellAdderMultiRoutineMix-8              1        2174805623 ns/op
```

# Usage

## JDK LongAdder

```go
import ga "github.com/linxGnu/go-adder"

func main() {
	adder := ga.NewLongAdder(ga.JDKAdderType)

    for i := 0; i < 100; i++ {
	   go func() {
          adder.Add(123)
       }()
    }
}
```

## RandomCellAdder

```
adder := ga.NewLongAdder(ga.RandomCellAdderType)
```

## AtomicAdder

```
adder := ga.NewLongAdder(ga.AtomicAdderType)
```