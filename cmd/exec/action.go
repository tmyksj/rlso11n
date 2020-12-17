package exec

import (
	"fmt"
	"github.com/tmyksj/rlso11n/cmd"
	"github.com/tmyksj/rlso11n/pkg/common/errors"
	"github.com/tmyksj/rlso11n/pkg/core"
	"github.com/tmyksj/rlso11n/pkg/loader"
	"github.com/tmyksj/rlso11n/pkg/service/exec"
	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {
	return cmd.Perform(func() error {
		cDocker := c.Bool("docker")
		cNodes := c.String("nodes")
		cServer := c.String("server")
		cArgs := c.Args().Slice()

		if len(cArgs) == 0 {
			return errors.By(nil, "bad request: arguments are required")
		}

		if err := loader.LoadAsClient(cServer); err != nil {
			return err
		}

		ctx := core.GetContext()
		res, err := exec.Action(ctx).Run(&exec.ReqActionRun{
			Docker: cDocker,
			Nodes:  cNodes,
			Args:   cArgs,
		})

		if err != nil {
			return err
		}

		for _, output := range res.Output {
			fmt.Println("[" + output.Node.Name + "]")
			fmt.Println(output.Stdout)
		}

		return nil
	})
}
