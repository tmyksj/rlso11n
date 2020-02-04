package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

var hostname string

func Errorf(pkg string, format string, args ...interface{}) {
	logrus.Errorf("%-16.16v| %-16.16v| "+format, append(append(make([]interface{}, 0), hostname, pkg), args...)...)
}

func Fatalf(pkg string, format string, args ...interface{}) {
	logrus.Fatalf("%-16.16v| %-16.16v| "+format, append(append(make([]interface{}, 0), hostname, pkg), args...)...)
}

func Infof(pkg string, format string, args ...interface{}) {
	logrus.Infof("%-16.16v| %-16.16v| "+format, append(append(make([]interface{}, 0), hostname, pkg), args...)...)
}

func Init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05.000000000Z07:00",
	})

	h, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	} else {
		hostname = h
	}
}

func Warnf(pkg string, format string, args ...interface{}) {
	logrus.Warnf("%-16.16v| %-16.16v| "+format, append(append(make([]interface{}, 0), hostname, pkg), args...)...)
}
