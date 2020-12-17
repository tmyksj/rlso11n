package up

import (
	"github.com/tmyksj/rlso11n/pkg/common/logger"
	"github.com/tmyksj/rlso11n/pkg/common/utility"
	"github.com/tmyksj/rlso11n/pkg/component/rpc"
	"github.com/tmyksj/rlso11n/pkg/core"
	"golang.org/x/sync/errgroup"
)

type ObjDockerd struct {
	ctx *core.Context
}

type ReqDockerdRun struct {
}

type ResDockerdRun struct {
}

func Dockerd(ctx *core.Context) *ObjDockerd {
	return &ObjDockerd{
		ctx: ctx,
	}
}

func (obj *ObjDockerd) Run(req *ReqDockerdRun) (*ResDockerdRun, error) {
	eg := errgroup.Group{}

	for _, node := range obj.ctx.NodeList() {
		node := node
		eg.Go(func() error {
			err := utility.Try(func() error {
				_, err := rpc.Dockerd(obj.ctx, node).Start(&rpc.ReqDockerdStart{})

				if err != nil {
					return err
				}

				return nil
			})

			if err != nil {
				logger.Error(pkg, "failed to start dockerd at %v, %v", node.Name, err)
				return err
			}

			logger.Info(pkg, "started dockerd at %v", node.Name)
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return &ResDockerdRun{}, nil
}
