package rpc

import (
	"github.com/tmyksj/rlso11n/pkg/context"
	"github.com/tmyksj/rlso11n/pkg/ext/docker"
)

const MtdDockerRun = "Rpc.DockerRun"

type ReqDockerRun struct {
	Args []string
}

type ResDockerRun struct {
	Output string
}

func (r *Rpc) DockerRun(req *ReqDockerRun, res *ResDockerRun) error {
	return do(func() error {
		if err := context.ReadyOrError(); err != nil {
			return err
		}

		o, err := docker.Run(req.Args...)
		if err != nil {
			return err
		}

		res.Output = o

		return nil
	})
}
