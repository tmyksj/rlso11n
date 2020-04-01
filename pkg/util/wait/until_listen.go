package wait

import (
	"github.com/tmyksj/rlso11n/pkg/util/attempt"
	"net"
)

func UntilListenTcp(address string) {
	untilListen("tcp", address)
}

func UntilListenUnix(address string) {
	untilListen("unix", address)
}

func untilListen(network string, address string) {
	attempt.UntilSucceed(func() error {
		conn, err := net.Dial(network, address)
		if err != nil {
			return err
		}

		_ = conn.Close()
		return nil
	})
}
