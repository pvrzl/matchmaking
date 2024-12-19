package log

import (
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func init() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetReportCaller(true)
	log.SetLevel(logrus.DebugLevel)

	Logger = log
}
