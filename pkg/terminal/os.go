package terminal

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"

	"semtag/pkg/output"
)

const (
	ShellName = "bash"
)

var (
	ErrEnvVarNotFound = errors.New("environment variable not found")
	ErrShellCommand   = errors.New("failed to execute shell command")
)

// GetEnv returns the environment variable value or it throws an error if the environment variable is empty/not defined
func GetEnv(key string) (string, error) {
	value := os.Getenv(key)
	if len(value) == 0 {
		return "", fmt.Errorf("%v: key=%s", ErrEnvVarNotFound, key)
	}
	return value, nil
}

// Shell executes a command in a shell and removes all new lines from the output
func Shell(cmd string) (string, error) {
	out, err := execute(cmd)
	if err != nil {
		return string(out), err
	}

	outputFormatted := strings.Replace(string(out), "\n", "", -1)
	return outputFormatted, nil
}

// ShellRaw executes a command in a shell. Similar to Shell but this one doesn't format the output
func ShellRaw(cmd string) (string, error) {
	out, err := execute(cmd)
	if err != nil {
		return string(out), err
	}

	return string(out), nil
}

// execute a shell command
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
		return []byte(""), fmt.Errorf("%v: cmd=%q, out=%q", ErrShellCommand, cmd, outAsString)
	}
	return out, nil
}

// Shellf executes a command in a shell. It uses fmt.Sprintf to format the command before execution
func Shellf(format string, v ...interface{}) (string, error) {
	return Shell(fmt.Sprintf(format, v...))
}
