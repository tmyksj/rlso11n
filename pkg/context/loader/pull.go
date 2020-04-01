package loader

import (
	"github.com/tmyksj/rlso11n/pkg/context"
	"github.com/tmyksj/rlso11n/pkg/context/loader_proxy"
	"github.com/tmyksj/rlso11n/pkg/rpc"
	"github.com/tmyksj/rlso11n/pkg/util/attempt"
	"github.com/tmyksj/rlso11n/pkg/util/wait"
	"strconv"
)

func Pull(host string) {
	wait.UntilListenTcp(host + ":" + strconv.Itoa(context.RpcPort()))

	res := rpc.ResContextPull{}
	attempt.UntilSucceed(func() error {
		return rpc.Call(host, rpc.MtdContextPull, &rpc.ReqContextPull{}, &res)
	})

	loader_proxy.Load(&loader_proxy.LoadReq{
		Dir:           res.Dir,
		HostList:      res.HostList,
		StarterAddr:   res.StarterAddr,
	}, &loader_proxy.SetupReq{
		Dir: false,
	})
}
