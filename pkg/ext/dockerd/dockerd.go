package dockerd

func Start() error {
	return startRootless()
}
