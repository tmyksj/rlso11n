package docker

import (
	"github.com/tmyksj/rlso11n/app/logger"
	"github.com/tmyksj/rlso11n/pkg/context"
	"os/exec"
	"strings"
)

func Run(args ...string) (string, error) {
	cmd := exec.Command("docker", append([]string{"-H", context.Addr()}, args...)...)
	cmd.Env = context.Env()

	b, err := cmd.Output()
	if err != nil {
		logger.Errorf("pkg/ext/docker", "fail to run docker command, %v", err)
		return "", err
	}

	r := strings.TrimSpace(string(b))

	logger.Infof("pkg/ext/docker", "succeed to run docker command")

	return r, nil
}
