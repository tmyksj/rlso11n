package serve

import (
	"bufio"
	"fmt"
	"github.com/tmyksj/rlso11n/app/logger"
	"github.com/tmyksj/rlso11n/cmd"
	"github.com/tmyksj/rlso11n/pkg/context"
	"github.com/tmyksj/rlso11n/pkg/context/loader"
	"github.com/tmyksj/rlso11n/pkg/errors"
	"github.com/tmyksj/rlso11n/pkg/util"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func Action(c *cli.Context) error {
	cDir := c.String("dir")
	cHosts := c.String("hosts")
	cHostfile := c.String("hostfile")

	var hosts []string
	if cHosts != "" {
		h, err := parseHosts(cHosts)
		if err != nil {
			logger.Error(pkg, "bad request, %v", err)
			return cmd.ExitError(err)
		}

		hosts = h
	} else if cHostfile != "" {
		h, err := parseHostfile(cHostfile)
		if err != nil {
			logger.Error(pkg, "bad request, %v", err)
			return cmd.ExitError(err)
		}

		hosts = h
	} else {
		err := errors.By(nil, "hosts or hostfile is required")

		logger.Error(pkg, "bad request, %v", err)
		return cmd.ExitError(err)
	}

	if err := loader.LoadAsManager(cDir, hosts); err != nil {
		logger.Error(pkg, "failed to load context, %v", err)
		return cmd.ExitError(err)
	}

	hostList := context.HostList()
	myAddr := context.MyAddr()

	eg := errgroup.Group{}

	for _, host := range hostList {
		host := host
		eg.Go(func() error {
			if host.Addr == myAddr {
				return run(host)
			} else {
				return runViaSsh(host)
			}
		})
	}

	if err := eg.Wait(); err != nil {
		return cmd.ExitError(err)
	}

	return nil
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

func run(host context.Host) error {
	svr := exec.Command("rlso11n", "serve-rpc")
	svr.Env = context.Env()
	svr.Stdout = os.Stderr
	svr.Stderr = os.Stderr

	if err := svr.Start(); err != nil {
		logger.Error(pkg, "failed to start worker of %v, %v", host.Name, err)
		return err
	}

	if err := loader.Push(host.Addr); err != nil {
		logger.Error(pkg, "failed to push context to %v, %v", host.Name, err)
		return err
	}

	util.WaitInterrupt()

	if err := svr.Process.Signal(syscall.SIGINT); err != nil {
		logger.Warn(pkg, "failed to send sigint to %v, %v", host.Name, err)

		if err := svr.Process.Kill(); err != nil {
			logger.Error(pkg, "failed to kill at %v, %v", host.Name, err)
			return err
		}
	}

	if err := svr.Wait(); err != nil {
		logger.Warn(pkg, "failed to wait at %v, %v", host.Name, err)
	}

	return nil
}

func runViaSsh(host context.Host) error {
	svr := exec.Command("ssh", "-tt", host.Addr)
	svr.Env = context.Env()
	svr.Stdout = os.Stderr
	svr.Stderr = os.Stderr
	svr.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	stdin, err := svr.StdinPipe()

	if err != nil {
		logger.Error(pkg, "failed to open pipe (stdin) to %v, %v", host.Name, err)
		return err
	}

	defer func() {
		if err := stdin.Close(); err != nil && err != io.EOF {
			logger.Error(pkg, "failed to close pipe (stdin) of %v, %v", host.Name, err)
		}
	}()

	if err := svr.Start(); err != nil {
		logger.Error(pkg, "failed to start ssh to %v, %v", host.Name, err)
		return err
	}

	if _, err := fmt.Fprintln(stdin, "rlso11n serve-rpc"); err != nil {
		logger.Error(pkg, "failed to start worker of %v, %v", host.Name, err)
		return err
	}

	if err := loader.Push(host.Addr); err != nil {
		logger.Error(pkg, "failed to push context to %v, %v", host.Name, err)
		return err
	}

	util.WaitInterrupt()

	if _, err := fmt.Fprintln(stdin, "\x03"); err != nil {
		logger.Error(pkg, "failed to send sigint to %v, %v", host.Name, err)
	}

	if _, err := fmt.Fprintln(stdin, "exit"); err != nil {
		logger.Warn(pkg, "failed to send exit at %v, %v", host.Name, err)
	}

	if err := svr.Wait(); err != nil {
		logger.Warn(pkg, "failed to wait ssh of %v, %v", host.Name, err)
	}

	return nil
}
