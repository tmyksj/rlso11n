package loader

import (
	"github.com/tmyksj/rlso11n/pkg/context"
	"github.com/tmyksj/rlso11n/pkg/context/loader_proxy"
	"github.com/tmyksj/rlso11n/pkg/errors"
	"github.com/tmyksj/rlso11n/pkg/rpc"
	"github.com/tmyksj/rlso11n/pkg/util"
	"strconv"
)

func Pull(addr string) error {
	if err := util.WaitListenTcp(addr + ":" + strconv.Itoa(context.RpcPort())); err != nil {
		return errors.By(err, "failed to pull context from %v", addr)
	}

	res := rpc.ResContextPull{}
	if err := util.Try(func() error {
		return rpc.Call(addr, rpc.MtdContextPull, &rpc.ReqContextPull{}, &res)
	}); err != nil {
		return errors.By(err, "failed to pull context")
	}

	if err := loader_proxy.Load(&loader_proxy.LoadReq{
		Dir:         res.Dir,
		HostList:    res.HostList,
		ManagerAddr: res.ManagerAddr,
	}, &loader_proxy.SetupReq{
		Dir: false,
	}); err != nil {
		return errors.By(err, "failed to pull context from %v", addr)
	}

	return nil
}
