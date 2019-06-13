// Package goadder contains a collection of thread-safe, concurrent data structures for reading and writing numeric int64-counter,
// inspired by OpenJDK9 LongAdder.
//
// Beside JDKAdder, ported version of OpenJDK9 LongAdder, package also provides other alternatives for various use cases.
package goadder

// Type of LongAdder
type Type int

const (
	// JDKAdderType is type for JDK-based LongAdder
	JDKAdderType Type = iota
	// RandomCellAdderType is type for RandomCellAdder
	RandomCellAdderType
	// AtomicAdderType is type for atomic-based adder
	AtomicAdderType
	// MutexAdderType is type for MutexAdder
	MutexAdderType
	// JDKF64AdderType is type for JDK-based DoubleAdder
	JDKF64AdderType
	// AtomicF64AdderType is type for atomic-based float64 adder
	AtomicF64AdderType
)

// LongAdder interface
type LongAdder interface {
	Add(x int64)
	Inc()
	Dec()
	Sum() int64
	Reset()
	SumAndReset() int64
	Store(v int64)
}

// Float64Adder interface
type Float64Adder interface {
	Add(x float64)
	Inc()
	Dec()
	Sum() float64
	Reset()
	SumAndReset() float64
	Store(v float64)
}

// NewLongAdder create new long adder upon type
func NewLongAdder(t Type) LongAdder {
	switch t {
	case MutexAdderType:
		return NewMutexAdder()
	case AtomicAdderType:
		return NewAtomicAdder()
	case RandomCellAdderType:
		return NewRandomCellAdder()
	default:
		return NewJDKAdder()
	}
}

// NewFloat64Adder create new float64 adder upon type
func NewFloat64Adder(t Type) Float64Adder {
	switch t {
	case AtomicF64AdderType:
		return NewAtomicF64Adder()
	default:
		return NewJDKF64Adder()
	}
}

// LongBinaryOperator represents an operation upon two int64-valued operands and producing an
// int64-valued result
type LongBinaryOperator interface {
	Apply(left, right int64) int64
}

// FloatBinaryOperator represents an operation upon two float64-valued operands and producing an
// float64-valued result
type FloatBinaryOperator interface {
	Apply(left, right float64) float64
}
