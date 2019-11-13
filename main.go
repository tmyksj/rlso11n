package main

import (
	"github.com/tmyksj/rootless-orchestration/cmd"
	"github.com/tmyksj/rootless-orchestration/cmd/bg"
	"github.com/tmyksj/rootless-orchestration/core"
	"github.com/tmyksj/rootless-orchestration/logger"
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "rootless-orchestration"
	app.Usage = "construct/destruct an orchestration"
	app.Version = "0.0.1"

	app.After = core.After
	app.Before = core.Before

	app.Action = bg.Action
	app.Commands = []cli.Command{
		{
			Name:   "bg/server",
			Action: bg.Server,
		},
		{
			Name:   "dockerd",
			Action: cmd.Dockerd,
		},
		{
			Name:   "swarm",
			Action: cmd.Swarm,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Fatalf("main", "%v", err)
	}
}
