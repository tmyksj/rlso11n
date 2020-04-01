package dockerd

import (
	"bufio"
	"github.com/tmyksj/rlso11n/app/core"
	"github.com/tmyksj/rlso11n/app/logger"
	"github.com/tmyksj/rlso11n/pkg/context"
	"github.com/tmyksj/rlso11n/pkg/util/wait"
	"io"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

const rlsDPidPrefix = "rlso11n/dockerd_pid:"
const rlsDScript = "rlsd.sh"

var rlsD *exec.Cmd
var rlsDPid = -1
var rlsMu sync.Mutex

func startRootless() error {
	rlsMu.Lock()
	defer rlsMu.Unlock()

	if rlsD != nil {
		logger.Infof("pkg/ext/dockerd", "dockerd is already running")
		return nil
	}

	scriptPath := context.Dir() + "/" + rlsDScript
	if err := writeRlsDScript(scriptPath); err != nil {
		logger.Errorf("pkg/ext/dockerd", "fail to create script file, %v", err)
		return err
	}

	rlsD = exec.Command(scriptPath, "--experimental",
		"--host", "unix://"+context.DockerSock(),
		"--storage-driver", "vfs")
	rlsD.Env = context.Env()
	rlsD.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	if out, err := rlsD.StdoutPipe(); err != nil {
		logger.Errorf("pkg/ext/dockerd", "fail to pipe stdout, %v", err)
	} else {
		go pipeStdout(out)
	}

	if out, err := rlsD.StderrPipe(); err != nil {
		logger.Errorf("pkg/ext/dockerd", "fail to pipe stderr, %v", err)
	} else {
		go pipeStderr(out)
	}

	if err := rlsD.Start(); err != nil {
		logger.Errorf("pkg/ext/dockerd", "fail to start dockerd, %v", err)
		return err
	}

	wait.UntilListenUnix(context.DockerSock())

	logger.Infof("pkg/ext/dockerd", "succeed to start dockerd")

	core.Finalize(rlsDFinalize)

	return nil
}

func pipeStderr(r io.ReadCloser) {
	defer func() {
		if err := r.Close(); err != nil {
			logger.Warnf("pkg/ext/dockerd", "fail to close pipe, %v", err)
		}
	}()

	s := bufio.NewScanner(r)
	for s.Scan() {
		logger.Infof("pkg/ext/dockerd", s.Text())
	}
}

func pipeStdout(r io.ReadCloser) {
	defer func() {
		if err := r.Close(); err != nil {
			logger.Warnf("pkg/ext/dockerd", "fail to close pipe, %v", err)
		}
	}()

	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		logger.Infof("pkg/ext/dockerd", line)

		if strings.HasPrefix(line, rlsDPidPrefix) {
			if pid, err := strconv.Atoi(line[len(rlsDPidPrefix):]); err == nil {
				rlsDPid = pid
			} else {
				logger.Errorf("pkg/ext/dockerd", "fail to parse dockerd pid, %v", err)
			}

			break
		}
	}

	for s.Scan() {
		logger.Infof("pkg/ext/dockerd", s.Text())
	}
}

func rlsDFinalize() {
	requested := false

	if rlsDPid >= 0 {
		if err := syscall.Kill(rlsDPid, syscall.SIGINT); err == nil {
			requested = true
		} else {
			logger.Warnf("pkg/ext/dockerd", "fail to send sigint to dockerd, %v", err)
		}
	}

	if !requested {
		if err := rlsD.Process.Signal(syscall.SIGINT); err == nil {
			requested = true
		} else {
			logger.Warnf("pkg/ext/dockerd", "fail to send sigint to rootlesskit, %v", err)
		}
	}

	if !requested {
		if err := rlsD.Process.Kill(); err == nil {
			requested = true
		} else {
			logger.Errorf("pkg/ext/dockerd", "fail to kill rootlesskit, %v", err)
			return
		}
	}

	if err := rlsD.Wait(); err != nil {
		logger.Warnf("pkg/ext/dockerd", "fail to wait dockerd, %v", err)
	}

	logger.Infof("pkg/ext/dockerd", "succeed to stop dockerd")
}

func writeRlsDScript(path string) error {
	return ioutil.WriteFile(path, []byte(`#!/bin/sh
set -e -x
if [ -z $_DOCKERD_ROOTLESS_CHILD ]; then
	_DOCKERD_ROOTLESS_CHILD=1
	export _DOCKERD_ROOTLESS_CHILD
	exec rootlesskit \
		--state-dir=`+context.Dir()+`/rootlesskit \
		--net=slirp4netns --mtu=65520 \
		--disable-host-loopback --port-driver=slirp4netns \
		--copy-up=/etc --copy-up=/run \
		--copy-up=/var/lib \
		$0 $@
else
	[ $_DOCKERD_ROOTLESS_CHILD = 1 ]
	echo "`+rlsDPidPrefix+`$$"
	rm -f /run/docker /run/xtables.lock
	rootlessctl --socket=$ROOTLESSKIT_STATE_DIR/api.sock add-ports `+context.Addr()+`:2377:2377/tcp
	rootlessctl --socket=$ROOTLESSKIT_STATE_DIR/api.sock add-ports `+context.Addr()+`:7946:7946/tcp
	rootlessctl --socket=$ROOTLESSKIT_STATE_DIR/api.sock add-ports `+context.Addr()+`:7946:7946/udp
	rootlessctl --socket=$ROOTLESSKIT_STATE_DIR/api.sock add-ports `+context.Addr()+`:4789:4789/udp
	rootlessctl --socket=$ROOTLESSKIT_STATE_DIR/api.sock list-ports
	exec dockerd $@
fi
`), 0755)
}
