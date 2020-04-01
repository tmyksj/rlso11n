package loader

import (
	"github.com/tmyksj/rlso11n/app/logger"
	"github.com/tmyksj/rlso11n/pkg/context/loader_proxy"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

const (
	envPrefix      = "RLSO11N_"
	envDir         = envPrefix + "DIR"
	envHostList    = envPrefix + "HOST_LIST"
)

func LoadAsCommander() {
	Pull("0.0.0.0")
	logger.Infof("pkg/context/loader", "set context as commander from existing rpc server")
}

func LoadAsStarter() {
	loader_proxy.Load(&loader_proxy.LoadReq{
		Dir:           parseDir(os.Getenv(envDir)),
		HostList:      parseHostList(os.Getenv(envHostList)),
		StarterAddr:   loader_proxy.CurrentAddr,
	}, &loader_proxy.SetupReq{
		Dir: false,
	})

	logger.Infof("pkg/context/loader", "set context as starter from environment")
}

func parseDir(val string) string {
	if val != "" {
		return val
	} else {
		return "/tmp/rlso11n-" + strconv.Itoa(rand.Int())
	}
}

func parseHostList(val string) []string {
	hostList := make([]string, 0)
	for _, host := range strings.Split(val, ",") {
		trimmed := strings.TrimSpace(host)
		if trimmed != "" {
			hostList = append(hostList, trimmed)
		}
	}

	return hostList
}
