package util

import (
	"net"
)

func waitUntilListen(network string, address string) {
	TryUntilSucceed(func() error {
		conn, err := net.Dial(network, address)
		if err != nil {
			return err
		}

		_ = conn.Close()
		return nil
	})
}
