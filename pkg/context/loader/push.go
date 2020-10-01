package loader

import (
	"github.com/tmyksj/rlso11n/pkg/context"
	"github.com/tmyksj/rlso11n/pkg/rpc"
	"github.com/tmyksj/rlso11n/pkg/util"
	"strconv"
)

func Push(host string) {
	util.WaitUntilListenTcp(host + ":" + strconv.Itoa(context.RpcPort()))

	util.TryUntilSucceed(func() error {
		return rpc.Call(host, rpc.MtdContextPush, &rpc.ReqContextPush{
			Dir:         context.Dir(),
			HostList:    context.HostList(),
			StarterAddr: context.StarterAddr(),
		}, &rpc.ResContextPush{})
	})
}
