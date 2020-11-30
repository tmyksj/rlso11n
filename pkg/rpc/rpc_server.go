package rpc

import (
	"github.com/tmyksj/rlso11n/app/logger"
	"github.com/tmyksj/rlso11n/pkg/context"
	"github.com/tmyksj/rlso11n/pkg/errors"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
	"sync"
)

var serverIsRunning = false
var serverMu sync.Mutex

func Serve() error {
	serverMu.Lock()
	defer serverMu.Unlock()

	if serverIsRunning {
		return nil
	}

	s := new(Rpc)
	if err := rpc.Register(s); err != nil {
		return errors.By(err, "failed to register rpc")
	}

	rpc.HandleHTTP()

	l, err := net.Listen("tcp", ":"+strconv.Itoa(context.RpcPort()))
	if err != nil {
		return errors.By(err, "failed to listen %v/tcp", context.RpcPort())
	}

	go func() {
		err := http.Serve(l, nil)
		logger.Info(pkg, "served, %v", err)
	}()

	serverIsRunning = true
	logger.Info(pkg, "succeed to start server")

	return nil
}
