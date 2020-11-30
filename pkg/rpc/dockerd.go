package rpc

import (
	"github.com/tmyksj/rlso11n/pkg/context"
	"github.com/tmyksj/rlso11n/pkg/ext/dockerd"
)

const MtdDockerdStart = "Rpc.DockerdStart"

type ReqDockerdStart struct {
}

type ResDockerdStart struct {
}

func (r *Rpc) DockerdStart(req *ReqDockerdStart, res *ResDockerdStart) error {
	return do(func() error {
		if err := context.ReadyOrError(); err != nil {
			return err
		}

		return dockerd.Start()
	})
}
