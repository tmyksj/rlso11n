package cmd

import (
	"github.com/tmyksj/rlso11n/pkg/common/logger"
	"github.com/tmyksj/rlso11n/pkg/core"
	"github.com/urfave/cli/v2"
)

func Perform(fun func() error) error {
	core.Initialize()

	err := fun()

	if err != nil {
		logger.Error(pkg, "%v", err)

		core.Finalize()
		return cli.Exit(err, 1)
	}

	core.Finalize()
	return nil
}
