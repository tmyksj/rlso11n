package rpc

import (
	"github.com/tmyksj/rlso11n/pkg/common/errors"
	"github.com/tmyksj/rlso11n/pkg/common/logger"
	"github.com/tmyksj/rlso11n/pkg/core"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
	"sync"
)

var serverIsRunning = false
var serverMu sync.Mutex

type ObjServer struct {
	ctx *core.Context
}

type ReqServerServe struct {
}

type ResServerServe struct {
}

func Server(ctx *core.Context) *ObjServer {
	return &ObjServer{
		ctx: ctx,
	}
}

func (obj *ObjServer) Serve(req *ReqServerServe) (*ResServerServe, error) {
	serverMu.Lock()
	defer serverMu.Unlock()

	if serverIsRunning {
		return &ResServerServe{}, nil
	}

	s := new(Rpc)
	if err := rpc.Register(s); err != nil {
		return nil, errors.By(err, "failed to register rpc")
	}

	rpc.HandleHTTP()

	l, err := net.Listen("tcp", ":"+strconv.Itoa(obj.ctx.RpcPort()))
	if err != nil {
		return nil, errors.By(err, "failed to listen %v/tcp", obj.ctx.RpcPort())
	}

	go func() {
		err := http.Serve(l, nil)
		logger.Info(pkg, "served, %v", err)
	}()

	serverIsRunning = true
	logger.Info(pkg, "succeed to start server")

	return &ResServerServe{}, nil
}
