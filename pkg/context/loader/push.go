package loader

import (
	"github.com/tmyksj/rlso11n/pkg/context"
	"github.com/tmyksj/rlso11n/pkg/errors"
	"github.com/tmyksj/rlso11n/pkg/rpc"
	"github.com/tmyksj/rlso11n/pkg/util"
	"strconv"
)

func Push(addr string) error {
	if err := util.WaitListenTcp(addr + ":" + strconv.Itoa(context.RpcPort())); err != nil {
		return errors.By(err, "failed to push context to %v", addr)
	}

	if err := util.Try(func() error {
		return rpc.Call(addr, rpc.MtdContextPush, &rpc.ReqContextPush{
			Dir:         context.Dir(),
			HostList:    context.HostList(),
			ManagerAddr: context.ManagerAddr(),
		}, &rpc.ResContextPush{})
	}); err != nil {
		return errors.By(err, "failed to push context to %v", addr)
	}

	return nil
}
