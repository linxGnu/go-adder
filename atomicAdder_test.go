package goadder

import (
	"testing"
)

func TestAtomicAdderNotRaceInc(t *testing.T) {
	testAdderNotRaceInc(t, AtomicAdderType)
}

func TestAtomicAdderRaceInc(t *testing.T) {
	testAdderRaceInc(t, AtomicAdderType)
}

func TestAtomicAdderNotRaceDec(t *testing.T) {
	testAdderNotRaceDec(t, AtomicAdderType)
}

func TestAtomicAdderRaceDec(t *testing.T) {
	testAdderRaceDec(t, AtomicAdderType)
}

func TestAtomicAdderNotRaceAdd(t *testing.T) {
	testAdderNotRaceAdd(t, AtomicAdderType)
}

func TestAtomicAdderRaceAdd(t *testing.T) {
	testAdderRaceAdd(t, AtomicAdderType)
}
