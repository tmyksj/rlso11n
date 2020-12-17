package serve_rpc

import (
	"github.com/tmyksj/rlso11n/cmd"
	"github.com/tmyksj/rlso11n/pkg/core"
	"github.com/tmyksj/rlso11n/pkg/loader"
	"github.com/tmyksj/rlso11n/pkg/service/serve_rpc"
	"github.com/urfave/cli/v2"
)

func Action(_ *cli.Context) error {
	return cmd.Perform(func() error {
		if err := loader.LoadAsWorker(); err != nil {
			return err
		}

		ctx := core.GetContext()
		_, err := serve_rpc.Action(ctx).Run(&serve_rpc.ReqActionRun{})

		if err != nil {
			return err
		}

		return nil
	})
}
