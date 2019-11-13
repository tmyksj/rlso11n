package cmd

import (
	"github.com/tmyksj/rootless-orchestration/context"
	"github.com/tmyksj/rootless-orchestration/ext/docker"
	"github.com/tmyksj/rootless-orchestration/logger"
	"github.com/tmyksj/rootless-orchestration/pkg/attempt"
	"github.com/urfave/cli"
	"sync"
	"time"
)

func Swarm(c *cli.Context) error {
	context.InitFromEnv()

	var joinToken string

	var wg1 sync.WaitGroup
	wg1.Add(len(context.HostList()))

	for _, h := range context.HostList() {
		go func(host string) {
			attempt.UntilSucceed(func() error {
				return docker.NetworkCreateDockerGwbridge(host)
			}, 100*time.Millisecond)
			logger.Infof("cmd", "created docker_gwbridge")

			if host == context.Addr() {
				attempt.UntilSucceed(func() error {
					return docker.SwarmInit(host)
				}, 100*time.Millisecond)
				logger.Infof("cmd", "initialized a swarm cluster")

				attempt.UntilSucceed(func() error {
					t, e := docker.SwarmJoinToken(context.Addr())
					if e == nil {
						joinToken = t
					}

					return e
				}, 100*time.Millisecond)
				logger.Infof("cmd", "got a join token")
			}

			wg1.Done()
		}(h)
	}

	wg1.Wait()

	var wg2 sync.WaitGroup
	wg2.Add(len(context.HostList()) - 1)

	for _, h := range context.HostList() {
		if h != context.Addr() {
			go func(host string) {
				attempt.UntilSucceed(func() error {
					return docker.SwarmJoin(host, joinToken)
				}, 100*time.Millisecond)
				logger.Infof("cmd", "joined to the swarm cluster")

				wg2.Done()
			}(h)
		}
	}

	wg2.Wait()

	logger.Infof("cmd", "the swarm cluster is ready")

	return nil
}
