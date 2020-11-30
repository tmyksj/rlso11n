package util

import (
	"github.com/tmyksj/rlso11n/app/logger"
	"github.com/tmyksj/rlso11n/pkg/errors"
	"time"
)

func Try(f func() error) error {
	for i := 25; ; i *= 2 {
		if err := f(); err == nil {
			return nil
		} else if i > 25600 {
			return errors.By(err, "timeout")
		}

		logger.Warn(pkg, "retries after %v milliseconds", i)
		time.Sleep(time.Duration(i) * time.Millisecond)
	}
}
