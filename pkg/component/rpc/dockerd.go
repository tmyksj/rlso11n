package rpc

import (
	"github.com/tmyksj/rlso11n/pkg/common/object"
	"github.com/tmyksj/rlso11n/pkg/core"
)

type ObjDockerd struct {
	ctx  *core.Context
	node *object.Node
}

type ReqDockerdStart struct {
}

type ResDockerdStart struct {
}

func Dockerd(ctx *core.Context, node *object.Node) *ObjDockerd {
	return &ObjDockerd{
		ctx:  ctx,
		node: node,
	}
}

func (obj *ObjDockerd) Start(req *ReqDockerdStart) (*ResDockerdStart, error) {
	res := &ResDockerdStart{}
	err := call(obj.ctx, obj.node, "Rpc.DockerdStart", req, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
