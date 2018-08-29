// +build !linux

package longadder

import "time"

const (
	limit = (1 << 31) - 1
)

var start = time.Unix(0, 0)

// getRandomInt based on nano second resolution
func getRandomInt() int {
	return int(time.Since(start).Nanoseconds() & limit)
}
