package dockerd

import (
	"github.com/tmyksj/rootless-orchestration/context"
	"github.com/tmyksj/rootless-orchestration/core"
	"github.com/tmyksj/rootless-orchestration/logger"
	"github.com/tmyksj/rootless-orchestration/pkg/wait"
	"io/ioutil"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

var rootlessDaemon *exec.Cmd
var rootlessMu sync.Mutex

func StartRootless() error {
	rootlessMu.Lock()
	defer rootlessMu.Unlock()

	if rootlessDaemon != nil {
		logger.Infof("ext/dockerd", "dockerd is already running")
		return nil
	}

	script := context.Dir() + "/rootless-dockerd.sh"
	err := ioutil.WriteFile(script, []byte(`#!/bin/sh
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
	rm -f /run/docker /run/xtables.lock
	rootlessctl --socket=$ROOTLESSKIT_STATE_DIR/api.sock add-ports `+context.Addr()+`:2375:2375/tcp
	rootlessctl --socket=$ROOTLESSKIT_STATE_DIR/api.sock add-ports `+context.Addr()+`:2377:2377/tcp
	rootlessctl --socket=$ROOTLESSKIT_STATE_DIR/api.sock add-ports `+context.Addr()+`:7946:7946/tcp
	rootlessctl --socket=$ROOTLESSKIT_STATE_DIR/api.sock add-ports `+context.Addr()+`:7946:7946/udp
	rootlessctl --socket=$ROOTLESSKIT_STATE_DIR/api.sock add-ports `+context.Addr()+`:4789:4789/udp
	rootlessctl --socket=$ROOTLESSKIT_STATE_DIR/api.sock list-ports
	exec dockerd $@
fi
`), 0755)

	if err != nil {
		logger.Errorf("ext/dockerd", "fail to create script file")
		logger.Errorf("ext/dockerd", "%v", err)
		return err
	}

	rootlessDaemon = exec.Command(script, "--experimental",
		"--host", "unix://"+os.Getenv("XDG_RUNTIME_DIR")+"/docker.sock",
		"--host", "tcp://0.0.0.0",
		"--storage-driver", "vfs")
	rootlessDaemon.Env = context.Env()

	err = rootlessDaemon.Start()
	if err != nil {
		logger.Errorf("ext/dockerd", "fail to start dockerd")
		logger.Errorf("ext/dockerd", "%v", err)
		return err
	}

	wait.UntilListen(context.Addr()+":2375", 100*time.Millisecond)

	logger.Infof("ext/dockerd", "succeed to start dockerd")

	core.Finalize(func() {
		if err := rootlessDaemon.Process.Signal(syscall.SIGINT); err != nil {
			logger.Warnf("ext/dockerd", "fail to send sigint to dockerd")
			logger.Warnf("ext/dockerd", "%v", err)

			if err := rootlessDaemon.Process.Kill(); err != nil {
				logger.Errorf("ext/dockerd", "fail to kill dockerd")
				logger.Errorf("ext/dockerd", "%v", err)
				return
			}
		}

		if err := rootlessDaemon.Wait(); err != nil {
			logger.Warnf("ext/dockerd", "fail to wait dockerd")
			logger.Warnf("ext/dockerd", "%v", err)
		}

		logger.Infof("ext/dockerd", "succeed to stop dockerd")
	})

	return nil
}
