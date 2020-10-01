package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

var hostname string

func Init() {
	logrus.SetFormatter(&logrus.TextFormatter{
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
	logrus.Infof("%-16.16v| %-16.16v| "+format, append(append(make([]interface{}, 0), hostname, pkg), args...)...)
}

func Warn(pkg string, format string, args ...interface{}) {
	logrus.Warnf("%-16.16v| %-16.16v| "+format, append(append(make([]interface{}, 0), hostname, pkg), args...)...)
}

func Error(pkg string, format string, args ...interface{}) {
	logrus.Errorf("%-16.16v| %-16.16v| "+format, append(append(make([]interface{}, 0), hostname, pkg), args...)...)
}

func Fatal(pkg string, format string, args ...interface{}) {
	logrus.Fatalf("%-16.16v| %-16.16v| "+format, append(append(make([]interface{}, 0), hostname, pkg), args...)...)
}
