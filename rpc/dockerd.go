package rpc

import "github.com/tmyksj/rootless-orchestration/ext/dockerd"

const MtdDockerdStart = "Rpc.DockerdStart"

type ReqDockerdStart struct {
}

type ResDockerdStart struct {
}

func (r *Rpc) DockerdStart(req *ReqDockerdStart, res *ResDockerdStart) error {
	return dockerd.StartRootless()
}
