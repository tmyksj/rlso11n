package rpc

import (
	"github.com/tmyksj/rootless-orchestration/context"
	"github.com/tmyksj/rootless-orchestration/logger"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
	"sync"
)

type Rpc int

var clientMap = make(map[string]*rpc.Client)
var clientMu sync.Mutex

var serverIsRunning = false
var serverMu sync.Mutex

func Call(host string, serviceMethod string, req interface{}, res interface{}) error {
	clientMu.Lock()

	client, ok := clientMap[host]
	if !ok {
		c, err := rpc.DialHTTP("tcp", host+":"+strconv.Itoa(context.RpcPort()))
		if err != nil {
			logger.Errorf("rpc", "fail to connect to %v", host)
			logger.Errorf("rpc", "%v", err)

			clientMu.Unlock()
			return err
		}

		client = c
		clientMap[host] = client

		logger.Infof("rpc", "succeed to connect to %v", host)
	}

	clientMu.Unlock()

	err := client.Call(serviceMethod, req, res)
	if err != nil {
		logger.Errorf("rpc", "fail to call %v", serviceMethod)
		logger.Errorf("rpc", "%v", err)
		return err
	}

	logger.Infof("rpc", "succeed to call %v", serviceMethod)

	return nil
}

func Serve() error {
	serverMu.Lock()
	defer serverMu.Unlock()

	if serverIsRunning {
		return nil
	}

	s := new(Rpc)
	if err := rpc.Register(s); err != nil {
		logger.Errorf("rpc", "fail to register rpc")
		logger.Errorf("rpc", "%v", err)
		return err
	}

	rpc.HandleHTTP()

	l, err := net.Listen("tcp", ":"+strconv.Itoa(context.RpcPort()))
	if err != nil {
		logger.Errorf("rpc", "fail to listen %v/tcp", context.RpcPort())
		logger.Errorf("rpc", "%v", err)
		return err
	}

	go func() {
		err := http.Serve(l, nil)
		logger.Infof("rpc", "served")
		logger.Infof("rpc", "%v", err)
	}()

	serverIsRunning = true
	logger.Infof("rpc", "succeed to start server")

	return nil
}
