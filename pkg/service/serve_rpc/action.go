package serve_rpc

import (
	"github.com/tmyksj/rlso11n/pkg/common/utility"
	"github.com/tmyksj/rlso11n/pkg/component/rpc"
	"github.com/tmyksj/rlso11n/pkg/core"
)

type ObjAction struct {
	ctx *core.Context
}

type ReqActionRun struct {
}

type ResActionRun struct {
}

func Action(ctx *core.Context) *ObjAction {
	return &ObjAction{
		ctx: ctx,
	}
}

func (obj *ObjAction) Run(req *ReqActionRun) (*ResActionRun, error) {
	_, err := rpc.Server(obj.ctx).Serve(&rpc.ReqServerServe{})
	if err != nil {
		return nil, err
	}

	utility.WaitInterrupt()

	return &ResActionRun{}, nil
}
