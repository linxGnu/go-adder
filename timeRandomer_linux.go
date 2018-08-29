// +build linux

package longadder

import "time"

const (
	limit = (1 << 31) - 1
)

// getRandomInt based on nano second resolution
func getRandomInt() int {
	return time.Now().Nanosecond() & limit
}
