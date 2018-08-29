package longadder

import (
	"github.com/valyala/fastrand"
)

const (
	limit = (1 << 31) - 1
)

func getRandomInt() int {
	return int(fastrand.Uint32() & limit)
}
