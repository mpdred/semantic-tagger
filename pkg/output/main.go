package output

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var debugEnvVar = os.Getenv("DEBUG")

func IsDebug() bool {
	return strings.ToLower(debugEnvVar) == "true"
}

func Debug(v ...interface{}) {
	if IsDebug() {
		log.Println(v...)
	}
}
func Debugf(format string, v ...interface{}) {
	Debug(fmt.Sprintf(format, v...))
}

func Info(v ...interface{}) {
	log.Println(v...)
}

func Infof(format string, v ...interface{}) {
	Info(fmt.Sprintf(format, v...))
}
