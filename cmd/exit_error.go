package cmd

import (
	"github.com/tmyksj/rlso11n/app/core"
	"github.com/urfave/cli/v2"
)

func ExitError(message interface{}) error {
	core.Finalize()
	return cli.Exit(message, 1)
}
