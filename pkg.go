/*
Package longadder contains a collection of thread-safe, concurrent data structures for reading and writing numeric int64-counter,
inspired by OpenJDK9 LongAdder.

Beside JDKAdder, ported version of OpenJDK9 LongAdder, package also provides other alternatives for various use cases.
*/
package longadder

// Type of LongAdder
type Type int

const (
	// JDKAdderType is type for JDK-based LongAdder
	JDKAdderType Type = iota
	// RandomCellAdderType is type for RandomCellAdder
	RandomCellAdderType
	// AtomicAdderType is type for AtomicAdder
	AtomicAdderType
	// MutexAdderType is type for MutexAdder
	MutexAdderType
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

// NewLongAdder create new LongAdder upon type
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

// LongBinaryOperator represents an operation upon two int64-valued operands and producing an
// int64-valued result
type LongBinaryOperator interface {
	Apply(left, right int64) int64
}
