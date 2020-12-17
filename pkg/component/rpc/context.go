package rpc

import (
	"github.com/tmyksj/rlso11n/pkg/common/object"
	"github.com/tmyksj/rlso11n/pkg/core"
)

type ObjContext struct {
	ctx  *core.Context
	node *object.Node
}

type ReqContextPull struct {
}

type ReqContextPush struct {
	Dir         string
	ManagerNode *object.Node
	NodeList    []*object.Node
}

type ResContextPull struct {
	Dir         string
	ManagerNode *object.Node
	NodeList    []*object.Node
}

type ResContextPush struct {
}

func Context(ctx *core.Context, node *object.Node) *ObjContext {
	return &ObjContext{
		ctx:  ctx,
		node: node,
	}
}

func (obj *ObjContext) Pull(req *ReqContextPull) (*ResContextPull, error) {
	res := &ResContextPull{}
	err := call(obj.ctx, obj.node, "Rpc.ContextPull", req, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (obj *ObjContext) Push(req *ReqContextPush) (*ResContextPush, error) {
	res := &ResContextPush{}
	err := call(obj.ctx, obj.node, "Rpc.ContextPush", req, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
