package rpc

import (
	"github.com/tmyksj/rlso11n/pkg/common/object"
	"github.com/tmyksj/rlso11n/pkg/core"
)

type ObjDocker struct {
	ctx  *core.Context
	node *object.Node
}

type ReqDockerRun struct {
	Args []string
}

type ResDockerRun struct {
	Stdout string
}

func Docker(ctx *core.Context, node *object.Node) *ObjDocker {
	return &ObjDocker{
		ctx:  ctx,
		node: node,
	}
}

func (obj *ObjDocker) Run(req *ReqDockerRun) (*ResDockerRun, error) {
	res := &ResDockerRun{}
	err := call(obj.ctx, obj.node, "Rpc.DockerRun", req, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
