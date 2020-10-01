package bg

import (
	"github.com/tmyksj/rlso11n/pkg/rpc"
	"github.com/tmyksj/rlso11n/pkg/util"
	"github.com/urfave/cli"
)

func Server(_ *cli.Context) error {
	if err := rpc.Serve(); err != nil {
		return err
	}

	util.WaitInterrupt()

	return nil
}
