package bg

import (
	"fmt"
	"github.com/tmyksj/rlso11n/app/logger"
	"github.com/tmyksj/rlso11n/pkg/context"
	"github.com/tmyksj/rlso11n/pkg/context/loader"
	"github.com/tmyksj/rlso11n/pkg/util/wait"
	"github.com/urfave/cli"
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"
)

func Start(_ *cli.Context) error {
	loader.LoadAsStarter()

	addr := context.Addr()
	hostList := context.HostList()

	var wg sync.WaitGroup
	wg.Add(len(hostList))

	for _, host := range hostList {
		if addr == host {
			go run(&wg, host)
		} else {
			go runViaSsh(&wg, host)
		}
	}

	wg.Wait()

	return nil
}

func run(wg *sync.WaitGroup, host string) {
	defer wg.Done()

	cmd := exec.Command("rlso11n", "bg/server")
	cmd.Env = context.Env()
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		logger.Errorf("cmd/bg", "failed to start server, %v", err)
		return
	}

	loader.Push(host)
	wait.Interrupt()

	if err := cmd.Process.Signal(syscall.SIGINT); err != nil {
		logger.Warnf("cmd/bg", "fail to send sigint to server, %v", err)

		if err := cmd.Process.Kill(); err != nil {
			logger.Errorf("cmd/bg", "fail to kill server, %v", err)
			return
		}
	}

	if err := cmd.Wait(); err != nil {
		logger.Warnf("cmd/bg", "fail to wait server, %v", err)
	}
}

func runViaSsh(wg *sync.WaitGroup, host string) {
	defer wg.Done()

	cmd := exec.Command("ssh", "-tt", host)
	cmd.Env = context.Env()
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	stdin, err := cmd.StdinPipe()

	if err != nil {
		logger.Errorf("cmd/bg", "fail to open pipe (stdin) to %v, %v", host, err)
		return
	}

	defer func() {
		if err := stdin.Close(); err != nil && err != io.EOF {
			logger.Errorf("cmd/bg", "fail to close pipe (stdin) of %v, %v", host, err)
		}
	}()

	if err := cmd.Start(); err != nil {
		logger.Errorf("cmd/bg", "failed to start ssh to %v, %v", host, err)
		return
	}

	if _, err := fmt.Fprintln(stdin, "rlso11n bg/server"); err != nil {
		logger.Errorf("cmd/bg", "failed to start worker of %v, %v", host, err)
		return
	}

	loader.Push(host)
	wait.Interrupt()

	if _, err := fmt.Fprintln(stdin, "\x03"); err != nil {
		logger.Errorf("cmd/bg", "fail to send sigint to %v, %v", host, err)
	}

	if _, err := fmt.Fprintln(stdin, "exit"); err != nil {
		logger.Warnf("cmd/bg", "fail to send exit to %v, %v", host, err)
	}

	if err := cmd.Wait(); err != nil {
		logger.Warnf("cmd/bg", "fail to wait ssh of %v, %v", host, err)
	}
}
