package rpc

import (
	"github.com/tmyksj/rlso11n/app/logger"
	"github.com/tmyksj/rlso11n/pkg/errors"
)

type Rpc int

func do(f func() error) error {
	if err := f(); err != nil {
		logger.Error(pkg, "rpc error, %v", err)
		return errors.By(err, "rpc error")
	} else {
		return nil
	}
}
