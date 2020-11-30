package loader_proxy

import (
	"github.com/tmyksj/rlso11n/app/core"
	"github.com/tmyksj/rlso11n/app/logger"
	"github.com/tmyksj/rlso11n/pkg/context"
	"github.com/tmyksj/rlso11n/pkg/errors"
	"net"
	"os"
	"os/exec"
)

type LoadReq struct {
	Dir         string
	HostList    []context.Host
	ManagerAddr string
}

type SetupReq struct {
	Dir bool
}

func Load(loadReq *LoadReq, setupReq *SetupReq) error {
	if err := load(loadReq); err != nil {
		return errors.By(err, "failed to load")
	}

	if err := setup(setupReq); err != nil {
		return errors.By(err, "failed to setup")
	}

	context.SetReady(true)

	return nil
}

func load(req *LoadReq) error {
	myAddr, err := findMyAddr(req.HostList)
	if err != nil {
		return err
	}

	context.SetDir(req.Dir)
	context.SetHostList(req.HostList)
	context.SetManagerAddr(req.ManagerAddr)
	context.SetMyAddr(myAddr)

	return nil
}

func setup(req *SetupReq) error {
	if req.Dir {
		if err := setupDir(); err != nil {
			return err
		}

		setupDirFinalizer()
	}

	return nil
}

func setupDir() error {
	if err := os.MkdirAll(context.Dir(), 0755); err != nil {
		return errors.By(err, "failed to make directory")
	}

	logger.Info(pkg, "succeed to make directory")
	return nil
}

func setupDirFinalizer() {
	core.RegisterFinalizer(func() {
		if _, err := os.Stat(context.Dir()); os.IsNotExist(err) {
			return
		}

		cmd := exec.Command("rootlesskit",
			"--state-dir="+context.Dir()+"/rootlesskit",
			"--net=host",
			"/bin/sh", "-c", "rm -rf $XDG_CONFIG_HOME $XDG_CACHE_HOME $XDG_DATA_HOME $XDG_RUNTIME_DIR")
		cmd.Env = context.Env()
		if err := cmd.Run(); err != nil {
			logger.Error(pkg, "failed to cleanup directory, %v", err)
		}

		if err := os.RemoveAll(context.Dir()); err != nil {
			logger.Error(pkg, "failed to remove directory, %v", err)
		} else {
			logger.Info(pkg, "succeed to remove directory")
		}
	})
}

func findMyAddr(hostList []context.Host) (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", errors.By(err, "failed to get interfaces")
	}

	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			logger.Warn(pkg, "failed to get addresses, %v", err)
			continue
		}

		for _, addr := range addrs {
			var ip string
			switch v := addr.(type) {
			case *net.IPAddr:
				ip = v.IP.String()
			case *net.IPNet:
				ip = v.IP.String()
			}

			for _, h := range hostList {
				if ip == h.Addr {
					return ip, nil
				}
			}
		}
	}

	return "", errors.By(nil, "failed to get my address")
}
