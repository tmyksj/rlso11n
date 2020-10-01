package util

import (
	"time"
)

func TryUntilSucceed(f func() error) {
	for {
		if err := f(); err != nil {
			time.Sleep(100 * time.Millisecond)
		} else {
			break
		}
	}
}
