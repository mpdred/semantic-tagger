package output

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var debug = os.Getenv("DEBUG")

func Debug(v ...interface{}) {
	if strings.ToLower(debug) == "true" {
		Info(v...)
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
