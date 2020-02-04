package cmd

import (
	"fmt"
	"github.com/tmyksj/rlso11n/app/logger"
	"github.com/tmyksj/rlso11n/cmd/exec"
	"github.com/urfave/cli"
	"strings"
)

func CommandNotFound(c *cli.Context, command string) {
	if strings.HasPrefix(command, "exec/docker@") {
		for i := 0; i < len(c.App.Commands); i++ {
			if c.App.Commands[i].Name == "exec/docker@..." {
				c.Command.Name = command
				c.Command.Usage = c.App.Commands[i].Usage
			}
		}

		if err := exec.DockerAt(c); err != nil {
			logger.Errorf("cmd", "fail to exec docker command, %v", err)
		}

		return
	}

	fmt.Printf("No help topic for %v\n", command)
}
