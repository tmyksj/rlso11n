package context

import (
	"errors"
	"github.com/tmyksj/rlso11n/app/logger"
	"os"
)

var addr string

func Addr() string {
	return addr
}

func SetAddr(val string) {
	addr = val
	logger.Info(pkg, "addr = %v", addr)
}

var dir string

func Dir() string {
	return dir
}

func SetDir(val string) {
	dir = val
	logger.Info(pkg, "dir = %v", dir)
}

func DockerSock() string {
	return Dir() + "/runtime/docker.sock"
}

func Env() []string {
	return append(os.Environ(),
		"XDG_CACHE_HOME="+Dir()+"/cache",
		"XDG_CONFIG_HOME="+Dir()+"/config",
		"XDG_DATA_HOME="+Dir()+"/data",
		"XDG_RUNTIME_DIR="+Dir()+"/runtime")
}

var hostList []string

func HostList() []string {
	return hostList
}

func SetHostList(val []string) {
	hostList = val
	logger.Info(pkg, "host list = %v", hostList)
}

var ready = false

func ReadyOrError() error {
	if ready {
		return nil
	} else {
		return errors.New("context is not ready")
	}
}

func SetReady(val bool) {
	ready = val
}

func RpcPort() int {
	return 50128
}

var starterAddr string

func StarterAddr() string {
	return starterAddr
}

func SetStarterAddr(val string) {
	starterAddr = val
	logger.Info(pkg, "starter addr = %v", starterAddr)
}
