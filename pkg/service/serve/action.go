package serve

import (
	"fmt"
	"github.com/tmyksj/rlso11n/pkg/common/errors"
	"github.com/tmyksj/rlso11n/pkg/common/logger"
	"github.com/tmyksj/rlso11n/pkg/common/object"
	"github.com/tmyksj/rlso11n/pkg/common/utility"
	"github.com/tmyksj/rlso11n/pkg/component/rpc"
	"github.com/tmyksj/rlso11n/pkg/core"
	"golang.org/x/sync/errgroup"
	"io"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

type ObjAction struct {
	ctx *core.Context
}

type ReqActionRun struct {
}

type ResActionRun struct {
}

func Action(ctx *core.Context) *ObjAction {
	return &ObjAction{
		ctx: ctx,
	}
}

func (obj *ObjAction) Run(req *ReqActionRun) (*ResActionRun, error) {
	myNode := obj.ctx.MyNode()
	nodeList := obj.ctx.NodeList()

	eg := errgroup.Group{}

	for _, node := range nodeList {
		node := node
		eg.Go(func() error {
			if myNode != nil && node.Addr == myNode.Addr {
				return obj.run(node)
			} else {
				return obj.runViaSsh(node)
			}
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return &ResActionRun{}, nil
}

func (obj *ObjAction) push(node *object.Node) error {
	if err := utility.WaitListenTcp(node.Addr + ":" + strconv.Itoa(obj.ctx.RpcPort())); err != nil {
		return errors.By(err, "failed to push context to %v", node.Addr)
	}

	err := utility.Try(func() error {
		_, err := rpc.Context(obj.ctx, node).Push(&rpc.ReqContextPush{
			Dir:         obj.ctx.Dir(),
			ManagerNode: obj.ctx.ManagerNode(),
			NodeList:    obj.ctx.NodeList(),
		})

		return err
	})

	if err != nil {
		return errors.By(err, "failed to push context to %v", node.Addr)
	}

	return nil
}

func (obj *ObjAction) run(node *object.Node) error {
	svr := exec.Command("rlso11n", "serve-rpc")
	svr.Env = obj.ctx.Env()
	svr.Stdout = os.Stderr
	svr.Stderr = os.Stderr

	if err := svr.Start(); err != nil {
		logger.Error(pkg, "failed to start worker of %v, %v", node.Name, err)
		return err
	}

	if err := obj.push(node); err != nil {
		logger.Error(pkg, "failed to push context to %v, %v", node.Name, err)
		return err
	}

	utility.WaitInterrupt()

	if err := svr.Process.Signal(syscall.SIGINT); err != nil {
		logger.Warn(pkg, "failed to send sigint to %v, %v", node.Name, err)

		if err := svr.Process.Kill(); err != nil {
			logger.Error(pkg, "failed to kill at %v, %v", node.Name, err)
			return err
		}
	}

	if err := svr.Wait(); err != nil {
		logger.Warn(pkg, "failed to wait at %v, %v", node.Name, err)
	}

	return nil
}

func (obj *ObjAction) runViaSsh(node *object.Node) error {
	svr := exec.Command("ssh", "-tt", node.Addr)
	svr.Env = obj.ctx.Env()
	svr.Stdout = os.Stderr
	svr.Stderr = os.Stderr
	svr.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	stdin, err := svr.StdinPipe()

	if err != nil {
		logger.Error(pkg, "failed to open pipe (stdin) to %v, %v", node.Name, err)
		return err
	}

	defer func() {
		if err := stdin.Close(); err != nil && err != io.EOF {
			logger.Error(pkg, "failed to close pipe (stdin) of %v, %v", node.Name, err)
		}
	}()

	if err := svr.Start(); err != nil {
		logger.Error(pkg, "failed to start ssh to %v, %v", node.Name, err)
		return err
	}

	if _, err := fmt.Fprintln(stdin, "rlso11n serve-rpc"); err != nil {
		logger.Error(pkg, "failed to start worker of %v, %v", node.Name, err)
		return err
	}

	if err := obj.push(node); err != nil {
		logger.Error(pkg, "failed to push context to %v, %v", node.Name, err)
		return err
	}

	utility.WaitInterrupt()

	if _, err := fmt.Fprintln(stdin, "\x03"); err != nil {
		logger.Error(pkg, "failed to send sigint to %v, %v", node.Name, err)
	}

	if _, err := fmt.Fprintln(stdin, "exit"); err != nil {
		logger.Warn(pkg, "failed to send exit at %v, %v", node.Name, err)
	}

	if err := svr.Wait(); err != nil {
		logger.Warn(pkg, "failed to wait ssh of %v, %v", node.Name, err)
	}

	return nil
}
