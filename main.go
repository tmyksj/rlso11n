package main

import (
	"fmt"
	"github.com/tmyksj/rlso11n/cmd/exec"
	"github.com/tmyksj/rlso11n/cmd/serve"
	"github.com/tmyksj/rlso11n/cmd/serve_rpc"
	"github.com/tmyksj/rlso11n/cmd/up"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "rlso11n"
	app.Usage = "manage engines in rootless mode"
	app.Version = "0.1.1"

	app.Commands = []*cli.Command{
		{
			Name:  "exec",
			Usage: "execute command",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "docker",
					Usage: "execute docker command",
				},
				&cli.StringFlag{
					Name:  "nodes",
					Usage: "nodes to use",
					Value: "all",
				},
				&cli.StringFlag{
					Name:  "server",
					Usage: "hostname or ip address to connect to",
					Value: "0.0.0.0",
				},
			},
			Action: exec.Action,
		},
		{
			Name:  "serve",
			Usage: "start server at all nodes",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "dir",
					Usage: "working directory to use",
				},
				&cli.StringFlag{
					Name:  "hostfile",
					Usage: "hostfile to use",
				},
				&cli.StringFlag{
					Name:  "hosts",
					Usage: "hostname or ip address of all nodes",
				},
			},
			Action: serve.Action,
		},
		{
			Name:   "serve-rpc",
			Usage:  "start rpc server",
			Hidden: true,
			Action: serve_rpc.Action,
		},
		{
			Name:  "up",
			Usage: "up engines",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "server",
					Usage: "hostname or ip address to connect to",
					Value: "0.0.0.0",
				},
			},
			Subcommands: []*cli.Command{
				{
					Name:   "dockerd",
					Usage:  "up docker daemon at all nodes",
					Action: up.Dockerd,
				},
				{
					Name:   "swarm",
					Usage:  "up swarm cluster using all nodes",
					Action: up.Swarm,
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		_ = fmt.Errorf("failed to start rlso11n, %v", err)
	}
}
