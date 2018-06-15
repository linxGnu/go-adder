# go-adder

[![Build Status](https://travis-ci.org/linxGnu/go-adder.svg?branch=master)](https://travis-ci.org/linxGnu/go-adder)
[![Go Report Card](https://goreportcard.com/badge/github.com/linxGnu/go-adder)](https://goreportcard.com/report/github.com/linxGnu/go-adder)
[![Coverage Status](https://coveralls.io/repos/github/linxGnu/go-adder/badge.svg?branch=master)](https://coveralls.io/github/linxGnu/go-adder?branch=master)
[![godoc](https://img.shields.io/badge/docs-GoDoc-green.svg)](https://godoc.org/github.com/linxGnu/go-adder)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/jmoiron/sqlx/master/LICENSE)

Thread-safe, high performance, contention-aware LongAdder for Go. Beside LongAdder ported from OpenJDK9, library includes other adders for various use case.

# Benchmark

* MacbookPro 2017 (MacBookPro14,3)

```
BenchmarkAtomicAdderSingleRoutine-8             2000000000               0.03 ns/op
BenchmarkJDKAdderSingleRoutine-8                2000000000               0.05 ns/op
BenchmarkRandomCellAdderSingleRoutine-8         2000000000               0.13 ns/op
BenchmarkAtomicAdderMultiRoutine-8                     1        4428545808 ns/op
BenchmarkJDKAdderMultiRoutine-8                        1        2181652003 ns/op
BenchmarkRandomCellAdderMultiRoutine-8                 1        1834870874 ns/op
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