package util

func WaitListenTcp(address string) error {
	return waitListen("tcp", address)
}
