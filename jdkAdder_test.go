package goadder

import (
	"testing"
)

func TestJDKAdderNotRaceInc(t *testing.T) {
	testAdderNotRaceInc(t, JDKAdderType)
}

func TestJDKAdderRaceInc(t *testing.T) {
	testAdderRaceInc(t, JDKAdderType)
}

func TestJDKAdderNotRaceDec(t *testing.T) {
	testAdderNotRaceDec(t, JDKAdderType)
}

func TestJDKAdderRaceDec(t *testing.T) {
	testAdderRaceDec(t, JDKAdderType)
}

func TestJDKAdderNotRaceAdd(t *testing.T) {
	testAdderNotRaceAdd(t, JDKAdderType)
}

func TestJDKAdderRaceAdd(t *testing.T) {
	testAdderRaceAdd(t, JDKAdderType)
}
