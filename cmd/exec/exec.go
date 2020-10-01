package exec

import (
	"fmt"
	"github.com/tmyksj/rlso11n/app/logger"
	"github.com/tmyksj/rlso11n/pkg/context"
	"github.com/tmyksj/rlso11n/pkg/context/loader"
	"github.com/tmyksj/rlso11n/pkg/rpc"
	"github.com/tmyksj/rlso11n/pkg/util"
	"github.com/urfave/cli"
	"strconv"
	"strings"
	"sync"
)

func Exec(c *cli.Context) error {
	loader.LoadAsCommander()

	var targetList []target
	var command []string

	targetSep := strings.Index(c.Command.Name, "@")
	commandSep := strings.Index(c.Command.Name, "/")
	if commandSep == -1 {
		targetList = parseTargetList(c.Command.Name[targetSep+1:])
		command = c.Args()[1:]
	} else {
		targetList = parseTargetList(c.Command.Name[targetSep+1 : commandSep])
		command = append([]string{c.Command.Name[commandSep+1:]}, c.Args()[1:]...)
	}

	var output = make([]string, len(targetList))

	var wg sync.WaitGroup
	wg.Add(len(targetList))

	for i, t := range targetList {
		go func(idx int, tgt target) {
			if commandSep == -1 {
				exec(&wg, tgt.host, command, &output[idx])
			} else {
				replaced := make([]string, len(command))
				for j, _ := range command {
					replaced[j] = strings.ReplaceAll(command[j], "$i", strconv.Itoa(tgt.i))
				}

				switch replaced[0] {
				case "docker":
					execDocker(&wg, tgt.host, replaced, &output[idx])
					break
				default:
					logger.Error(pkg, "-> %v, %v is not supported", tgt.host, replaced[0])
					wg.Done()
					break
				}
			}
		}(i, t)
	}

	wg.Wait()

	logger.Info(pkg, "executed command at all targets")

	for i, t := range targetList {
		fmt.Println("[" + t.host + "]")
		fmt.Println(output[i])
	}

	return nil
}

func exec(wg *sync.WaitGroup, host string, command []string, output *string) {
	defer wg.Done()

	util.TryUntilSucceed(func() error {
		res := rpc.ResExtRun{}
		if err := rpc.Call(host, rpc.MtdExtRun, &rpc.ReqExtRun{Args: command}, &res); err != nil {
			return err
		}

		*output = res.Stdout

		return nil
	})

	logger.Info(pkg, "-> %v, executed command", host)
}

type target struct {
	i    int
	host string
}

func parseTargetList(t string) []target {
	switch t {
	case "all":
		return parseTargetListAll()
	case "manager":
		return parseTargetListManager()
	case "worker":
		return parseTargetListWorker()
	default:
		return parseTargetListSpecified(t)
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
		if h == context.StarterAddr() {
			tl = append(tl, target{i: i, host: h})
		}
	}

	return tl
}

func parseTargetListWorker() []target {
	var tl []target
	for i, h := range context.HostList() {
		if h != context.StarterAddr() {
			tl = append(tl, target{i: i, host: h})
		}
	}

	return tl
}

func parseTargetListSpecified(t string) []target {
	hl := context.HostList()

	var tl []target
	for _, tSep := range strings.Split(t, ",") {
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
