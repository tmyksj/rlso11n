package context

import (
	"github.com/tmyksj/rootless-orchestration/core"
	"github.com/tmyksj/rootless-orchestration/logger"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
)

const (
	envPrefix      = "ROOTLESS_ORCHESTRATION_"
	envDir         = envPrefix + "DIR"
	envHosts       = envPrefix + "HOSTS"
	envManagerAddr = envPrefix + "MANAGER_ADDR"
	envSshUsername = envPrefix + "SSH_USERNAME"
	envSshIdentity = envPrefix + "SSH_IDENTITY"
)

var addr string

var dir string

var hostList []string

var managerAddr string

var sshUsername string

var sshAuthMethod []ssh.AuthMethod

func Addr() string {
	return addr
}

func Dir() string {
	return dir
}

func Env() []string {
	return append(os.Environ(),
		"DOCKER_HOST=unix://"+os.Getenv("XDG_RUNTIME_DIR")+"/docker.sock",
		"XDG_CONFIG_HOME="+Dir()+"/config",
		"XDG_CACHE_HOME="+Dir()+"/cache",
		"XDG_DATA_HOME="+Dir()+"/home")
}

func HostList() []string {
	return hostList
}

func ManagerAddr() string {
	return managerAddr
}

func RpcPort() int {
	return 50128
}

func SshUsername() string {
	return sshUsername
}

func SshAuthMethod() []ssh.AuthMethod {
	return sshAuthMethod
}

func Init(d string, hl string, ma string) {
	initDir(d)
	initAddr(hl)
	initManagerAddr(ma)
	logger.Infof("context", "set context")
}

func InitFromEnv() {
	initDir(os.Getenv(envDir))
	initAddr(os.Getenv(envHosts))
	initManagerAddr(os.Getenv(envManagerAddr))
	initSsh(os.Getenv(envSshUsername), os.Getenv(envSshIdentity))
	logger.Infof("context", "set context from environment")
}

func initDir(v string) {
	if v != "" {
		dir = v
	} else {
		dir = "/tmp/rootless-orchestration-" + strconv.Itoa(rand.Int())
	}

	logger.Infof("context", "dir = %v", dir)

	if err := os.MkdirAll(dir, 0755); err != nil {
		logger.Errorf("context", "fail to make directory")
		logger.Errorf("context", "%v", err)
	} else {
		logger.Infof("context", "succeed to make directory")
	}

	core.Finalize(func() {
		if _, err := os.Stat(Dir()); os.IsNotExist(err) {
			return
		}

		script := Dir() + "/clean.sh"
		err := ioutil.WriteFile(script, []byte(`#!/bin/sh
set -e -x
if [ -z $_DOCKERD_ROOTLESS_CHILD ]; then
	_DOCKERD_ROOTLESS_CHILD=1
	export _DOCKERD_ROOTLESS_CHILD
	exec rootlesskit \
		--state-dir=`+Dir()+`/rootlesskit \
		--net=host \
		$0 $@
else
	[ $_DOCKERD_ROOTLESS_CHILD = 1 ]
	exec rm -rf $XDG_CONFIG_HOME $XDG_CACHE_HOME $XDG_DATA_HOME
fi
`), 0755)

		if err != nil {
			logger.Errorf("context", "fail to create cleaning script")
			logger.Errorf("context", "%v", err)
		} else {
			c := exec.Command(script)
			c.Env = Env()

			if err := c.Run(); err != nil {
				logger.Errorf("context", "fail to run cleaning script")
				logger.Errorf("context", "%v", err)
			}
		}

		if err := os.RemoveAll(Dir()); err != nil {
			logger.Errorf("context", "fail to remove directory")
			logger.Errorf("context", "%v", err)
		} else {
			logger.Infof("context", "succeed to remove directory")
		}
	})
}

func initAddr(v string) {
	hostList = make([]string, 0)
	for _, host := range strings.Split(v, ",") {
		if host != "" {
			hostList = append(hostList, host)
		}
	}

	if ifaces, err := net.Interfaces(); err != nil {
		logger.Errorf("context", "fail to get interfaces")
		logger.Errorf("context", "%v", err)
	} else {
		for _, i := range ifaces {
			if addrs, err := i.Addrs(); err != nil {
				logger.Errorf("context", "fail to get addrs")
				logger.Errorf("context", "%v", err)
			} else {
				for _, ad := range addrs {
					var ip string
					switch v := ad.(type) {
					case *net.IPAddr:
						ip = v.IP.String()
					case *net.IPNet:
						ip = v.IP.String()
					}

					for _, h := range hostList {
						if ip == h {
							addr = ip
							logger.Infof("context", "addr = %v", addr)
						}
					}
				}
			}
		}
	}

	logger.Infof("context", "hosts = %v", hostList)
}

func initManagerAddr(v string) {
	if v != "" {
		managerAddr = v
	} else {
		managerAddr = Addr()
	}

	logger.Infof("context", "manager addr = %v", managerAddr)
}

func initSsh(un string, id string) {
	if un != "" {
		sshUsername = un
	} else if c, err := user.Current(); err == nil {
		sshUsername = c.Username
	} else {
		logger.Errorf("context", "fail to get username (to use ssh)")
		logger.Errorf("context", "%v", err)
		sshUsername = ""
	}

	logger.Infof("context", "ssh username = %v", sshUsername)

	sshAuthMethod = []ssh.AuthMethod{}
	if id != "" {
		key, err := ioutil.ReadFile(id)
		if err != nil {
			logger.Errorf("context", "fail to read identity file (to use ssh)")
			logger.Errorf("context", "%v", err)
			goto initializedSshAuth
		}

		identity, err := ssh.ParsePrivateKey(key)
		if err != nil {
			logger.Errorf("cmd/bg", "fail to parse private key")
			logger.Errorf("cmd/bg", "%v", err)
			goto initializedSshAuth
		}

		sshAuthMethod = []ssh.AuthMethod{
			ssh.PublicKeys(identity),
		}
	}

initializedSshAuth:
	logger.Infof("context", "ssh auth method = #%v", len(sshAuthMethod))
}
