package rpc

import (
	"github.com/tmyksj/rlso11n/app/logger"
	"github.com/tmyksj/rlso11n/pkg/context"
	"github.com/tmyksj/rlso11n/pkg/errors"
	"net/rpc"
	"strconv"
	"sync"
)

var clientMap = make(map[string]*rpc.Client)
var clientMu sync.Mutex

func Call(addr string, serviceMethod string, req interface{}, res interface{}) error {
	clientMu.Lock()

	client, ok := clientMap[addr]
	if !ok {
		c, err := rpc.DialHTTP("tcp", addr+":"+strconv.Itoa(context.RpcPort()))
		if err != nil {
			clientMu.Unlock()
			return errors.By(err, "failed to connect to %v", addr)
		}

		client = c
		clientMap[addr] = client

		logger.Info(pkg, "succeed to connect to %v", addr)
	}

	clientMu.Unlock()

	if err := client.Call(serviceMethod, req, res); err != nil {
		return errors.By(err, "failed to call %v", serviceMethod)
	}

	return nil
}
