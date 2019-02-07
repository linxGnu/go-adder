package goadder

import (
	"testing"
)

func TestJDKF64AdderNotRaceInc(t *testing.T) {
	testF64AdderNotRaceInc(t, JDKF64AdderType)
}

func TestJDKF64AdderRaceInc(t *testing.T) {
	testF64AdderRaceInc(t, JDKF64AdderType)
}

func TestJDKF64AdderNotRaceDec(t *testing.T) {
	testF64AdderNotRaceDec(t, JDKF64AdderType)
}

func TestJDKF64AdderRaceDec(t *testing.T) {
	testF64AdderRaceDec(t, JDKF64AdderType)
}

func TestJDKF64AdderNotRaceAdd(t *testing.T) {
	testF64AdderNotRaceAdd(t, JDKF64AdderType)
}

func TestJDKF64AdderRaceAdd(t *testing.T) {
	testF64AdderRaceAdd(t, JDKF64AdderType)
}
