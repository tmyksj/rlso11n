package wait

import (
	"github.com/tmyksj/rlso11n/pkg/util/attempt"
	"net"
)

func UntilListen(addr string) {
	attempt.UntilSucceed(func() error {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			return err
		}

		_ = conn.Close()
		return nil
	})
}
