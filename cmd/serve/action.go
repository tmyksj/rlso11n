package serve

import (
	"bufio"
	"github.com/tmyksj/rlso11n/cmd"
	"github.com/tmyksj/rlso11n/pkg/common/errors"
	"github.com/tmyksj/rlso11n/pkg/common/logger"
	"github.com/tmyksj/rlso11n/pkg/core"
	"github.com/tmyksj/rlso11n/pkg/loader"
	"github.com/tmyksj/rlso11n/pkg/service/serve"
	"github.com/urfave/cli/v2"
	"os"
	"strings"
)

func Action(c *cli.Context) error {
	return cmd.Perform(func() error {
		cDir := c.String("dir")
		cHosts := c.String("hosts")
		cHostfile := c.String("hostfile")

		var hosts []string
		if cHosts != "" {
			h, err := parseHosts(cHosts)
			if err != nil {
				return errors.By(nil, "bad request: %v", err)
			}

			hosts = h
		} else if cHostfile != "" {
			h, err := parseHostfile(cHostfile)
			if err != nil {
				return errors.By(nil, "bad request: %v", err)
			}

			hosts = h
		} else {
			return errors.By(nil, "bad request: hosts or hostfile is required")
		}

		if err := loader.LoadAsManager(cDir, hosts); err != nil {
			return err
		}

		ctx := core.GetContext()
		_, err := serve.Action(ctx).Run(&serve.ReqActionRun{})

		if err != nil {
			return err
		}

		return nil
	})
}

func parseHostfile(val string) ([]string, error) {
	file, err := os.Open(val)
	if err != nil {
		return nil, errors.By(err, "failed to open hostfile")
	}

	defer func() {
		if err := file.Close(); err != nil {
			logger.Warn(pkg, "failed to close hostfile, %v", err)
		}
	}()

	hostList := make([]string, 0)
	for s := bufio.NewScanner(file); s.Scan(); {
		t := s.Text()
		if strings.HasPrefix(t, "#") {
			continue
		}

		host := strings.Split(strings.Split(t, " ")[0], ":")[0]
		if host != "" {
			hostList = append(hostList, host)
		}
	}

	if len(hostList) == 0 {
		return nil, errors.By(nil, "host size must be equals or greater than 1")
	}

	return hostList, nil
}

func parseHosts(val string) ([]string, error) {
	hostList := make([]string, 0)
	for _, host := range strings.Split(val, ",") {
		if trimmed := strings.TrimSpace(host); trimmed != "" {
			hostList = append(hostList, trimmed)
		}
	}

	if len(hostList) == 0 {
		return nil, errors.By(nil, "host size must be equals or greater than 1")
	}

	return hostList, nil
}
