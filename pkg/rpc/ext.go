package rpc

import (
	"github.com/tmyksj/rlso11n/pkg/context"
	"github.com/tmyksj/rlso11n/pkg/ext"
)

const MtdExtRun = "Rpc.ExtRun"

type ReqExtRun struct {
	Args []string
}

type ResExtRun struct {
	Stdout string
}

func (r *Rpc) ExtRun(req *ReqExtRun, res *ResExtRun) error {
	if err := context.ReadyOrError(); err != nil {
		return err
	}

	o, err := ext.Run(req.Args...)
	if err != nil {
		return err
	}

	res.Stdout = o

	return nil
}
