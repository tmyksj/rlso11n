package serve_rpc

import (
	"github.com/tmyksj/rlso11n/app/logger"
	"github.com/tmyksj/rlso11n/cmd"
	"github.com/tmyksj/rlso11n/pkg/rpc"
	"github.com/tmyksj/rlso11n/pkg/util"
	"github.com/urfave/cli/v2"
)

func Action(_ *cli.Context) error {
	if err := rpc.Serve(); err != nil {
		logger.Error(pkg, "failed to serve, %v", err)
		return cmd.ExitError(err)
	}

	util.WaitInterrupt()

	return nil
}
