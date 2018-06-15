package goadder

// LongAdderType type of long adder
type LongAdderType int

const (
	// JDKAdderType recommended long adder from JDK
	JDKAdderType LongAdderType = iota
	// RandomCellAdderType long adder with simple strategy of preallocating atomic cell
	// and select random cell to add.
	//
	// RandomCellAdder is faster than JDKAdder in multi routine race benchmark but much
	// slower in case of single routine (no race).
	//
	// RandomCellAdder consume 2KB for storing cells, which is often larger than JDKAdder
	// which number of cells is dynamic.
	RandomCellAdderType
	// AtomicAdderType simple atomic adder. Fastest at single routine but slowest at multi routine benchmark.
	AtomicAdderType
)

// LongAdder interface
type LongAdder interface {
	Add(x int64)
	Inc()
	Dec()
	Sum() int64
}

// NewLongAdder create new long adder base on type
func NewLongAdder(t LongAdderType) LongAdder {
	switch t {
	case AtomicAdderType:
		return NewAtomicAdder()
	case RandomCellAdderType:
		return NewRandomCellAdder()
	default:
		return NewJDKLongAdder()
	}
}
