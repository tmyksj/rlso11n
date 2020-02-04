package exec

import (
	"fmt"
	"github.com/tmyksj/rlso11n/app/logger"
	"github.com/tmyksj/rlso11n/pkg/context"
	"github.com/tmyksj/rlso11n/pkg/context/loader"
	"github.com/tmyksj/rlso11n/pkg/rpc"
	"github.com/tmyksj/rlso11n/pkg/util/attempt"
	"github.com/urfave/cli"
	"strconv"
	"strings"
	"sync"
)

func DockerAt(c *cli.Context) error {
	loader.LoadAsCommander()

	target := parseTarget(c.Command.Name[len("exec/docker@"):])
	logger.Infof("cmd/exec", "execute docker command at %v", target)

	var output = make([]string, len(target))

	var wg sync.WaitGroup
	wg.Add(len(target))

	for i, h := range target {
		go func(idx int, host string) {
			attempt.UntilSucceed(func() error {
				res := rpc.ResDockerRun{}
				if err := rpc.Call(host, rpc.MtdDockerRun, &rpc.ReqDockerRun{Args: c.Args()[1:]}, &res); err != nil {
					return err
				}

				output[idx] = res.Output

				return nil
			})
			logger.Infof("cmd/exec", "executed docker command")

			wg.Done()
		}(i, h)
	}

	wg.Wait()

	logger.Infof("cmd/exec", "executed docker command at all target")

	for i, h := range target {
		fmt.Println("[" + h + "]")
		fmt.Println(output[i])
	}

	return nil
}

func parseTarget(t string) []string {
	switch t {
	case "all":
		return parseTargetAll()
	case "manager":
		return parseTargetManager()
	case "worker":
		return parseTargetWorker()
	default:
		return parseTargetSpecified(t)
	}
}

func parseTargetAll() []string {
	return context.HostList()
}

func parseTargetManager() []string {
	return []string{context.StarterAddr()}
}

func parseTargetSpecified(t string) []string {
	h := context.HostList()

	var arr []string
	for _, tSep := range strings.Split(t, ",") {
		idx := strings.Index(tSep, "-")
		if idx == -1 {
			if n, err := strconv.Atoi(tSep); err == nil {
				arr = append(arr, h[n%len(h)])
			}
		} else {
			tSep0 := tSep[0:idx]
			tSep1 := tSep[idx+1:]

			n0, err0 := strconv.Atoi(tSep0)
			n1, err1 := strconv.Atoi(tSep1)

			if err0 == nil && err1 == nil {
				for i := n0; i <= n1; i++ {
					arr = append(arr, h[i%len(h)])
				}
			}
		}
	}

	return arr
}

func parseTargetWorker() []string {
	var arr []string
	for _, h := range context.HostList() {
		if h != context.StarterAddr() {
			arr = append(arr, h)
		}
	}

	return arr
}
