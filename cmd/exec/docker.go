package exec

import (
	"github.com/tmyksj/rlso11n/app/logger"
	"github.com/tmyksj/rlso11n/pkg/rpc"
	"github.com/tmyksj/rlso11n/pkg/util/attempt"
	"sync"
)

func execDocker(wg *sync.WaitGroup, host string, command []string, output *string) {
	defer wg.Done()

	attempt.UntilSucceed(func() error {
		res := rpc.ResDockerRun{}
		if err := rpc.Call(host, rpc.MtdDockerRun, &rpc.ReqDockerRun{Args: command[1:]}, &res); err != nil {
			return err
		}

		*output = res.Output

		return nil
	})
	logger.Infof("cmd/exec", "executed command")
}
