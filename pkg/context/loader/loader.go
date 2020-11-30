package loader

import (
	"github.com/tmyksj/rlso11n/app/logger"
	"github.com/tmyksj/rlso11n/pkg/context"
	"github.com/tmyksj/rlso11n/pkg/context/loader_proxy"
	"github.com/tmyksj/rlso11n/pkg/errors"
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

	if err := Pull(addr); err != nil {
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

	if err := loader_proxy.Load(&loader_proxy.LoadReq{
		Dir:         parsedDir,
		HostList:    parsedHostList,
		ManagerAddr: decideManagerAddr(parsedHostList),
	}, &loader_proxy.SetupReq{
		Dir: false,
	}); err != nil {
		return errors.By(err, "failed to load context as manager")
	}

	logger.Info(pkg, "set context as manager")
	return nil
}

func parseDir(val string) string {
	if val != "" {
		return val
	} else {
		return "/tmp/rlso11n-" + strconv.Itoa(rand.Int())
	}
}

func parseHostList(val []string) ([]context.Host, error) {
	hostList := make([]context.Host, 0)
	for _, host := range val {
		if net.ParseIP(host) == nil {
			addr, err := lookup(host)
			if err != nil {
				return nil, err
			}

			hostList = append(hostList, context.Host{
				Addr: addr,
				Name: host,
			})
		} else {
			hostList = append(hostList, context.Host{
				Addr: host,
				Name: host,
			})
		}
	}

	return hostList, nil
}

func lookup(host string) (string, error) {
	ips, err := net.LookupIP(host)
	if err != nil {
		return "", err
	}

	for _, ip := range ips {
		if !ip.IsLoopback() {
			return ip.String(), nil
		}
	}

	ifaces, err := net.Interfaces()
	if err != nil {
		logger.Warn(pkg, "failed to lookup %v", host)
		return "", errors.By(err, "failed to get interfaces")
	}

	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			logger.Warn(pkg, "failed to lookup %v", host)
			return "", errors.By(err, "failed to get addresses")
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

	logger.Warn(pkg, "failed to lookup %v", host)
	return "", errors.By(nil, "failed to get address")
}

func decideManagerAddr(hostList []context.Host) string {
	ifaces, err := net.Interfaces()
	if err != nil {
		logger.Warn(pkg, "failed to get interfaces, %v", err)
		return hostList[0].Addr
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
					return ip
				}
			}
		}
	}

	return hostList[0].Addr
}
