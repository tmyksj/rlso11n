package rpc

import (
	"github.com/tmyksj/rlso11n/pkg/common/errors"
	"github.com/tmyksj/rlso11n/pkg/common/logger"
	"github.com/tmyksj/rlso11n/pkg/component/exec"
	"github.com/tmyksj/rlso11n/pkg/core"
)

type Rpc struct {
}

func (rpc *Rpc) AnyRun(req *ReqAnyRun, res *ResAnyRun) error {
	return rpc.perform(func() error {
		ctx := core.GetContext()
		if err := ctx.CheckReady(); err != nil {
			return err
		}

		ret, err := exec.Any(ctx).Run(&exec.ReqAnyRun{
			Name: req.Name,
			Args: req.Args,
		})

		if err != nil {
			return err
		}

		res.Stdout = ret.Stdout

		return nil
	})
}

func (rpc *Rpc) ContextPull(req *ReqContextPull, res *ResContextPull) error {
	return rpc.perform(func() error {
		ctx := core.GetContext()
		if err := ctx.CheckReady(); err != nil {
			return err
		}

		res.Dir = ctx.Dir()
		res.ManagerNode = ctx.ManagerNode()
		res.NodeList = ctx.NodeList()

		return nil
	})
}

func (rpc *Rpc) ContextPush(req *ReqContextPush, res *ResContextPush) error {
	return rpc.perform(func() error {
		constructor := core.Constructor{
			Dir:         req.Dir,
			ManagerNode: req.ManagerNode,
			NodeList:    req.NodeList,
		}

		if err := constructor.Load(); err != nil {
			return err
		}

		if err := constructor.Setup(); err != nil {
			return err
		}

		if err := constructor.ToReady(); err != nil {
			return err
		}

		return nil
	})
}

func (rpc *Rpc) DockerRun(req *ReqDockerRun, res *ResDockerRun) error {
	return rpc.perform(func() error {
		ctx := core.GetContext()
		if err := ctx.CheckReady(); err != nil {
			return err
		}

		ret, err := exec.Docker(ctx).Run(&exec.ReqDockerRun{
			Args: req.Args,
		})

		if err != nil {
			return err
		}

		res.Stdout = ret.Stdout

		return nil
	})
}

func (rpc *Rpc) DockerdStart(req *ReqDockerdStart, res *ResDockerdStart) error {
	return rpc.perform(func() error {
		ctx := core.GetContext()
		if err := ctx.CheckReady(); err != nil {
			return err
		}

		_, err := exec.Dockerd(ctx).Start(&exec.ReqDockerdStart{})

		if err != nil {
			return err
		}

		return nil
	})
}

func (rpc *Rpc) perform(fun func() error) error {
	if err := fun(); err != nil {
		logger.Error(pkg, "rpc error, %v", err)
		return errors.By(err, "rpc error")
	}

	return nil
}
