package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

var hostname string

func Init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:   true,
		ForceQuote:      true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05.000000000Z07:00",
	})

	if h, err := os.Hostname(); err != nil {
		hostname = "unknown"
	} else {
		hostname = h
	}
}

func Info(pkg string, format string, args ...interface{}) {
	logrus.WithFields(logrus.Fields{
		"host": hostname,
		"pkg":  pkg,
	}).Infof(format, args...)
}

func Warn(pkg string, format string, args ...interface{}) {
	logrus.WithFields(logrus.Fields{
		"host": hostname,
		"pkg":  pkg,
	}).Warnf(format, args...)
}

func Error(pkg string, format string, args ...interface{}) {
	logrus.WithFields(logrus.Fields{
		"host": hostname,
		"pkg":  pkg,
	}).Errorf(format, args...)
}
