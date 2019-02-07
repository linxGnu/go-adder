package goadder

import (
	"testing"
)

func TestAtomicF64AdderNotRaceInc(t *testing.T) {
	testF64AdderNotRaceInc(t, AtomicF64AdderType)
}

func TestAtomicF64AdderRaceInc(t *testing.T) {
	testF64AdderRaceInc(t, AtomicF64AdderType)
}

func TestAtomicF64AdderNotRaceDec(t *testing.T) {
	testF64AdderNotRaceDec(t, AtomicF64AdderType)
}

func TestAtomicF64AdderRaceDec(t *testing.T) {
	testF64AdderRaceDec(t, AtomicF64AdderType)
}

func TestAtomicF64AdderNotRaceAdd(t *testing.T) {
	testF64AdderNotRaceAdd(t, AtomicF64AdderType)
}

func TestAtomicF64AdderRaceAdd(t *testing.T) {
	testF64AdderRaceAdd(t, AtomicF64AdderType)
}
