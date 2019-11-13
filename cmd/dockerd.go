package cmd

import (
	"github.com/tmyksj/rootless-orchestration/context"
	"github.com/tmyksj/rootless-orchestration/logger"
	"github.com/tmyksj/rootless-orchestration/pkg/attempt"
	"github.com/tmyksj/rootless-orchestration/rpc"
	"github.com/urfave/cli"
	"sync"
	"time"
)

func Dockerd(c *cli.Context) error {
	context.InitFromEnv()

	var wg sync.WaitGroup
	wg.Add(len(context.HostList()))

	for _, h := range context.HostList() {
		go func(host string) {
			attempt.UntilSucceed(func() error {
				return rpc.Call(host, rpc.MtdDockerdStart, &rpc.ReqDockerdStart{}, &rpc.ResDockerdStart{})
			}, 100*time.Millisecond)
			logger.Infof("cmd", "started dockerd")

			wg.Done()
		}(h)
	}

	wg.Wait()

	logger.Infof("cmd", "started all dockerd")

	return nil
}
