package cmd

import (
	"github.com/tmyksj/rlso11n/app/core"
	"github.com/urfave/cli/v2"
)

func After(_ *cli.Context) error {
	core.Finalize()
	return nil
}
