package ext

import (
	"github.com/tmyksj/rlso11n/app/logger"
	"github.com/tmyksj/rlso11n/pkg/context"
	"github.com/tmyksj/rlso11n/pkg/errors"
	"os/exec"
	"strings"
)

func Run(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Env = context.Env()

	b, err := cmd.Output()
	if err != nil {
		return "", errors.By(err, "failed to run command")
	}

	r := strings.TrimSpace(string(b))

	logger.Info(pkg, "succeed to run command")

	return r, nil
}
