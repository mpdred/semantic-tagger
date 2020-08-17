package terminal

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"semtag/pkg"
	"semtag/pkg/output"
)

const (
	ShellName = "sh"

	NoEnvVarValue             = ""
	BadShellResponse   string = ""
	EmptyShellResponse string = ""
)

var (
	ErrEnvVarNotFound = errors.New("environment variable not found")
	ErrShellCommand   = errors.New("failed to execute shell command")
)

func GetEnv(key string) (string, error) {
	value := os.Getenv(key)
	if len(value) == 0 {
		return NoEnvVarValue, pkg.NewErrorDetails(ErrEnvVarNotFound, key)
	}
	return value, nil
}

func Shell(cmd string) (string, error) {
	c := exec.Command(ShellName, "-c", cmd)
	c.Stderr = os.Stderr
	out, err := c.Output()
	outAsString := string(out)

	debugDetails := fmt.Sprintf("\n$ %s\n%s\n", c, outAsString)
	output.Debug(debugDetails)

	outputFormatted := strings.Replace(string(out), "\n", "", -1)
	if err != nil {
		return BadShellResponse, pkg.NewErrorDetails(
			ErrShellCommand,
			debugDetails)

	}
	return outputFormatted, nil
}

func Shellf(format string, v ...interface{}) (string, error) {
	return Shell(fmt.Sprintf(format, v...))
}
