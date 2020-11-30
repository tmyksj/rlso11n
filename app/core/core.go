package core

import (
	"github.com/tmyksj/rlso11n/app/logger"
	"math/rand"
	"sync"
	"time"
)

func Initialize() {
	logger.Init()
	rand.Seed(time.Now().UnixNano())

	logger.Info(pkg, "initialized")
}

var finalizerList []func()
var finalizerMu sync.Mutex

func Finalize() {
	for i := len(finalizerList) - 1; i >= 0; i-- {
		finalizerList[i]()
	}

	logger.Info(pkg, "exit")
}

func RegisterFinalizer(f func()) {
	finalizerMu.Lock()
	defer finalizerMu.Unlock()

	finalizerList = append(finalizerList, f)
}
