package bg

import (
	"github.com/tmyksj/rlso11n/app/logger"
	"github.com/tmyksj/rlso11n/pkg/rpc"
	"github.com/tmyksj/rlso11n/pkg/util/wait"
	"github.com/urfave/cli"
)

func Server(_ *cli.Context) error {
	if err := rpc.Serve(); err != nil {
		logger.Errorf("cmd/bg", "fail to start server, %v", err)
		return err
	}

	logger.Infof("cmd/bg", "succeed to start server")

	wait.Interrupt()

	return nil
}
