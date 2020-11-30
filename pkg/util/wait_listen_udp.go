package util

func WaitListenUnix(address string) error {
	return waitListen("unix", address)
}
