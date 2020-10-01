package dockerd

import (
	"bufio"
	"github.com/tmyksj/rlso11n/app/core"
	"github.com/tmyksj/rlso11n/app/logger"
	"github.com/tmyksj/rlso11n/pkg/context"
	"github.com/tmyksj/rlso11n/pkg/util"
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
		logger.Info(pkg, "dockerd is already running")
		return nil
	}

	scriptPath := context.Dir() + "/" + rlsDScript
	if err := writeRlsDScript(scriptPath); err != nil {
		logger.Error(pkg, "failed to create script file, %v", err)
		return err
	}

	rlsD = exec.Command(scriptPath)
	rlsD.Env = context.Env()
	rlsD.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	if out, err := rlsD.StdoutPipe(); err != nil {
		logger.Error(pkg, "failed to pipe stdout, %v", err)
	} else {
		go pipeStdout(out)
	}

	if out, err := rlsD.StderrPipe(); err != nil {
		logger.Error(pkg, "failed to pipe stderr, %v", err)
	} else {
		go pipeStderr(out)
	}

	if err := rlsD.Start(); err != nil {
		logger.Error(pkg, "failed to start dockerd, %v", err)
		return err
	}

	core.RegisterFinalizer(rlsDFinalize)
	util.WaitUntilListenUnix(context.DockerSock())

	logger.Info(pkg, "succeed to start dockerd")

	return nil
}

func pipeStderr(r io.ReadCloser) {
	defer func() {
		if err := r.Close(); err != nil {
			logger.Warn(pkg, "failed to close pipe, %v", err)
		}
	}()

	s := bufio.NewScanner(r)
	for s.Scan() {
		logger.Info(pkg, "<- dockerd, %v", s.Text())
	}
}

func pipeStdout(r io.ReadCloser) {
	defer func() {
		if err := r.Close(); err != nil {
			logger.Warn(pkg, "failed to close pipe, %v", err)
		}
	}()

	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		logger.Info(pkg, "<- dockerd, %v", line)

		if strings.HasPrefix(line, rlsDPidPrefix) {
			if pid, err := strconv.Atoi(line[len(rlsDPidPrefix):]); err == nil {
				rlsDPid = pid
			} else {
				logger.Error(pkg, "failed to parse dockerd pid, %v", err)
			}

			break
		}
	}

	for s.Scan() {
		logger.Info(pkg, "<- dockerd, %v", s.Text())
	}
}

func rlsDFinalize() {
	requested := false

	if rlsDPid >= 0 {
		if err := syscall.Kill(rlsDPid, syscall.SIGINT); err == nil {
			requested = true
		} else {
			logger.Warn(pkg, "failed to send sigint to dockerd, %v", err)
		}
	}

	if !requested {
		if err := rlsD.Process.Signal(syscall.SIGINT); err == nil {
			requested = true
		} else {
			logger.Warn(pkg, "failed to send sigint to rootlesskit, %v", err)
		}
	}

	if !requested {
		if err := rlsD.Process.Kill(); err == nil {
			requested = true
		} else {
			logger.Error(pkg, "failed to kill rootlesskit, %v", err)
			return
		}
	}

	if err := rlsD.Wait(); err != nil {
		logger.Warn(pkg, "failed to wait dockerd, %v", err)
	}

	logger.Info(pkg, "succeed to stop dockerd")
}

func writeRlsDScript(path string) error {
	return ioutil.WriteFile(path, []byte(`#!/bin/sh
set -e -x
if [ -z $_DOCKERD_ROOTLESS_CHILD ]; then
	_DOCKERD_ROOTLESS_CHILD=1
	export _DOCKERD_ROOTLESS_CHILD

	mkdir -p `+context.Dir()+`/lower `+context.Dir()+`/upper `+context.Dir()+`/work `+context.Dir()+`/merged
	if rootlesskit mount -t overlay overlay \
			-olowerdir=`+context.Dir()+`/lower,upperdir=`+context.Dir()+`/upper,workdir=`+context.Dir()+`/work \
			`+context.Dir()+`/merged > /dev/null 2>&1; then
		STORAGE_DRIVER=overlay2
	else
		STORAGE_DRIVER=vfs
	fi
	rm -rf `+context.Dir()+`/lower `+context.Dir()+`/upper `+context.Dir()+`/work `+context.Dir()+`/merged

	exec rootlesskit \
		--copy-up=/etc \
		--copy-up=/run \
		--copy-up=/var/lib \
		--disable-host-loopback \
		--mtu=65520 \
		--net=slirp4netns \
		--port-driver=slirp4netns \
		--state-dir=`+context.Dir()+`/rootlesskit \
		$0 \
			--experimental \
			--host unix://`+context.DockerSock()+` \
			--storage-driver $STORAGE_DRIVER
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
