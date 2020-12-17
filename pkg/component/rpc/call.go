package rpc

import (
	"github.com/tmyksj/rlso11n/pkg/common/errors"
	"github.com/tmyksj/rlso11n/pkg/common/logger"
	"github.com/tmyksj/rlso11n/pkg/common/object"
	"github.com/tmyksj/rlso11n/pkg/core"
	"net/rpc"
	"strconv"
	"sync"
)

var clientMap = make(map[string]*rpc.Client)
var clientMu sync.Mutex

func call(ctx *core.Context, node *object.Node, serviceMethod string, req interface{}, res interface{}) error {
	clientMu.Lock()

	client, ok := clientMap[node.Addr]
	if !ok {
		c, err := rpc.DialHTTP("tcp", node.Addr+":"+strconv.Itoa(ctx.RpcPort()))
		if err != nil {
			clientMu.Unlock()
			return errors.By(err, "failed to connect to %v", node.Name)
		}

		client = c
		clientMap[node.Addr] = client

		logger.Info(pkg, "succeed to connect to %v", node.Name)
	}

	clientMu.Unlock()

	if err := client.Call(serviceMethod, req, res); err != nil {
		return errors.By(err, "failed to call %v", serviceMethod)
	}

	return nil
}
