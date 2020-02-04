package up

import (
	"github.com/tmyksj/rlso11n/app/logger"
	"github.com/tmyksj/rlso11n/pkg/context"
	"github.com/tmyksj/rlso11n/pkg/context/loader"
	"github.com/tmyksj/rlso11n/pkg/rpc"
	"github.com/tmyksj/rlso11n/pkg/util/attempt"
	"github.com/urfave/cli"
	"sync"
)

func Dockerd(_ *cli.Context) error {
	loader.LoadAsCommander()

	var wg sync.WaitGroup
	wg.Add(len(context.HostList()))

	for _, h := range context.HostList() {
		go func(host string) {
			attempt.UntilSucceed(func() error {
				return rpc.Call(host, rpc.MtdDockerdStart, &rpc.ReqDockerdStart{}, &rpc.ResDockerdStart{})
			})
			logger.Infof("cmd/up", "started dockerd")

			wg.Done()
		}(h)
	}

	wg.Wait()

	logger.Infof("cmd/up", "started all dockerd")

	return nil
}
