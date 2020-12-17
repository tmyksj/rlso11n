package loader

import (
	"github.com/tmyksj/rlso11n/pkg/common/errors"
	"github.com/tmyksj/rlso11n/pkg/common/logger"
	"github.com/tmyksj/rlso11n/pkg/common/object"
	"github.com/tmyksj/rlso11n/pkg/common/utility"
	"github.com/tmyksj/rlso11n/pkg/component/rpc"
	"github.com/tmyksj/rlso11n/pkg/core"
	"math/rand"
	"net"
	"strconv"
)

func LoadAsClient(host string) error {
	var addr string
	if net.ParseIP(host) == nil {
		a, err := lookup(host)
		if err != nil {
			return errors.By(err, "failed to load context as client")
		}

		addr = a
	} else {
		addr = host
	}

	err := pull(&object.Node{
		Addr: addr,
		Name: host,
	})

	if err != nil {
		return errors.By(err, "failed to load context as client")
	}

	logger.Info(pkg, "set context as client")
	return nil
}

func LoadAsManager(dir string, hosts []string) error {
	parsedDir := parseDir(dir)
	parsedHostList, err := parseHostList(hosts)
	if err != nil {
		return errors.By(err, "failed to load context as manager")
	}

	constructor := core.Constructor{
		Dir:         parsedDir,
		ManagerNode: decideManagerNode(parsedHostList),
		NodeList:    parsedHostList,
	}

	if err := constructor.Load(); err != nil {
		return errors.By(err, "failed to load context as manager")
	}

	if err := constructor.ToReady(); err != nil {
		return errors.By(err, "failed to load context as manager")
	}

	logger.Info(pkg, "set context as manager")
	return nil
}

func LoadAsWorker() error {
	return nil
}

func decideManagerNode(nodeList []*object.Node) *object.Node {
	ifaces, err := net.Interfaces()
	if err != nil {
		logger.Warn(pkg, "failed to get interfaces, %v", err)
		return nodeList[0]
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
					return node
				}
			}
		}
	}

	return nodeList[0]
}

func lookup(host string) (string, error) {
	ips, err := net.LookupIP(host)
	if err != nil {
		return "", errors.By(err, "failed to lookup %v", host)
	}

	for _, ip := range ips {
		if !ip.IsLoopback() {
			return ip.String(), nil
		}
	}

	ifaces, err := net.Interfaces()
	if err != nil {
		logger.Warn(pkg, "failed to get interfaces, %v", err)
		return "", errors.By(err, "failed to lookup %v", host)
	}

	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			logger.Warn(pkg, "failed to get addresses, %v", err)
			return "", errors.By(err, "failed to lookup %v", host)
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPAddr:
				ip = v.IP
			case *net.IPNet:
				ip = v.IP
			}

			if ip != nil && !ip.IsLoopback() {
				ipStr := ip.String()

				logger.Warn(pkg, "failed to lookup %v, use %v", host, ipStr)
				return ipStr, nil
			}
		}
	}

	return "", errors.By(err, "failed to lookup %v", host)
}

func parseDir(val string) string {
	if val != "" {
		return val
	} else {
		return "/tmp/rlso11n-" + strconv.Itoa(rand.Int())
	}
}

func parseHostList(val []string) ([]*object.Node, error) {
	nodeList := make([]*object.Node, 0)
	for _, host := range val {
		if net.ParseIP(host) == nil {
			addr, err := lookup(host)
			if err != nil {
				return nil, err
			}

			nodeList = append(nodeList, &object.Node{
				Addr: addr,
				Name: host,
			})
		} else {
			nodeList = append(nodeList, &object.Node{
				Addr: host,
				Name: host,
			})
		}
	}

	return nodeList, nil
}

func pull(node *object.Node) error {
	ctx := core.GetContext()

	if err := utility.WaitListenTcp(node.Addr + ":" + strconv.Itoa(ctx.RpcPort())); err != nil {
		return errors.By(err, "failed to pull context from %v", node.Name)
	}

	var res *rpc.ResContextPull
	err := utility.Try(func() error {
		ret, err := rpc.Context(ctx, node).Pull(&rpc.ReqContextPull{})
		if err != nil {
			return err
		}

		res = ret

		return nil
	})

	if err != nil {
		return errors.By(err, "failed to pull context from %v", node.Name)
	}

	constructor := core.Constructor{
		Dir:         res.Dir,
		ManagerNode: res.ManagerNode,
		NodeList:    res.NodeList,
	}

	if err := constructor.Load(); err != nil {
		return errors.By(err, "failed to pull context from %v", node.Name)
	}

	if err := constructor.ToReady(); err != nil {
		return errors.By(err, "failed to pull context from %v", node.Name)
	}

	return nil
}
