package bg

import (
	"github.com/tmyksj/rootless-orchestration/logger"
	"github.com/tmyksj/rootless-orchestration/pkg/wait"
	"github.com/tmyksj/rootless-orchestration/rpc"
	"github.com/urfave/cli"
)

func Server(c *cli.Context) error {
	if err := rpc.Serve(); err != nil {
		logger.Errorf("cmd/bg", "fail to start server")
		logger.Errorf("cmd/bg", "%v", err)
		return err
	}

	logger.Infof("cmd/bg", "succeed to start server")

	wait.Interrupt()

	return nil
}
