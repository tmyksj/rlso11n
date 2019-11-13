package attempt

import (
	"time"
)

func UntilSucceed(f func() error, d time.Duration) {
	for {
		if err := f(); err != nil {
			time.Sleep(d)
		} else {
			break
		}
	}
}
