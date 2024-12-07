package util

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Logger = logrus.New()

func InitLogger() {
	Logger.SetOutput(os.Stdout)
	Logger.SetFormatter(&logrus.JSONFormatter{})
	Logger.SetLevel(logrus.DebugLevel)
}
