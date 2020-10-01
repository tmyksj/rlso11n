package loader

import (
	"github.com/tmyksj/rlso11n/pkg/context"
	"github.com/tmyksj/rlso11n/pkg/context/loader_proxy"
	"github.com/tmyksj/rlso11n/pkg/rpc"
	"github.com/tmyksj/rlso11n/pkg/util"
	"strconv"
)

func Pull(host string) {
	util.WaitUntilListenTcp(host + ":" + strconv.Itoa(context.RpcPort()))

	res := rpc.ResContextPull{}
	util.TryUntilSucceed(func() error {
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
