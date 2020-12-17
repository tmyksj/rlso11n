package up

import (
	"github.com/tmyksj/rlso11n/pkg/common/logger"
	"github.com/tmyksj/rlso11n/pkg/common/object"
	"github.com/tmyksj/rlso11n/pkg/common/utility"
	"github.com/tmyksj/rlso11n/pkg/component/rpc"
	"github.com/tmyksj/rlso11n/pkg/core"
	"golang.org/x/sync/errgroup"
)

type ObjSwarm struct {
	ctx   *core.Context
	token string
}

type ReqSwarmRun struct {
}

type ResSwarmRun struct {
}

func Swarm(ctx *core.Context) *ObjSwarm {
	return &ObjSwarm{
		ctx: ctx,
	}
}

func (obj *ObjSwarm) Run(req *ReqSwarmRun) (*ResSwarmRun, error) {
	eg1 := errgroup.Group{}

	for _, node := range obj.ctx.NodeList() {
		node := node
		eg1.Go(func() error {
			err := utility.Try(func() error {
				return obj.networkCreateDockerGwbridge(node)
			})

			if err != nil {
				return err
			}

			if node.Addr == obj.ctx.ManagerNode().Addr {
				err := utility.Try(func() error {
					return obj.swarmInit(node)
				})

				if err != nil {
					return err
				}

				err = utility.Try(func() error {
					return obj.swarmJoinToken(node)
				})

				if err != nil {
					return err
				}
			}

			return nil
		})
	}

	if err := eg1.Wait(); err != nil {
		return nil, err
	}

	eg2 := errgroup.Group{}

	for _, node := range obj.ctx.NodeList() {
		node := node
		if node.Addr != obj.ctx.ManagerNode().Addr {
			eg2.Go(func() error {
				err := utility.Try(func() error {
					return obj.swarmJoin(node)
				})

				if err != nil {
					return err
				}

				return nil
			})
		}
	}

	if err := eg2.Wait(); err != nil {
		return nil, err
	}

	logger.Info(pkg, "the swarm cluster is ready")

	return &ResSwarmRun{}, nil
}

func (obj *ObjSwarm) networkCreateDockerGwbridge(node *object.Node) error {
	_, err := rpc.Docker(obj.ctx, node).Run(&rpc.ReqDockerRun{
		Args: []string{
			"network", "create",
			"--subnet", "172.20.0.0/20",
			"-o", "com.docker.network.bridge.enable_icc=true",
			"-o", "com.docker.network.bridge.enable_ip_masquerade=true",
			"-o", "com.docker.network.bridge.host_binding_ipv4=0.0.0.0",
			"-o", "com.docker.network.bridge.name=docker_gwbridge",
			"-o", "com.docker.network.driver.mtu=65520", "docker_gwbridge",
		},
	})

	if err != nil {
		logger.Error(pkg, "failed to create docker_gwbridge at %v, %v", node.Name, err)
		return err
	}

	logger.Info(pkg, "succeed to create docker_gwbridge at %v", node.Name)

	return nil
}

func (obj *ObjSwarm) swarmInit(node *object.Node) error {
	_, err := rpc.Docker(obj.ctx, node).Run(&rpc.ReqDockerRun{
		Args: []string{
			"swarm", "init",
			"--advertise-addr", node.Addr + ":2377",
		},
	})

	if err != nil {
		logger.Error(pkg, "failed to init swarm cluster at %v, %v", node.Name, err)
		return err
	}

	logger.Info(pkg, "succeed to init swarm cluster at %v", node.Name)

	return nil
}

func (obj *ObjSwarm) swarmJoin(node *object.Node) error {
	_, err := rpc.Docker(obj.ctx, node).Run(&rpc.ReqDockerRun{
		Args: []string{
			"swarm", "join",
			"--token", obj.token,
			"--advertise-addr", node.Addr + ":2377",
			obj.ctx.ManagerNode().Addr + ":2377",
		},
	})

	if err != nil {
		logger.Error(pkg, "failed to join to swarm cluster at %v, %v", node.Name, err)
		return err
	}

	logger.Info(pkg, "succeed to join to swarm cluster at %v", node.Name)

	return nil
}

func (obj *ObjSwarm) swarmJoinToken(node *object.Node) error {
	res, err := rpc.Docker(obj.ctx, node).Run(&rpc.ReqDockerRun{
		Args: []string{
			"swarm", "join-token",
			"-q",
			"worker",
		},
	})

	if err != nil {
		logger.Error(pkg, "failed to get join token at %v, %v", node.Name, err)
		return err
	}

	obj.token = res.Stdout

	logger.Info(pkg, "succeed to get join token at %v", node.Name)

	return nil
}
