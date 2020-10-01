package cmd

import (
	"fmt"
	"github.com/tmyksj/rlso11n/app/logger"
	"github.com/tmyksj/rlso11n/cmd/exec"
	"github.com/urfave/cli"
	"strings"
)

func Proxy(c *cli.Context, command string) {
	if proxy(c, command, "exec@", "exec@{...}", exec.Exec) {
		return
	}

	fmt.Printf("No help topic for %v\n", command)
}

func proxy(c *cli.Context, command string, prefix string, fullCommand string, f func(*cli.Context) error) bool {
	if strings.HasPrefix(command, prefix) {
		for i := 0; i < len(c.App.Commands); i++ {
			if c.App.Commands[i].Name == fullCommand {
				c.Command.Name = command
				c.Command.Usage = c.App.Commands[i].Usage
			}
		}

		if err := f(c); err != nil {
			logger.Error(pkg, "failed to execute command, %v", err)
		}

		return true
	}

	return false
}
