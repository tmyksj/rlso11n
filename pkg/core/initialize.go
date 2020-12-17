package core

import (
	"github.com/tmyksj/rlso11n/pkg/common/logger"
	"math/rand"
	"time"
)

func Initialize() {
	logger.Init()
	rand.Seed(time.Now().UnixNano())

	logger.Info(pkg, "initialized")
}
