package loader

import (
	"github.com/tmyksj/rlso11n/pkg/context"
	"github.com/tmyksj/rlso11n/pkg/rpc"
	"github.com/tmyksj/rlso11n/pkg/util/attempt"
	"github.com/tmyksj/rlso11n/pkg/util/wait"
	"strconv"
)

func Push(host string) {
	wait.UntilListen(host + ":" + strconv.Itoa(context.RpcPort()))

	attempt.UntilSucceed(func() error {
		return rpc.Call(host, rpc.MtdContextPush, &rpc.ReqContextPush{
			Dir:         context.Dir(),
			HostList:    context.HostList(),
			StarterAddr: context.StarterAddr(),
		}, &rpc.ResContextPush{})
	})
}
