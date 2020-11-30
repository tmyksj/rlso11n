package up

import (
	"github.com/tmyksj/rlso11n/app/logger"
	"github.com/tmyksj/rlso11n/cmd"
	"github.com/tmyksj/rlso11n/pkg/context"
	"github.com/tmyksj/rlso11n/pkg/context/loader"
	"github.com/tmyksj/rlso11n/pkg/rpc"
	"github.com/tmyksj/rlso11n/pkg/util"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
)

func Swarm(c *cli.Context) error {
	cServer := c.String("server")

	if err := loader.LoadAsClient(cServer); err != nil {
		logger.Error(pkg, "failed to load context, %v", err)
		return cmd.ExitError(err)
	}

	var joinToken string

	eg1 := errgroup.Group{}

	for _, host := range context.HostList() {
		host := host
		eg1.Go(func() error {
			if err := util.Try(func() error {
				return networkCreateDockerGwbridge(host)
			}); err != nil {
				return err
			}

			if host.Addr == context.ManagerAddr() {
				if err := util.Try(func() error { return swarmInit(host) }); err != nil {
					return err
				}

				if err := util.Try(func() error {
					t, e := swarmJoinToken(host)
					if e == nil {
						joinToken = t
					}

					return e
				}); err != nil {
					return err
				}
			}

			return nil
		})
	}

	if err := eg1.Wait(); err != nil {
		return cmd.ExitError(err)
	}

	eg2 := errgroup.Group{}

	for _, host := range context.HostList() {
		host := host
		if host.Addr != context.ManagerAddr() {
			eg2.Go(func() error {
				if err := util.Try(func() error {
					return swarmJoin(host, joinToken)
				}); err != nil {
					return err
				}

				return nil
			})
		}
	}

	if err := eg2.Wait(); err != nil {
		return cmd.ExitError(err)
	}

	logger.Info(pkg, "the swarm cluster is ready")

	return nil
}

func networkCreateDockerGwbridge(host context.Host) error {
	err := rpc.Call(host.Addr, rpc.MtdDockerRun, &rpc.ReqDockerRun{
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
		logger.Error(pkg, "failed to create docker_gwbridge at %v, %v", host.Name, err)
		return err
	}

	logger.Info(pkg, "succeed to create docker_gwbridge at %v", host.Name)

	return nil
}

func swarmInit(host context.Host) error {
	err := rpc.Call(host.Addr, rpc.MtdDockerRun, &rpc.ReqDockerRun{
		Args: []string{
			"swarm", "init",
			"--advertise-addr", host.Addr + ":2377",
		},
	}, &rpc.ResDockerRun{})
	if err != nil {
		logger.Error(pkg, "failed to init swarm cluster at %v, %v", host.Name, err)
		return err
	}

	logger.Info(pkg, "succeed to init swarm cluster at %v", host.Name)

	return nil
}

func swarmJoin(host context.Host, token string) error {
	err := rpc.Call(host.Addr, rpc.MtdDockerRun, &rpc.ReqDockerRun{
		Args: []string{
			"swarm", "join",
			"--token", token,
			"--advertise-addr", host.Addr + ":2377",
			context.ManagerAddr() + ":2377",
		},
	}, &rpc.ResDockerRun{})
	if err != nil {
		logger.Error(pkg, "failed to join to swarm cluster at %v, %v", host.Name, err)
		return err
	}

	logger.Info(pkg, "succeed to join to swarm cluster at %v", host.Name)

	return nil
}

func swarmJoinToken(host context.Host) (string, error) {
	res := rpc.ResDockerRun{}
	err := rpc.Call(host.Addr, rpc.MtdDockerRun, &rpc.ReqDockerRun{
		Args: []string{
			"swarm", "join-token",
			"-q",
			"worker",
		},
	}, &res)
	if err != nil {
		logger.Error(pkg, "failed to get join token at %v, %v", host.Name, err)
		return "", err
	}

	logger.Info(pkg, "succeed to get join token at %v", host.Name)

	return res.Output, nil
}
