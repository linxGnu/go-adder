package goadder

import (
	"runtime"
	"sync"
	"testing"
)

var benchNumRoutine = 32
var benchDelta = 100000
var benchDeltaSingleRoute = 1000000

var atomicAdder1 = NewLongAdder(AtomicAdderType)
var mutexAdder1 = NewLongAdder(MutexAdderType)
var jdkAdder1 = NewLongAdder(JDKAdderType)
var randomCellAdder1 = NewLongAdder(RandomCellAdderType)

var atomicAdder2 = NewLongAdder(AtomicAdderType)
var mutexAdder2 = NewLongAdder(MutexAdderType)
var jdkAdder2 = NewLongAdder(JDKAdderType)
var randomCellAdder2 = NewLongAdder(RandomCellAdderType)

var atomicAdder3 = NewLongAdder(AtomicAdderType)
var mutexAdder3 = NewLongAdder(MutexAdderType)
var jdkAdder3 = NewLongAdder(JDKAdderType)
var randomCellAdder3 = NewLongAdder(RandomCellAdderType)

func init() {
	// set max procs to thread contention
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func BenchmarkMutexAdderSingleRoutine(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchAdderSingleRoutine(mutexAdder1)
	}
}

func BenchmarkAtomicAdderSingleRoutine(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchAdderSingleRoutine(atomicAdder1)
	}
}

func BenchmarkRandomCellAdderSingleRoutine(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchAdderSingleRoutine(randomCellAdder1)
	}
}

func BenchmarkJDKAdderSingleRoutine(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchAdderSingleRoutine(jdkAdder1)
	}
}

func BenchmarkMutexAdderMultiRoutine(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchAdderMultiRoutine(mutexAdder2)
	}
}

func BenchmarkAtomicAdderMultiRoutine(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchAdderMultiRoutine(atomicAdder2)
	}
}

func BenchmarkRandomCellAdderMultiRoutine(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchAdderMultiRoutine(randomCellAdder2)
	}
}

func BenchmarkJDKAdderMultiRoutine(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchAdderMultiRoutine(jdkAdder2)
	}
}

func BenchmarkMutexAdderMultiRoutineMix(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchAdderMultiRoutineMix(mutexAdder3)
	}
}

func BenchmarkAtomicAdderMultiRoutineMix(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchAdderMultiRoutineMix(atomicAdder3)
	}
}
func BenchmarkRandomCellAdderMultiRoutineMix(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchAdderMultiRoutineMix(randomCellAdder3)
	}
}

func BenchmarkJDKAdderMultiRoutineMix(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchAdderMultiRoutineMix(jdkAdder3)
	}
}

func benchAdderSingleRoutine(adder LongAdder) {
	for i := 0; i < benchDeltaSingleRoute; i++ {
		adder.Add(1)
	}
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

func benchAdderMultiRoutineMix(adder LongAdder) {
	var wg sync.WaitGroup
	for i := 0; i < benchNumRoutine; i++ {
		wg.Add(1)
		go func() {
			var sum int64
			for j := 0; j < benchDelta; j++ {
				adder.Add(1)
				if j%50 == 0 {
					sum += adder.Sum()
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
