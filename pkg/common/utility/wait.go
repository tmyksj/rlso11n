package utility

import (
	"net"
	"os"
	"os/signal"
	"syscall"
)

func WaitInterrupt() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}

func WaitListenTcp(address string) error {
	return waitListen("tcp", address)
}

func WaitListenUnix(address string) error {
	return waitListen("unix", address)
}

func waitListen(network string, address string) error {
	return Try(func() error {
		conn, err := net.Dial(network, address)
		if err != nil {
			return err
		}

		_ = conn.Close()
		return nil
	})
}
