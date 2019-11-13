package rpc

import "github.com/tmyksj/rootless-orchestration/context"

const MtdContextSync = "Rpc.ContextSync"

type ReqContextSync struct {
	Dir         string
	HostList    string
	ManagerAddr string
}

type ResContextSync struct {
}

func (r *Rpc) ContextSync(req *ReqContextSync, res *ResContextSync) error {
	context.Init(req.Dir, req.HostList, req.ManagerAddr)
	return nil
}
