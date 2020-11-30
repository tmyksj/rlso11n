package exec

import (
	"fmt"
	"github.com/tmyksj/rlso11n/app/logger"
	"github.com/tmyksj/rlso11n/cmd"
	"github.com/tmyksj/rlso11n/pkg/context"
	"github.com/tmyksj/rlso11n/pkg/context/loader"
	"github.com/tmyksj/rlso11n/pkg/errors"
	"github.com/tmyksj/rlso11n/pkg/rpc"
	"github.com/tmyksj/rlso11n/pkg/util"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
	"strconv"
	"strings"
)

func Action(c *cli.Context) error {
	cDocker := c.Bool("docker")
	cNodes := c.String("nodes")
	cServer := c.String("server")
	cArgs := c.Args().Slice()

	if len(cArgs) == 0 {
		err := errors.By(nil, "arguments are required")

		logger.Error(pkg, "bad request, %v", err)
		return cmd.ExitError(err)
	}

	if err := loader.LoadAsClient(cServer); err != nil {
		logger.Error(pkg, "failed to load context, %v", err)
		return cmd.ExitError(err)
	}

	targetList := parseTargetList(cNodes)

	eg := errgroup.Group{}
	output := make([]string, len(targetList))

	for i, target := range targetList {
		i := i
		target := target
		eg.Go(func() error {
			args := make([]string, len(cArgs))
			for j := range cArgs {
				args[j] = strings.ReplaceAll(cArgs[j], "$i", strconv.Itoa(target.i))
			}

			if cDocker {
				return execDocker(target.host, args, &output[i])
			} else {
				return exec(target.host, args, &output[i])
			}
		})
	}

	if err := eg.Wait(); err != nil {
		return cmd.ExitError(err)
	}

	logger.Info(pkg, "executed command at all targets")

	for i, target := range targetList {
		fmt.Println("[" + target.host.Name + "]")
		fmt.Println(output[i])
	}

	return nil
}

func exec(host context.Host, command []string, output *string) error {
	err := util.Try(func() error {
		res := rpc.ResExtRun{}
		if err := rpc.Call(host.Addr, rpc.MtdExtRun, &rpc.ReqExtRun{Args: command[1:], Name: command[0]}, &res); err != nil {
			return err
		}

		*output = res.Stdout

		return nil
	})
	if err != nil {
		logger.Error(pkg, "failed to execute command at %v, %v", host.Name, err)
		return err
	}

	logger.Info(pkg, "executed command at %v", host.Name)

	return nil
}

func execDocker(host context.Host, command []string, output *string) error {
	err := util.Try(func() error {
		res := rpc.ResDockerRun{}
		if err := rpc.Call(host.Addr, rpc.MtdDockerRun, &rpc.ReqDockerRun{Args: command}, &res); err != nil {
			return err
		}

		*output = res.Output

		return nil
	})
	if err != nil {
		logger.Error(pkg, "failed to execute command at %v, %v", host.Name, err)
		return err
	}

	logger.Info(pkg, "executed command at %v", host.Name)

	return nil
}

type target struct {
	i    int
	host context.Host
}

func parseTargetList(list string) []target {
	switch list {
	case "all":
		return parseTargetListAll()
	case "manager":
		return parseTargetListManager()
	default:
		return parseTargetListSpecified(list)
	}
}

func parseTargetListAll() []target {
	var tl []target
	for i, h := range context.HostList() {
		tl = append(tl, target{i: i, host: h})
	}

	return tl
}

func parseTargetListManager() []target {
	var tl []target
	for i, h := range context.HostList() {
		if h.Addr == context.ManagerAddr() {
			tl = append(tl, target{i: i, host: h})
		}
	}

	return tl
}

func parseTargetListSpecified(list string) []target {
	hl := context.HostList()

	var tl []target
	for _, tSep := range strings.Split(list, ",") {
		idx := strings.Index(tSep, "-")
		if idx == -1 {
			if n, err := strconv.Atoi(tSep); err == nil {
				tl = append(tl, target{i: n, host: hl[n%len(hl)]})
			}
		} else {
			tSep0 := tSep[0:idx]
			tSep1 := tSep[idx+1:]

			n0, err0 := strconv.Atoi(tSep0)
			n1, err1 := strconv.Atoi(tSep1)

			if err0 == nil && err1 == nil {
				for i := n0; i <= n1; i++ {
					tl = append(tl, target{i: i, host: hl[i%len(hl)]})
				}
			}
		}
	}

	return tl
}
