package util

func WaitUntilListenTcp(address string) {
	waitUntilListen("tcp", address)
}
