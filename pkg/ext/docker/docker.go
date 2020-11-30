package docker

import (
	"github.com/tmyksj/rlso11n/app/logger"
	"github.com/tmyksj/rlso11n/pkg/context"
	"github.com/tmyksj/rlso11n/pkg/errors"
	"os/exec"
	"strings"
)

func Run(args ...string) (string, error) {
	cmd := exec.Command("docker",
		append([]string{"--host", "unix://" + context.DockerSock()}, args...)...)
	cmd.Env = context.Env()

	b, err := cmd.Output()
	if err != nil {
		return "", errors.By(err, "failed to run docker command")
	}

	r := strings.TrimSpace(string(b))

	logger.Info(pkg, "succeed to run docker command")

	return r, nil
}
