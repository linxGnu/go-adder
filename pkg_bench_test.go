package goadder

import (
	"sync"
	"testing"
)

var benchNumRoutine = 200
var benchDelta = 1000000
var benchDeltaSingleRoute = 10000000

var atomicAdder1 = NewLongAdder(AtomicAdderType)
var jdkAdder1 = NewLongAdder(JDKAdderType)
var randomCellAdder1 = NewLongAdder(RandomCellAdderType)

var atomicAdder2 = NewLongAdder(AtomicAdderType)
var jdkAdder2 = NewLongAdder(JDKAdderType)
var randomCellAdder2 = NewLongAdder(RandomCellAdderType)

func BenchmarkAtomicAdderSingleRoutine(t *testing.B) {
	for i := 0; i < benchDeltaSingleRoute; i++ {
		atomicAdder1.Add(1)
	}
}

func BenchmarkJDKAdderSingleRoutine(t *testing.B) {
	for i := 0; i < benchDeltaSingleRoute; i++ {
		jdkAdder1.Add(1)
	}
}

func BenchmarkRandomCellAdderSingleRoutine(t *testing.B) {
	for i := 0; i < benchDeltaSingleRoute; i++ {
		randomCellAdder1.Add(1)
	}
}

func BenchmarkAtomicAdderMultiRoutine(t *testing.B) {
	benchAdderMultiRoutine(atomicAdder2)
}

func BenchmarkJDKAdderMultiRoutine(t *testing.B) {
	benchAdderMultiRoutine(jdkAdder2)
}

func BenchmarkRandomCellAdderMultiRoutine(t *testing.B) {
	benchAdderMultiRoutine(randomCellAdder2)
}

func benchAdderMultiRoutine(adder LongAdder) {
	var wg sync.WaitGroup
	for i := 0; i < benchNumRoutine; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < benchDelta; j++ {
				adder.Add(1)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
