package loader

import (
	"github.com/tmyksj/rlso11n/app/logger"
	"github.com/tmyksj/rlso11n/pkg/context/loader_proxy"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"math/rand"
	"os"
	"os/user"
	"strconv"
	"strings"
)

const (
	envPrefix      = "RLSO11N_"
	envDir         = envPrefix + "DIR"
	envHostList    = envPrefix + "HOST_LIST"
	envSshUsername = envPrefix + "SSH_USERNAME"
	envSshIdentity = envPrefix + "SSH_IDENTITY"
)

func LoadAsCommander() {
	Pull("0.0.0.0")
	logger.Infof("pkg/context/loader", "set context as commander from existing rpc server")
}

func LoadAsStarter() {
	loader_proxy.Load(&loader_proxy.LoadReq{
		Dir:           parseDir(os.Getenv(envDir)),
		HostList:      parseHostList(os.Getenv(envHostList)),
		SshAuthMethod: parseSshAuthMethod(os.Getenv(envSshIdentity)),
		SshUsername:   parseSshUsername(os.Getenv(envSshUsername)),
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

func parseSshAuthMethod(val string) []ssh.AuthMethod {
	var sshAuthMethod []ssh.AuthMethod
	if val != "" {
		key, err := ioutil.ReadFile(val)
		if err != nil {
			logger.Errorf("pkg/context", "fail to read identity file (to use ssh), %v", err)
			goto initializedSshAuth
		}

		identity, err := ssh.ParsePrivateKey(key)
		if err != nil {
			logger.Errorf("pkg/context", "fail to parse private key, %v", err)
			goto initializedSshAuth
		}

		sshAuthMethod = []ssh.AuthMethod{
			ssh.PublicKeys(identity),
		}
	}

initializedSshAuth:
	return sshAuthMethod
}

func parseSshUsername(val string) string {
	if val != "" {
		return val
	} else if c, err := user.Current(); err == nil {
		return c.Username
	} else {
		logger.Errorf("pkg/context/loader", "fail to get username (to use ssh), %v", err)
		return ""
	}
}
