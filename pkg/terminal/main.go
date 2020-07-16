package terminal

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"semtag/pkg/output"
)

func Shell(cmd string) (string, error) {
	const ShellToUse = "sh"
	c := exec.Command(ShellToUse, "-c", cmd)
	output.Debug(c.Args)
	c.Stderr = os.Stderr
	out, err := c.Output()
	output.Debug(string(out))
	return strings.Replace(string(out), "\n", "", -1), err
}

func Shellf(format string, v ...interface{}) (string, error) {
	return Shell(fmt.Sprintf(format, v...))
}
