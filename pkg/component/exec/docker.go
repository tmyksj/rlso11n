package exec

import (
	"github.com/tmyksj/rlso11n/pkg/common/errors"
	"github.com/tmyksj/rlso11n/pkg/common/logger"
	"github.com/tmyksj/rlso11n/pkg/core"
	"os/exec"
	"strings"
)

type ObjDocker struct {
	ctx *core.Context
}

type ReqDockerRun struct {
	Args []string
}

type ResDockerRun struct {
	Stdout string
}

func Docker(ctx *core.Context) *ObjDocker {
	return &ObjDocker{
		ctx: ctx,
	}
}

func (obj *ObjDocker) Run(req *ReqDockerRun) (*ResDockerRun, error) {
	cmd := exec.Command("docker", append([]string{"--host", "unix://" + obj.ctx.DockerSock()}, req.Args...)...)
	cmd.Env = obj.ctx.Env()

	b, err := cmd.Output()
	if err != nil {
		return nil, errors.By(err, "failed to run docker command")
	}

	stdout := strings.TrimSpace(string(b))

	logger.Info(pkg, "succeed to run docker command")

	return &ResDockerRun{
		Stdout: stdout,
	}, nil
}
