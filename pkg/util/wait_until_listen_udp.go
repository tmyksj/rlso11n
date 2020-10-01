package util

func WaitUntilListenUnix(address string) {
	waitUntilListen("unix", address)
}
