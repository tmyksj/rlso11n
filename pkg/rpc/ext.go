package rpc

import (
	"github.com/tmyksj/rlso11n/pkg/context"
	"github.com/tmyksj/rlso11n/pkg/ext"
)

const MtdExtRun = "Rpc.ExtRun"

type ReqExtRun struct {
	Args []string
	Name string
}

type ResExtRun struct {
	Stdout string
}

func (r *Rpc) ExtRun(req *ReqExtRun, res *ResExtRun) error {
	return do(func() error {
		if err := context.ReadyOrError(); err != nil {
			return err
		}

		o, err := ext.Run(req.Name, req.Args...)
		if err != nil {
			return err
		}

		res.Stdout = o

		return nil
	})
}
