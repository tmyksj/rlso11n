package rpc

import (
	"github.com/tmyksj/rlso11n/pkg/context"
	"github.com/tmyksj/rlso11n/pkg/context/loader_proxy"
)

const MtdContextPull = "Rpc.ContextPull"

type ReqContextPull struct {
}

type ResContextPull struct {
	Dir         string
	HostList    []context.Host
	ManagerAddr string
}

func (r *Rpc) ContextPull(req *ReqContextPull, res *ResContextPull) error {
	return do(func() error {
		if err := context.ReadyOrError(); err != nil {
			return err
		}

		res.Dir = context.Dir()
		res.HostList = context.HostList()
		res.ManagerAddr = context.ManagerAddr()
		return nil
	})
}

const MtdContextPush = "Rpc.ContextPush"

type ReqContextPush struct {
	Dir         string
	HostList    []context.Host
	ManagerAddr string
}

type ResContextPush struct {
}

func (r *Rpc) ContextPush(req *ReqContextPush, res *ResContextPush) error {
	return do(func() error {
		return loader_proxy.Load(&loader_proxy.LoadReq{
			Dir:         req.Dir,
			HostList:    req.HostList,
			ManagerAddr: req.ManagerAddr,
		}, &loader_proxy.SetupReq{
			Dir: true,
		})
	})
}
