package exec

import (
	"github.com/tmyksj/rlso11n/pkg/common/errors"
	"github.com/tmyksj/rlso11n/pkg/common/logger"
	"github.com/tmyksj/rlso11n/pkg/core"
	"os/exec"
	"strings"
)

type ObjAny struct {
	ctx *core.Context
}

type ReqAnyRun struct {
	Name string
	Args []string
}

type ResAnyRun struct {
	Stdout string
}

func Any(ctx *core.Context) *ObjAny {
	return &ObjAny{
		ctx: ctx,
	}
}

func (obj *ObjAny) Run(req *ReqAnyRun) (*ResAnyRun, error) {
	cmd := exec.Command(req.Name, req.Args...)
	cmd.Env = obj.ctx.Env()

	b, err := cmd.Output()
	if err != nil {
		return nil, errors.By(err, "failed to run command")
	}

	stdout := strings.TrimSpace(string(b))

	logger.Info(pkg, "succeed to run command")

	return &ResAnyRun{
		Stdout: stdout,
	}, nil
}
