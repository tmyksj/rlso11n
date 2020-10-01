package core

import (
	"github.com/tmyksj/rlso11n/app/logger"
	"github.com/urfave/cli"
	"math/rand"
	"time"
)

func Initialize(_ *cli.Context) error {
	logger.Init()
	rand.Seed(time.Now().UnixNano())

	logger.Info(pkg, "initialized")
	return nil
}

var finalizerList []func()

func Finalize(_ *cli.Context) error {
	for i := len(finalizerList) - 1; i >= 0; i-- {
		finalizerList[i]()
	}

	logger.Info(pkg, "exit")
	return nil
}

func RegisterFinalizer(f func()) {
	finalizerList = append(finalizerList, f)
}
