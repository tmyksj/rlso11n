package exec

import (
	"bufio"
	"github.com/tmyksj/rlso11n/pkg/common/errors"
	"github.com/tmyksj/rlso11n/pkg/common/logger"
	"github.com/tmyksj/rlso11n/pkg/common/utility"
	"github.com/tmyksj/rlso11n/pkg/core"
	"io"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

const rlsDPidPrefix string = "rlso11n/dockerd_pid:"
const rlsDScript string = "rlsd.sh"

var rlsD *exec.Cmd
var rlsDMu sync.Mutex
var rlsDPid = -1

type ObjDockerd struct {
	ctx *core.Context
}

type ReqDockerdStart struct {
}

type ResDockerdStart struct {
}

func Dockerd(ctx *core.Context) *ObjDockerd {
	return &ObjDockerd{
		ctx: ctx,
	}
}

func (obj *ObjDockerd) Start(req *ReqDockerdStart) (*ResDockerdStart, error) {
	rlsDMu.Lock()
	defer rlsDMu.Unlock()

	if obj.ctx.MyNode() == nil {
		return nil, errors.By(nil, "illegal state")
	}

	if rlsD != nil {
		logger.Info(pkg, "dockerd is already running")
		return &ResDockerdStart{}, nil
	}

	scriptPath := obj.ctx.Dir() + "/" + rlsDScript
	if err := obj.writeRlsDScript(scriptPath); err != nil {
		return nil, errors.By(err, "failed to create script file")
	}

	rlsD = exec.Command(scriptPath)
	rlsD.Env = obj.ctx.Env()
	rlsD.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	if out, err := rlsD.StdoutPipe(); err != nil {
		return nil, errors.By(err, "failed to pipe stdout")
	} else {
		go obj.pipeStdout(out)
	}

	if out, err := rlsD.StderrPipe(); err != nil {
		return nil, errors.By(err, "failed to pipe stderr")
	} else {
		go obj.pipeStderr(out)
	}

	if err := rlsD.Start(); err != nil {
		return nil, errors.By(err, "failed to start dockerd")
	}

	obj.ctx.AddFinalizer(obj.rlsDFinalize)
	if err := utility.WaitListenUnix(obj.ctx.DockerSock()); err != nil {
		logger.Warn(pkg, "failed to check listening, %v", err)
	}

	logger.Info(pkg, "succeed to start dockerd")

	return &ResDockerdStart{}, nil
}

func (obj *ObjDockerd) pipeStderr(r io.ReadCloser) {
	defer func() {
		if err := r.Close(); err != nil {
			logger.Warn(pkg, "failed to close pipe, %v", err)
		}
	}()

	s := bufio.NewScanner(r)
	for s.Scan() {
		logger.Info(pkg, "%v", s.Text())
	}
}

func (obj *ObjDockerd) pipeStdout(r io.ReadCloser) {
	defer func() {
		if err := r.Close(); err != nil {
			logger.Warn(pkg, "failed to close pipe, %v", err)
		}
	}()

	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		logger.Info(pkg, "%v", line)

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
		logger.Info(pkg, "%v", s.Text())
	}
}

func (obj *ObjDockerd) rlsDFinalize() {
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

func (obj *ObjDockerd) writeRlsDScript(path string) error {
	return ioutil.WriteFile(path, []byte(`#!/bin/sh
set -e -x
if [ -z $_DOCKERD_ROOTLESS_CHILD ]; then
	_DOCKERD_ROOTLESS_CHILD=1
	export _DOCKERD_ROOTLESS_CHILD

	mkdir -p `+obj.ctx.Dir()+`/lower `+obj.ctx.Dir()+`/upper `+obj.ctx.Dir()+`/work `+obj.ctx.Dir()+`/merged
	if rootlesskit mount -t overlay overlay \
			-olowerdir=`+obj.ctx.Dir()+`/lower,upperdir=`+obj.ctx.Dir()+`/upper,workdir=`+obj.ctx.Dir()+`/work \
			`+obj.ctx.Dir()+`/merged > /dev/null 2>&1; then
		STORAGE_DRIVER=overlay2
	else
		STORAGE_DRIVER=vfs
	fi
	rm -rf `+obj.ctx.Dir()+`/lower `+obj.ctx.Dir()+`/upper `+obj.ctx.Dir()+`/work `+obj.ctx.Dir()+`/merged

	exec rootlesskit \
		--copy-up=/etc \
		--copy-up=/run \
		--copy-up=/var/lib \
		--disable-host-loopback \
		--mtu=65520 \
		--net=slirp4netns \
		--port-driver=slirp4netns \
		--state-dir=`+obj.ctx.Dir()+`/rootlesskit \
		$0 \
			--experimental \
			--host unix://`+obj.ctx.DockerSock()+` \
			--storage-driver $STORAGE_DRIVER
else
	[ $_DOCKERD_ROOTLESS_CHILD = 1 ]
	echo "`+rlsDPidPrefix+`$$"
	rm -f /run/containerd /run/docker /run/xtables.lock /var/lib/docker
	rootlessctl --socket=$ROOTLESSKIT_STATE_DIR/api.sock add-ports `+obj.ctx.MyNode().Addr+`:2377:2377/tcp
	rootlessctl --socket=$ROOTLESSKIT_STATE_DIR/api.sock add-ports `+obj.ctx.MyNode().Addr+`:7946:7946/tcp
	rootlessctl --socket=$ROOTLESSKIT_STATE_DIR/api.sock add-ports `+obj.ctx.MyNode().Addr+`:7946:7946/udp
	rootlessctl --socket=$ROOTLESSKIT_STATE_DIR/api.sock add-ports `+obj.ctx.MyNode().Addr+`:4789:4789/udp
	rootlessctl --socket=$ROOTLESSKIT_STATE_DIR/api.sock list-ports
	exec dockerd $@
fi
`), 0755)
}
