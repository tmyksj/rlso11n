package util

import (
	"net"
)

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
