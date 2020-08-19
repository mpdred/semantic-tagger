package output

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	envVarDebug = "DEBUG"
)

var log *logrus.Logger

func initLogger() {
	log = logrus.New()
	log.Formatter = new(logrus.TextFormatter)

	logLevel := getLogLevel()
	log.Level = logLevel
	if logLevel <= logrus.DebugLevel {
		log.SetReportCaller(true)
	}
	log.Out = os.Stdout
	log.Trace("logger initialized")
}

func getLogLevel() logrus.Level {
	out := os.Getenv(envVarDebug)

	switch strings.ToLower(out) {
	case "trace":
		return logrus.TraceLevel
	case "true":
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	default:
		return logrus.InfoLevel
	}
	Logger().Panic()
	return -1
}

func Logger() *logrus.Logger {
	if log == nil {
		initLogger()
	}
	return log
}
