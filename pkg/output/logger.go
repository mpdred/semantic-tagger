package output

import (
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	envVarDebug = "DEBUG_SEMTAG"
)

var (
	log       *logrus.Logger
	sessionId string
)

// initLogger creates a singleton of the log framework
func initLogger() {
	if log != nil {
		return
	}
	log = logrus.New()
	log.Formatter = &logrus.TextFormatter{
		ForceQuote:             true,
		DisableQuote:           false,
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		PadLevelText:           true,
		QuoteEmptyFields:       true,
	}

	logLevel := getLogLevel()
	if logLevel >= logrus.DebugLevel {
		log.SetReportCaller(true)
	}

	log.Out = os.Stdout

	Logger().WithFields(logrus.Fields{
		"logLevel":        log.Level.String(),
		"logReportCaller": log.ReportCaller,
	}).Info("logger initialized")

	log.Level = logLevel
}

func getSessionId() string {
	if sessionId != "" {
		return sessionId
	}

	newUuid, err := uuid.NewUUID()
	if err != nil {
		log.Fatal("unable to generate a session ID", err)
	}
	sessionId := newUuid.String()
	return sessionId
}

// getLogLevel looks for an environment variable to see if the user wants a specific logging level; if the environment variable is not found, it defaults to INFO
func getLogLevel() logrus.Level {
	out := os.Getenv(envVarDebug)
	logLevel := strings.ToLower(out)

	Logger().WithField("logLevelFromUser", logLevel).Info("read log level")

	switch logLevel {
	case "trace":
		return logrus.TraceLevel
	case "debug":
		return logrus.DebugLevel
	case "":
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
		return logrus.TraceLevel
	}
	return logrus.InfoLevel
}

// Logger returns the pointer of the log instance
func Logger() *logrus.Entry {
	if log == nil {
		initLogger()
	}
	return log.WithFields(logrus.Fields{"sessionId": getSessionId()})
}
