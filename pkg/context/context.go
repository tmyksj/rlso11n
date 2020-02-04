package context

import (
	"errors"
	"github.com/tmyksj/rlso11n/app/logger"
	"golang.org/x/crypto/ssh"
	"os"
)

var addr string

func Addr() string {
	return addr
}

func SetAddr(val string) {
	addr = val
	logger.Infof("pkg/context", "addr = %v", addr)
}

var dir string

func Dir() string {
	return dir
}

func SetDir(val string) {
	dir = val
	logger.Infof("pkg/context", "dir = %v", dir)
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
	logger.Infof("pkg/context", "host list = %v", hostList)
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

var sshAuthMethod []ssh.AuthMethod

func SshAuthMethod() []ssh.AuthMethod {
	return sshAuthMethod
}

func SetSshAuthMethod(val []ssh.AuthMethod) {
	sshAuthMethod = val
	logger.Infof("pkg/context", "ssh auth method = #%v", len(sshAuthMethod))
}

var sshUsername string

func SshUsername() string {
	return sshUsername
}

func SetSshUsername(val string) {
	sshUsername = val
	logger.Infof("pkg/context", "ssh username = %v", sshUsername)
}

var starterAddr string

func StarterAddr() string {
	return starterAddr
}

func SetStarterAddr(val string) {
	starterAddr = val
	logger.Infof("pkg/context", "starter addr = %v", starterAddr)
}