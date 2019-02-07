package goadder

import (
	"testing"
)

func TestMutexAdderNotRaceInc(t *testing.T) {
	testAdderNotRaceInc(t, MutexAdderType)
}

func TestMutexAdderRaceInc(t *testing.T) {
	testAdderRaceInc(t, MutexAdderType)
}

func TestMutexAdderNotRaceDec(t *testing.T) {
	testAdderNotRaceDec(t, MutexAdderType)
}

func TestMutexAdderRaceDec(t *testing.T) {
	testAdderRaceDec(t, MutexAdderType)
}

func TestMutexAdderNotRaceAdd(t *testing.T) {
	testAdderNotRaceAdd(t, MutexAdderType)
}

func TestMutexAdderRaceAdd(t *testing.T) {
	testAdderRaceAdd(t, MutexAdderType)
}
