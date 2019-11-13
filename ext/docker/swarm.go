package docker

import (
	"github.com/tmyksj/rootless-orchestration/context"
	"github.com/tmyksj/rootless-orchestration/logger"
	"os/exec"
	"strings"
)

func SwarmInit(host string) error {
	cmd := exec.Command("docker",
		"-H", host,
		"swarm", "init",
		"--advertise-addr", host+":2377")
	cmd.Env = context.Env()

	err := cmd.Run()
	if err != nil {
		logger.Errorf("ext/docker", "fail to init swarm cluster")
		logger.Errorf("ext/docker", "%v", err)
		return err
	}

	logger.Infof("ext/docker", "succeed to init swarm cluster")

	return nil
}

func SwarmJoin(host string, token string) error {
	cmd := exec.Command("docker",
		"-H", host,
		"swarm", "join",
		"--token", token,
		"--advertise-addr", host+":2377",
		context.ManagerAddr()+":2377")
	cmd.Env = context.Env()

	err := cmd.Run()
	if err != nil {
		logger.Errorf("ext/docker", "fail to join swarm cluster")
		logger.Errorf("ext/docker", "%v", err)
		return err
	}

	logger.Infof("ext/docker", "succeed to join swarm cluster", )

	return nil
}

func SwarmJoinToken(host string) (string, error) {
	cmd := exec.Command("docker",
		"-H", host,
		"swarm", "join-token",
		"-q",
		"worker")
	cmd.Env = context.Env()

	b, err := cmd.Output()
	if err != nil {
		logger.Errorf("ext/docker", "fail to get join token")
		logger.Errorf("ext/docker", "%v", err)
		return "", err
	}

	r := strings.TrimSpace(string(b))

	logger.Infof("ext/docker", "succeed to get join token")

	return r, nil
}
