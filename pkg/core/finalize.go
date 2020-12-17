package core

import (
	"github.com/tmyksj/rlso11n/pkg/common/logger"
)

func Finalize() {
	for i := len(context.finalizerList) - 1; i >= 0; i-- {
		context.finalizerList[i]()
	}

	logger.Info(pkg, "exit")
}
