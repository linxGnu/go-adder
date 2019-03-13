package goadder

import (
	"sync"
	"testing"
)

var atomicF64Adder1 = NewFloat64Adder(AtomicF64AdderType)
var jdkF64Adder1 = NewFloat64Adder(JDKF64AdderType)

var atomicF64Adder2 = NewFloat64Adder(AtomicF64AdderType)
var jdkF64Adder2 = NewFloat64Adder(JDKF64AdderType)

var atomicF64Adder3 = NewFloat64Adder(AtomicF64AdderType)
var jdkF64Adder3 = NewFloat64Adder(JDKF64AdderType)

func BenchmarkAtomicF64AdderSingleRoutine(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchF64AdderSingleRoutine(atomicF64Adder1)
	}
}

func BenchmarkJDKF64AdderSingleRoutine(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchF64AdderSingleRoutine(jdkF64Adder1)
	}
}

func BenchmarkAtomicF64AdderMultiRoutine(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchF64AdderMultiRoutine(atomicF64Adder2)
	}
}

func BenchmarkJDKF64AdderMultiRoutine(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchF64AdderMultiRoutine(jdkF64Adder2)
	}
}

func BenchmarkAtomicF64AdderMultiRoutineMix(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchF64AdderMultiRoutineMix(atomicF64Adder3)
	}
}

func BenchmarkJDKF64AdderMultiRoutineMix(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchF64AdderMultiRoutineMix(jdkF64Adder3)
	}
}

func benchF64AdderSingleRoutine(adder Float64Adder) {
	for i := 0; i < benchDeltaSingleRoute; i++ {
		adder.Add(1.1)
	}
}

func benchF64AdderMultiRoutine(adder Float64Adder) {
	var wg sync.WaitGroup
	for i := 0; i < benchNumRoutine; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < benchDelta; j++ {
				adder.Add(1.1)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func benchF64AdderMultiRoutineMix(adder Float64Adder) {
	var wg sync.WaitGroup
	for i := 0; i < benchNumRoutine; i++ {
		wg.Add(1)
		go func() {
			var sum float64
			for j := 0; j < benchDelta; j++ {
				adder.Add(1.1)
				if j%50 == 0 {
					sum += adder.Sum()
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
