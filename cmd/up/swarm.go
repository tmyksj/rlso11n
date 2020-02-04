package up

import (
	"github.com/tmyksj/rlso11n/app/logger"
	"github.com/tmyksj/rlso11n/pkg/context"
	"github.com/tmyksj/rlso11n/pkg/context/loader"
	"github.com/tmyksj/rlso11n/pkg/rpc"
	"github.com/tmyksj/rlso11n/pkg/util/attempt"
	"github.com/urfave/cli"
	"sync"
)

func Swarm(_ *cli.Context) error {
	loader.LoadAsCommander()

	var joinToken string

	var wg1 sync.WaitGroup
	wg1.Add(len(context.HostList()))

	for _, h := range context.HostList() {
		go func(host string) {
			attempt.UntilSucceed(func() error {
				return networkCreateDockerGwbridge(host)
			})
			logger.Infof("cmd/up", "created docker_gwbridge")

			if host == context.StarterAddr() {
				attempt.UntilSucceed(func() error {
					return swarmInit(host)
				})
				logger.Infof("cmd/up", "initialized a swarm cluster")

				attempt.UntilSucceed(func() error {
					t, e := swarmJoinToken(host)
					if e == nil {
						joinToken = t
					}

					return e
				})
				logger.Infof("cmd/up", "got a join token")
			}

			wg1.Done()
		}(h)
	}

	wg1.Wait()

	var wg2 sync.WaitGroup
	wg2.Add(len(context.HostList()) - 1)

	for _, h := range context.HostList() {
		if h != context.StarterAddr() {
			go func(host string) {
				attempt.UntilSucceed(func() error {
					return swarmJoin(host, joinToken)
				})
				logger.Infof("cmd/up", "joined to the swarm cluster")

				wg2.Done()
			}(h)
		}
	}

	wg2.Wait()

	logger.Infof("cmd/up", "the swarm cluster is ready")

	return nil
}

func networkCreateDockerGwbridge(host string) error {
	err := rpc.Call(host, rpc.MtdDockerRun, &rpc.ReqDockerRun{
		Args: []string{
			"network", "create",
			"--subnet", "172.20.0.0/20",
			"-o", "com.docker.network.bridge.enable_icc=true",
			"-o", "com.docker.network.bridge.enable_ip_masquerade=true",
			"-o", "com.docker.network.bridge.host_binding_ipv4=0.0.0.0",
			"-o", "com.docker.network.bridge.name=docker_gwbridge",
			"-o", "com.docker.network.driver.mtu=65520", "docker_gwbridge",
		},
	}, &rpc.ResDockerRun{})
	if err != nil {
		logger.Errorf("cmd/up", "fail to create docker_gwbridge, %v", err)
		return err
	}

	logger.Infof("cmd/up", "succeed to create docker_gwbridge")

	return nil
}

func swarmInit(host string) error {
	err := rpc.Call(host, rpc.MtdDockerRun, &rpc.ReqDockerRun{
		Args: []string{
			"swarm", "init",
			"--advertise-addr", host + ":2377",
		},
	}, &rpc.ResDockerRun{})
	if err != nil {
		logger.Errorf("cmd/up", "fail to init swarm cluster, %v", err)
		return err
	}

	logger.Infof("cmd/up", "succeed to init swarm cluster")

	return nil
}

func swarmJoin(host string, token string) error {
	err := rpc.Call(host, rpc.MtdDockerRun, &rpc.ReqDockerRun{
		Args: []string{
			"swarm", "join",
			"--token", token,
			"--advertise-addr", host + ":2377",
			context.StarterAddr() + ":2377",
		},
	}, &rpc.ResDockerRun{})
	if err != nil {
		logger.Errorf("cmd/up", "fail to join swarm cluster, %v", err)
		return err
	}

	logger.Infof("cmd/up", "succeed to join swarm cluster", )

	return nil
}

func swarmJoinToken(host string) (string, error) {
	res := rpc.ResDockerRun{}
	err := rpc.Call(host, rpc.MtdDockerRun, &rpc.ReqDockerRun{
		Args: []string{
			"swarm", "join-token",
			"-q",
			"worker",
		},
	}, &res)
	if err != nil {
		logger.Errorf("cmd/up", "fail to get join token, %v", err)
		return "", err
	}

	logger.Infof("cmd/up", "succeed to get join token")

	return res.Output, nil
}
