package main

import (
	"fmt"
	"github.com/tmyksj/rlso11n/app/core"
	"github.com/tmyksj/rlso11n/app/logger"
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
	app.Version = "0.0.2"

	app.After = core.After
	app.Before = core.Before

	app.Commands = []cli.Command{
		{
			Name:   "bg/server",
			Usage:  "(executed by " + app.Name + " only) starts rpc server",
			Action: bg.Server,
		},
		{
			Name:   "bg/start",
			Usage:  "starts rpc server at all nodes",
			Action: bg.Start,
		},
		{
			Name: "exec/docker@...",
			Usage: "executes docker command at given nodes\n" +
				"\t\te.g.) " + app.Name + " exec/docker@all run ...        # at all nodes\n" +
				"\t\t      " + app.Name + " exec/docker@manager run ...    # at manager nodes\n" +
				"\t\t      " + app.Name + " exec/docker@worker run ...     # at worker nodes\n" +
				"\t\t      " + app.Name + " exec/docker@0,2-7 run ...      # at 0, 2-7 nodes\n" +
				"\t\t      " + app.Name + " exec/docker@2-10 run ...       # at 2-10 nodes",
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

	app.CommandNotFound = cmd.CommandNotFound

	err := app.Run(os.Args)
	if err != nil {
		logger.Fatalf("main", "%v", err)
	}
}
