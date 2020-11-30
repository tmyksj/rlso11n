package up

import (
	"github.com/tmyksj/rlso11n/app/logger"
	"github.com/tmyksj/rlso11n/cmd"
	"github.com/tmyksj/rlso11n/pkg/context"
	"github.com/tmyksj/rlso11n/pkg/context/loader"
	"github.com/tmyksj/rlso11n/pkg/rpc"
	"github.com/tmyksj/rlso11n/pkg/util"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
)

func Dockerd(c *cli.Context) error {
	cServer := c.String("server")

	if err := loader.LoadAsClient(cServer); err != nil {
		logger.Error(pkg, "failed to load context, %v", err)
		return cmd.ExitError(err)
	}

	eg := errgroup.Group{}

	for _, host := range context.HostList() {
		host := host
		eg.Go(func() error {
			if err := util.Try(func() error {
				return rpc.Call(host.Addr, rpc.MtdDockerdStart, &rpc.ReqDockerdStart{}, &rpc.ResDockerdStart{})
			}); err != nil {
				logger.Error(pkg, "failed to start dockerd at %v, %v", host.Name, err)
				return err
			}

			logger.Info(pkg, "started dockerd at %v", host.Name)
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return cmd.ExitError(err)
	}

	logger.Info(pkg, "started dockerd at all hosts")

	return nil
}
