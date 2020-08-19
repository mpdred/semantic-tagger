package output

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

const envVarDebug = "DEBUG"

var log *logrus.Logger

func initLogger() {
	log = logrus.New()
	log.Formatter = new(logrus.TextFormatter)
	if isDebug() {
		log.Level = logrus.TraceLevel
		log.SetReportCaller(true)
	}
	log.Out = os.Stdout
	log.Trace("logger initialized")
}

func isDebug() bool {
	out := os.Getenv(envVarDebug)
	b := strings.ToLower(out) == "true"
	Logger().Trace("debug mode:", b)
	return b
}

func Logger() *logrus.Logger {
	if log == nil {
		initLogger()
	}
	return log
}
