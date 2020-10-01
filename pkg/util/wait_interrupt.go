package util

import (
	"os"
	"os/signal"
	"syscall"
)

func WaitInterrupt() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}
