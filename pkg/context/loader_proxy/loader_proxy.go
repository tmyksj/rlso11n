package loader_proxy

import (
	"github.com/tmyksj/rlso11n/app/core"
	"github.com/tmyksj/rlso11n/app/logger"
	"github.com/tmyksj/rlso11n/pkg/context"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
	"os/exec"
)

type LoadReq struct {
	Dir           string
	HostList      []string
	SshAuthMethod []ssh.AuthMethod
	SshUsername   string
	StarterAddr   string
}

const CurrentAddr = "CurrentAddr"

type SetupReq struct {
	Dir bool
}

func Load(loadReq *LoadReq, setupReq *SetupReq) {
	load(loadReq)
	setup(setupReq)

	context.SetReady(true)
}

func load(req *LoadReq) {
	context.SetAddr(parseAddr(req.HostList))
	context.SetDir(req.Dir)
	context.SetHostList(req.HostList)
	context.SetSshAuthMethod(req.SshAuthMethod)
	context.SetSshUsername(req.SshUsername)

	if req.StarterAddr == CurrentAddr {
		context.SetStarterAddr(context.Addr())
	} else {
		context.SetStarterAddr(req.StarterAddr)
	}
}

func setup(req *SetupReq) {
	if req.Dir {
		setupDir()
		setupDirFinalizer()
	}
}

func setupDir() {
	if err := os.MkdirAll(context.Dir(), 0755); err != nil {
		logger.Errorf("pkg/context/loader_proxy", "fail to make directory, %v", err)
	} else {
		logger.Infof("pkg/context/loader_proxy", "succeed to make directory")
	}
}

func setupDirFinalizer() {
	core.Finalize(func() {
		if _, err := os.Stat(context.Dir()); os.IsNotExist(err) {
			return
		}

		cmd := exec.Command("rootlesskit",
			"--state-dir="+context.Dir()+"/rootlesskit",
			"--net=host",
			"/bin/sh", "-c", "rm -rf $XDG_CONFIG_HOME $XDG_CACHE_HOME $XDG_DATA_HOME $XDG_RUNTIME_DIR")
		cmd.Env = context.Env()
		if err := cmd.Run(); err != nil {
			logger.Errorf("pkg/context/loader_proxy", "fail to cleanup directory, %v", err)
		}

		if err := os.RemoveAll(context.Dir()); err != nil {
			logger.Errorf("pkg/context/loader_proxy", "fail to remove directory, %v", err)
		} else {
			logger.Infof("pkg/context/loader_proxy", "succeed to remove directory")
		}
	})
}

func parseAddr(hostList []string) string {
	if ifaces, err := net.Interfaces(); err != nil {
		logger.Errorf("pkg/context/loader", "fail to get interfaces, %v", err)
	} else {
		for _, i := range ifaces {
			if addrs, err := i.Addrs(); err != nil {
				logger.Errorf("pkg/context/loader", "fail to get addrs, %v", err)
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
							return ip
						}
					}
				}
			}
		}
	}

	return ""
}
