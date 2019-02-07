package goadder

import (
	"testing"
)

func TestRandomCellAdderNotRaceInc(t *testing.T) {
	testAdderNotRaceInc(t, RandomCellAdderType)
}

func TestRandomCellAdderRaceInc(t *testing.T) {
	testAdderRaceInc(t, RandomCellAdderType)
}

func TestRandomCellAdderNotRaceDec(t *testing.T) {
	testAdderNotRaceDec(t, RandomCellAdderType)
}

func TestRandomCellAdderRaceDec(t *testing.T) {
	testAdderRaceDec(t, RandomCellAdderType)
}

func TestRandomCellAdderNotRaceAdd(t *testing.T) {
	testAdderNotRaceAdd(t, RandomCellAdderType)
}

func TestRandomCellAdderRaceAdd(t *testing.T) {
	testAdderRaceAdd(t, RandomCellAdderType)
}
