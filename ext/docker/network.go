package docker

import (
	"github.com/tmyksj/rootless-orchestration/context"
	"github.com/tmyksj/rootless-orchestration/logger"
	"os/exec"
)

func NetworkCreateDockerGwbridge(host string) error {
	cmd := exec.Command("docker",
		"-H", host,
		"network", "create",
		"--subnet", "172.20.0.0/20",
		"-o", "com.docker.network.bridge.enable_icc=true",
		"-o", "com.docker.network.bridge.enable_ip_masquerade=true",
		"-o", "com.docker.network.bridge.host_binding_ipv4=0.0.0.0",
		"-o", "com.docker.network.bridge.name=docker_gwbridge",
		"-o", "com.docker.network.driver.mtu=65520", "docker_gwbridge")
	cmd.Env = context.Env()

	err := cmd.Run()
	if err != nil {
		logger.Errorf("ext/docker", "fail to create docker_gwbridge")
		logger.Errorf("ext/docker", "%v", err)
		return err
	}

	logger.Infof("ext/docker", "succeed to create docker_gwbridge")

	return nil
}
