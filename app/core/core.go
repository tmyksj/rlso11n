package core

import (
	"github.com/tmyksj/rlso11n/app/logger"
	"github.com/urfave/cli"
	"math/rand"
	"time"
)

func After(_ *cli.Context) error {
	finalize()
	logger.Infof("app/core", "exit")

	return nil
}

func Before(_ *cli.Context) error {
	initialize()
	logger.Infof("app/core", "completed to initialize")

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
