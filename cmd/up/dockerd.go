package up

import (
	"github.com/tmyksj/rlso11n/cmd"
	"github.com/tmyksj/rlso11n/pkg/core"
	"github.com/tmyksj/rlso11n/pkg/loader"
	"github.com/tmyksj/rlso11n/pkg/service/up"
	"github.com/urfave/cli/v2"
)

func Dockerd(c *cli.Context) error {
	return cmd.Perform(func() error {
		cServer := c.String("server")

		if err := loader.LoadAsClient(cServer); err != nil {
			return err
		}

		ctx := core.GetContext()
		_, err := up.Dockerd(ctx).Run(&up.ReqDockerdRun{})

		if err != nil {
			return err
		}

		return nil
	})
}
