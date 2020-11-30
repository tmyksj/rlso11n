package cmd

import (
	"github.com/tmyksj/rlso11n/app/core"
	"github.com/urfave/cli/v2"
)

func Before(_ *cli.Context) error {
	core.Initialize()
	return nil
}
