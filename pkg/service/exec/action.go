package exec

import (
	"github.com/tmyksj/rlso11n/pkg/common/logger"
	"github.com/tmyksj/rlso11n/pkg/common/object"
	"github.com/tmyksj/rlso11n/pkg/common/utility"
	"github.com/tmyksj/rlso11n/pkg/component/rpc"
	"github.com/tmyksj/rlso11n/pkg/core"
	"golang.org/x/sync/errgroup"
	"strconv"
	"strings"
)

type ObjAction struct {
	ctx *core.Context
}

type ReqActionRun struct {
	Docker bool
	Nodes  string
	Args   []string
}

type ResActionRun struct {
	Output []*ResActionRunOutput
}

type ResActionRunOutput struct {
	Node   *object.Node
	Stdout string
}

func Action(ctx *core.Context) *ObjAction {
	return &ObjAction{
		ctx: ctx,
	}
}

func (obj *ObjAction) Run(req *ReqActionRun) (*ResActionRun, error) {
	targetList := obj.parseNodes(req.Nodes)

	eg := errgroup.Group{}
	output := make([]*ResActionRunOutput, len(targetList))

	for i, target := range targetList {
		i := i
		target := target

		output[i] = &ResActionRunOutput{
			Node: target.node,
		}

		eg.Go(func() error {
			args := make([]string, len(req.Args))
			for j := range req.Args {
				args[j] = strings.ReplaceAll(req.Args[j], "$i", strconv.Itoa(target.index))
			}

			if req.Docker {
				return obj.execDocker(target.node, args, output[i])
			} else {
				return obj.exec(target.node, args, output[i])
			}
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return &ResActionRun{
		Output: output,
	}, nil
}

type target struct {
	index int
	node  *object.Node
}

func (obj *ObjAction) exec(node *object.Node, args []string, output *ResActionRunOutput) error {
	err := utility.Try(func() error {
		res, err := rpc.Any(obj.ctx, node).Run(&rpc.ReqAnyRun{
			Name: args[0],
			Args: args[1:],
		})

		if err != nil {
			return err
		}

		output.Stdout = res.Stdout

		return nil
	})

	if err != nil {
		logger.Error(pkg, "failed to execute command at %v, %v", node.Name, err)
		return err
	}

	logger.Info(pkg, "executed command at %v", node.Name)

	return nil
}

func (obj *ObjAction) execDocker(node *object.Node, args []string, output *ResActionRunOutput) error {
	err := utility.Try(func() error {
		res, err := rpc.Docker(obj.ctx, node).Run(&rpc.ReqDockerRun{
			Args: args,
		})

		if err != nil {
			return err
		}

		output.Stdout = res.Stdout

		return nil
	})

	if err != nil {
		logger.Error(pkg, "failed to execute command at %v, %v", node.Name, err)
		return err
	}

	logger.Info(pkg, "executed command at %v", node.Name)

	return nil
}

func (obj *ObjAction) parseNodes(nodes string) []target {
	switch nodes {
	case "all":
		return obj.parseNodesAsAll()
	case "manager":
		return obj.parseNodesAsManager()
	default:
		return obj.parseNodesAsSpecified(nodes)
	}
}

func (obj *ObjAction) parseNodesAsAll() []target {
	var targetList []target
	for i, node := range obj.ctx.NodeList() {
		targetList = append(targetList, target{
			index: i,
			node:  node,
		})
	}

	return targetList
}

func (obj *ObjAction) parseNodesAsManager() []target {
	var targetList []target
	for i, node := range obj.ctx.NodeList() {
		if node.Addr == obj.ctx.ManagerNode().Addr {
			targetList = append(targetList, target{
				index: i,
				node:  node,
			})
		}
	}

	return targetList
}

func (obj *ObjAction) parseNodesAsSpecified(nodes string) []target {
	nodeList := obj.ctx.NodeList()

	var targetList []target
	for _, tSep := range strings.Split(nodes, ",") {
		idx := strings.Index(tSep, "-")
		if idx == -1 {
			if n, err := strconv.Atoi(tSep); err == nil {
				targetList = append(targetList, target{
					index: n,
					node:  nodeList[n%len(nodeList)],
				})
			}
		} else {
			tSep0 := tSep[0:idx]
			tSep1 := tSep[idx+1:]

			n0, err0 := strconv.Atoi(tSep0)
			n1, err1 := strconv.Atoi(tSep1)

			if err0 == nil && err1 == nil {
				for i := n0; i <= n1; i++ {
					targetList = append(targetList, target{
						index: i,
						node:  nodeList[i%len(nodeList)],
					})
				}
			}
		}
	}

	return targetList
}
