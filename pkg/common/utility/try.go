package utility

import (
	"github.com/tmyksj/rlso11n/pkg/common/errors"
	"github.com/tmyksj/rlso11n/pkg/common/logger"
	"time"
)

func Try(fun func() error) error {
	for i := 25; ; i *= 2 {
		if err := fun(); err == nil {
			return nil
		} else if i > 25600 {
			return errors.By(err, "timeout")
		}

		logger.Warn(pkg, "retries after %v milliseconds", i)
		time.Sleep(time.Duration(i) * time.Millisecond)
	}
}
