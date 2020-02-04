package wait

import (
	"os"
	"os/signal"
	"syscall"
)

func Interrupt() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}
