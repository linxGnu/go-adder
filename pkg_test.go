package goadder

import (
	"sync"
	"testing"
)

var numRoutine = 9
var delta = 5237659

func testAdderNotRaceInc(t *testing.T, ty LongAdderType) {
	adder := NewLongAdder(ty)

	for i := 0; i < delta; i++ {
		adder.Inc()
	}

	tmp := int64(delta)
	if adder.Sum() != tmp || adder.SumAndReset() != tmp || adder.Sum() != 0 {
		t.Errorf("Adder(%d) logic is wrong", ty)
	}
}

func testAdderRaceInc(t *testing.T, ty LongAdderType) {
	adder := NewLongAdder(ty)

	var wg sync.WaitGroup
	for i := 0; i < numRoutine; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < delta; j++ {
				adder.Inc()
			}
			wg.Done()
		}()
	}
	wg.Wait()

	tmp := int64(delta) * int64(numRoutine)
	if adder.Sum() != tmp || adder.SumAndReset() != tmp || adder.Sum() != 0 {
		t.Errorf("Adder(%d) logic is wrong", ty)
	}
}

func testAdderNotRaceDec(t *testing.T, ty LongAdderType) {
	adder := NewLongAdder(ty)

	for i := 0; i < delta; i++ {
		adder.Dec()
	}

	if adder.Sum() != -int64(delta) {
		t.Errorf("Adder(%d) logic is wrong", ty)
	}

	adder.Reset()
	if adder.Sum() != 0 {
		t.Errorf("Adder(%d) logic is wrong", ty)
	}
}

func testAdderRaceDec(t *testing.T, ty LongAdderType) {
	adder := NewLongAdder(ty)

	var wg sync.WaitGroup
	for i := 0; i < numRoutine; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < delta; j++ {
				adder.Dec()
			}
			wg.Done()
		}()
	}
	wg.Wait()

	if adder.Sum() != -int64(delta)*int64(numRoutine) {
		t.Errorf("Adder(%d) logic is wrong", ty)
	}

	adder.Reset()
	if adder.Sum() != 0 {
		t.Errorf("Adder(%d) logic is wrong", ty)
	}
}

func testAdderNotRaceAdd(t *testing.T, ty LongAdderType) {
	adder := NewLongAdder(ty)

	for i := 0; i < delta; i++ {
		adder.Add(int64(i))
	}

	tmp := int64(delta)
	if adder.Sum() != tmp*(tmp-1)/2 {
		t.Errorf("Adder(%d) logic is wrong", ty)
	}
}

func testAdderRaceAdd(t *testing.T, ty LongAdderType) {
	adder := NewLongAdder(ty)

	var wg sync.WaitGroup
	for i := 0; i < numRoutine; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < delta; j++ {
				adder.Add(int64(j))
			}
			wg.Done()
		}()
	}
	wg.Wait()

	tmp := int64(delta)
	if adder.Sum() != (tmp*(tmp-1)/2)*int64(numRoutine) {
		t.Errorf("Adder(%d) logic is wrong", ty)
	}
}

// var adder *V1
// var atomicInt64 int64
// var adder2 *LongAdder

// func init() {
// 	adder = NewDefaultV1()
// 	adder2 = NewLongAdder()
// }

// // BenchmarkUint64Adder bench uint64 adder
// func BenchmarkUint64Adder(b *testing.B) {
// 	start := time.Now()
// 	var wg sync.WaitGroup
// 	for i := 0; i < numGR; i++ {
// 		wg.Add(1)
// 		go func() {
// 			for i := 0; i < numOP; i++ {
// 				adder.Add(1)
// 			}
// 			wg.Done()
// 		}()
// 	}
// 	wg.Wait()
// 	fmt.Println("A1", time.Now().Sub(start).Seconds())
// }

// func BenchmarkUint64AdderV2(b *testing.B) {
// 	start := time.Now()
// 	var wg sync.WaitGroup
// 	for i := 0; i < numGR; i++ {
// 		wg.Add(1)
// 		go func() {
// 			for i := 0; i < numOP; i++ {
// 				adder2.Add(1)
// 			}
// 			wg.Done()
// 		}()
// 	}
// 	wg.Wait()
// 	fmt.Println("A2", time.Now().Sub(start).Seconds())
// }

// // BenchmarkAtomicAdder bench atomic adder
// func BenchmarkAtomicAdder(b *testing.B) {
// 	start := time.Now()
// 	var wg sync.WaitGroup
// 	for i := 0; i < numGR; i++ {
// 		wg.Add(1)
// 		go func() {
// 			for i := 0; i < numOP; i++ {
// 				atomic.AddInt64(&atomicInt64, 1)
// 			}
// 			wg.Done()
// 		}()
// 	}
// 	wg.Wait()
// 	fmt.Println("B", time.Now().Sub(start).Seconds())
// }
