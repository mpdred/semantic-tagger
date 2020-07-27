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
func DebugF(format string, v ...interface{}) {
	Debug(fmt.Sprintf(format, v...))
}

func Info(v ...interface{}) {
	log.Println(v...)
}

func InfoF(format string, v ...interface{}) {
	Info(fmt.Sprintf(format, v...))
}
