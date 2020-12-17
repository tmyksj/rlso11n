package core

import (
	"github.com/tmyksj/rlso11n/pkg/common/errors"
	"github.com/tmyksj/rlso11n/pkg/common/logger"
	"github.com/tmyksj/rlso11n/pkg/common/object"
	"net"
	"os"
	"os/exec"
	"strings"
)

type Constructor struct {
	Dir         string
	ManagerNode *object.Node
	NodeList    []*object.Node
}

func (constructor *Constructor) Load() error {
	myNode, err := constructor.findMyNode(constructor.NodeList)
	if err != nil {
		return errors.By(err, "failed to load")
	}

	context.dir = constructor.Dir
	logger.Info(pkg, "dir = %v", context.dir)

	context.managerNode = constructor.ManagerNode
	logger.Info(pkg, "manager node = %v", constructor.toStringNode(context.managerNode))

	context.myNode = myNode
	logger.Info(pkg, "my node = %v", constructor.toStringNode(context.myNode))

	context.nodeList = constructor.NodeList
	logger.Info(pkg, "node list = %v", constructor.toStringNodeList(context.nodeList))

	return nil
}

func (constructor *Constructor) Setup() error {
	if err := constructor.setupDir(); err != nil {
		return errors.By(err, "failed to setup")
	}

	constructor.setupDirFinalizer()

	return nil
}

func (constructor *Constructor) ToReady() error {
	context.ready = true
	return nil
}

func (constructor *Constructor) findMyNode(nodeList []*object.Node) (*object.Node, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		logger.Warn(pkg, "failed to get interfaces, %v", err)
		return nil, nil
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

			for _, node := range nodeList {
				if ip == node.Addr {
					return node, nil
				}
			}
		}
	}

	return nil, nil
}

func (constructor *Constructor) setupDir() error {
	if err := os.MkdirAll(context.Dir(), 0755); err != nil {
		return errors.By(err, "failed to make directory")
	}

	logger.Info(pkg, "succeed to make directory")
	return nil
}

func (constructor *Constructor) setupDirFinalizer() {
	context.AddFinalizer(func() {
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

func (constructor *Constructor) toStringNode(node *object.Node) string {
	if node == nil {
		return "<nil>"
	}

	return "{" + node.Name + " " + node.Addr + "}"
}

func (constructor *Constructor) toStringNodeList(nodeList []*object.Node) string {
	ret := make([]string, len(nodeList))
	for i, node := range nodeList {
		ret[i] = constructor.toStringNode(node)
	}

	return "[" + strings.Join(ret, " ") + "]"
}
