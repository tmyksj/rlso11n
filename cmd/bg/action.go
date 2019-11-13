package bg

import (
	"fmt"
	"github.com/tmyksj/rootless-orchestration/context"
	"github.com/tmyksj/rootless-orchestration/logger"
	"github.com/tmyksj/rootless-orchestration/pkg/attempt"
	"github.com/tmyksj/rootless-orchestration/pkg/wait"
	"github.com/tmyksj/rootless-orchestration/rpc"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

func Action(c *cli.Context) error {
	context.InitFromEnv()

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

	cmd := exec.Command("rootless-orchestration", "bg/server")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		logger.Errorf("cmd/bg", "failed to start server")
		logger.Errorf("cmd/bg", "%v", err)
		return
	}

	syncContext(host)
	wait.Interrupt()

	if err := cmd.Process.Signal(syscall.SIGINT); err != nil {
		logger.Warnf("cmd/bg", "fail to send sigint to server")
		logger.Warnf("cmd/bg", "%v", err)

		if err := cmd.Process.Kill(); err != nil {
			logger.Errorf("cmd/bg", "fail to kill server")
			logger.Errorf("cmd/bg", "%v", err)
			return
		}
	}

	if err := cmd.Wait(); err != nil {
		logger.Warnf("cmd/bg", "fail to wait server")
		logger.Warnf("cmd/bg", "%v", err)
	}
}

func runViaSsh(wg *sync.WaitGroup, host string) {
	defer wg.Done()

	client, err := ssh.Dial("tcp", host+":22", &ssh.ClientConfig{
		User: context.SshUsername(),
		Auth: context.SshAuthMethod(),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})

	if err != nil {
		logger.Errorf("cmd/bg", "fail to connect to %v via ssh", host)
		logger.Errorf("cmd/bg", "%v", err)
		return
	}

	defer func() {
		if err := client.Close(); err != nil {
			logger.Errorf("cmd/bg", "fail to close connection of %v", host)
			logger.Errorf("cmd/bg", "%v", err)
		}
	}()

	session, err := client.NewSession()

	if err != nil {
		logger.Errorf("cmd/bg", "fail to establish session to %v", host)
		logger.Errorf("cmd/bg", "%v", err)
		return
	}

	defer func() {
		if err := session.Close(); err != nil && err != io.EOF {
			logger.Errorf("cmd/bg", "fail to close session of %v", host)
			logger.Errorf("cmd/bg", "%v", err)
		}
	}()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm", 40, 80, modes); err != nil {
		logger.Errorf("cmd/bg", "request for pseudo terminal failed to %v", host)
		logger.Errorf("cmd/bg", "%v", err)
		return
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		logger.Errorf("cmd/bg", "fail to open pipe (stdin) to %v", host)
		logger.Errorf("cmd/bg", "%v", err)
		return
	}

	defer func() {
		if err := stdin.Close(); err != nil && err != io.EOF {
			logger.Errorf("cmd/bg", "fail to close pipe (stdin) of %v", host)
			logger.Errorf("cmd/bg", "%v", err)
		}
	}()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	if err := session.Shell(); err != nil {
		logger.Errorf("cmd/bg", "failed to start shell of %v", host)
		logger.Errorf("cmd/bg", "%v", err)
		return
	}

	if _, err := fmt.Fprintln(stdin, "rootless-orchestration bg/server"); err != nil {
		logger.Errorf("cmd/bg", "failed to start worker of %v", host)
		logger.Errorf("cmd/bg", "%v", err)
		return
	}

	syncContext(host)
	wait.Interrupt()

	if _, err := fmt.Fprintln(stdin, "\x03"); err != nil {
		logger.Errorf("cmd/bg", "fail to send sigint to session of %v", host)
		logger.Errorf("cmd/bg", "%v", err)
	}

	if _, err := fmt.Fprintln(stdin, "exit"); err != nil {
		logger.Warnf("cmd/bg", "fail to exit to session of %v", host)
		logger.Warnf("cmd/bg", "%v", err)
	}

	if err := session.Wait(); err != nil {
		logger.Warnf("cmd/bg", "fail to wait session of %v", host)
		logger.Warnf("cmd/bg", "%v", err)
	}
}

func syncContext(host string) {
	wait.UntilListen(host+":"+strconv.Itoa(context.RpcPort()), 100*time.Millisecond)

	attempt.UntilSucceed(func() error {
		return rpc.Call(host, rpc.MtdContextSync, &rpc.ReqContextSync{
			Dir:         context.Dir(),
			HostList:    strings.Join(context.HostList(), ","),
			ManagerAddr: context.ManagerAddr(),
		}, &rpc.ResContextSync{})
	}, 100*time.Millisecond)
}
