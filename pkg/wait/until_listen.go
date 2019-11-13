package wait

import (
	"github.com/tmyksj/rootless-orchestration/pkg/attempt"
	"net"
	"time"
)

func UntilListen(addr string, d time.Duration) {
	attempt.UntilSucceed(func() error {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			return err
		}

		_ = conn.Close()
		return nil
	}, d)
}
