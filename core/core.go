package core

import (
	"github.com/tmyksj/rootless-orchestration/logger"
	"github.com/urfave/cli"
	"math/rand"
	"time"
)

func After(c *cli.Context) error {
	finalize()
	logger.Infof("core", "exit")

	return nil
}

func Before(c *cli.Context) error {
	initialize()
	logger.Infof("core", "completed to initialize")

	return nil
}

var finalizeFuncList []func()

func Finalize(f func()) {
	finalizeFuncList = append(finalizeFuncList, f)
}

func finalize() {
	for i := len(finalizeFuncList) - 1; i >= 0; i-- {
		finalizeFuncList[i]()
	}
}

func initialize() {
	logger.Init()
	rand.Seed(time.Now().UnixNano())
}
