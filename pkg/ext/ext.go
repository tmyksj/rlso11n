package ext

import (
	"github.com/tmyksj/rlso11n/app/logger"
	"github.com/tmyksj/rlso11n/pkg/context"
	"os/exec"
	"strings"
)

func Run(args ...string) (string, error) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = context.Env()

	b, err := cmd.Output()
	if err != nil {
		logger.Error(pkg, "failed to run command, %v", err)
		return "", err
	}

	r := strings.TrimSpace(string(b))

	logger.Info(pkg, "succeed to run command")

	return r, nil
}
