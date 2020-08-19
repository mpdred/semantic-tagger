package terminal

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"

	"semtag/pkg"
	"semtag/pkg/output"
)

const (
	ShellName = "bash"

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
	out, err := execute(cmd)
	if err != nil {
		return string(out), err
	}

	outputFormatted := strings.Replace(string(out), "\n", "", -1)
	return outputFormatted, nil
}

func ShellRaw(cmd string) (string, error) {
	out, err := execute(cmd)
	if err != nil {
		return string(out), err
	}

	return string(out), nil
}

func execute(cmd string) ([]byte, error) {
	c := exec.Command(ShellName, "-c", cmd)
	c.Stderr = os.Stderr
	out, err := c.Output()
	outAsString := string(out)

	output.Logger().WithFields(logrus.Fields{
		"shellCommand": cmd,
		"shellOutput":  outAsString,
	}).Debug("execute shell command")

	if err != nil {
		return []byte(BadShellResponse), pkg.NewErrorDetails(
			ErrShellCommand,
			c, outAsString)
	}
	return out, nil
}

func Shellf(format string, v ...interface{}) (string, error) {
	return Shell(fmt.Sprintf(format, v...))
}
