package main

import (
	"fmt"
	"github.com/tmyksj/rlso11n/app/core"
	"github.com/tmyksj/rlso11n/cmd"
	"github.com/tmyksj/rlso11n/cmd/bg"
	"github.com/tmyksj/rlso11n/cmd/up"
	"github.com/urfave/cli"
	"os"
)

func main() {
	cli.AppHelpTemplate = fmt.Sprintf("\n\n%s\n", cli.AppHelpTemplate)

	app := cli.NewApp()
	app.Name = "rlso11n"
	app.Usage = "construct/destruct an orchestration"
	app.Version = "0.0.4"

	app.Before = core.Initialize
	app.After = core.Finalize

	app.Commands = []cli.Command{
		{
			Name:   "bg/server",
			Usage:  "(executed by " + app.Name + " only) starts rpc server",
			Action: bg.Server,
			Hidden: true,
		},
		{
			Name:   "bg/start",
			Usage:  "starts rpc server at all nodes",
			Action: bg.Start,
		},
		{
			Name: "exec@{...}",
			Usage: "executes command at given nodes\n" +
				"\t\te.g.) " + app.Name + " exec@all hostname          # at all nodes\n" +
				"\t\t      " + app.Name + " exec@manager hostname      # at manager nodes\n" +
				"\t\t      " + app.Name + " exec@worker hostname       # at worker nodes\n" +
				"\t\t      " + app.Name + " exec@0,2-7 hostname        # at 0, 2-7 nodes\n" +
				"\t\t      " + app.Name + " exec@2-10 hostname         # at 2-10 nodes",
		},
		{
			Name: "exec@{...}/docker",
			Usage: "executes docker command at given nodes\n" +
				"\t\te.g.) " + app.Name + " exec@all/docker run ...    # at all nodes",
		},
		{
			Name:   "up/dockerd",
			Usage:  "ups docker daemon in rootless mode at all nodes",
			Action: up.Dockerd,
		},
		{
			Name:   "up/swarm",
			Usage:  "ups swarm cluster using all nodes",
			Action: up.Swarm,
		},
	}
	app.CommandNotFound = cmd.Proxy

	if err := app.Run(os.Args); err != nil {
		_ = fmt.Errorf("failed to start rlso11n, %v", err)
	}
}
