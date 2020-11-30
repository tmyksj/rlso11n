package context

import (
	"github.com/tmyksj/rlso11n/app/logger"
	"github.com/tmyksj/rlso11n/pkg/errors"
	"os"
)

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

type Host struct {
	Addr string
	Name string
}

var hostList []Host

func HostList() []Host {
	return hostList
}

func SetHostList(val []Host) {
	hostList = val
	logger.Info(pkg, "hosts = %v", hostList)
}

var managerAddr string

func ManagerAddr() string {
	return managerAddr
}

func SetManagerAddr(val string) {
	managerAddr = val
	logger.Info(pkg, "manager addr = %v", managerAddr)
}

var myAddr string

func MyAddr() string {
	return myAddr
}

func SetMyAddr(val string) {
	myAddr = val
	logger.Info(pkg, "my addr = %v", myAddr)
}

var ready = false

func ReadyOrError() error {
	if ready {
		return nil
	} else {
		return errors.By(nil, "context is not ready")
	}
}

func SetReady(val bool) {
	ready = val
}

func RpcPort() int {
	return 50128
}
