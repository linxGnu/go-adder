package longadder

import (
	"runtime"
	"sync"
	"testing"
)

var benchNumRoutine = 200
var benchDelta = 1000000
var benchDeltaSingleRoute = 10000000

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
	runtime.GOMAXPROCS(200)
}

func BenchmarkMutexAdderSingleRoutine(t *testing.B) {
	benchAdderSingleRoutine(mutexAdder1)
}

func BenchmarkAtomicAdderSingleRoutine(t *testing.B) {
	benchAdderSingleRoutine(atomicAdder1)
}

func BenchmarkRandomCellAdderSingleRoutine(t *testing.B) {
	benchAdderSingleRoutine(randomCellAdder1)
}

func BenchmarkJDKAdderSingleRoutine(t *testing.B) {
	benchAdderSingleRoutine(jdkAdder1)
}

func BenchmarkMutexAdderMultiRoutine(t *testing.B) {
	benchAdderMultiRoutine(mutexAdder2)
}

func BenchmarkAtomicAdderMultiRoutine(t *testing.B) {
	benchAdderMultiRoutine(atomicAdder2)
}

func BenchmarkRandomCellAdderMultiRoutine(t *testing.B) {
	benchAdderMultiRoutine(randomCellAdder2)
}

func BenchmarkJDKAdderMultiRoutine(t *testing.B) {
	benchAdderMultiRoutine(jdkAdder2)
}

func BenchmarkMutexAdderMultiRoutineMix(t *testing.B) {
	benchAdderMultiRoutineMix(mutexAdder3)
}

func BenchmarkAtomicAdderMultiRoutineMix(t *testing.B) {
	benchAdderMultiRoutineMix(atomicAdder3)
}
func BenchmarkRandomCellAdderMultiRoutineMix(t *testing.B) {
	benchAdderMultiRoutineMix(randomCellAdder3)
}

func BenchmarkJDKAdderMultiRoutineMix(t *testing.B) {
	benchAdderMultiRoutineMix(jdkAdder3)
}

func benchAdderSingleRoutine(adder LongAdder) {
	for i := 0; i < benchDeltaSingleRoute; i++ {
		jdkAdder1.Add(1)
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
