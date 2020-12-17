package rpc

import (
	"github.com/tmyksj/rlso11n/pkg/common/object"
	"github.com/tmyksj/rlso11n/pkg/core"
)

type ObjAny struct {
	ctx  *core.Context
	node *object.Node
}

type ReqAnyRun struct {
	Name string
	Args []string
}

type ResAnyRun struct {
	Stdout string
}

func Any(ctx *core.Context, node *object.Node) *ObjAny {
	return &ObjAny{
		ctx:  ctx,
		node: node,
	}
}

func (obj *ObjAny) Run(req *ReqAnyRun) (*ResAnyRun, error) {
	res := &ResAnyRun{}
	err := call(obj.ctx, obj.node, "Rpc.AnyRun", req, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
