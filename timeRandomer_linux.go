// +build linux

package longadder

import "time"

// getRandomInt based on nano second resolution
func getRandomInt() int {
	return time.Now().Nanosecond()
}
