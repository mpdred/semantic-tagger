package pkg

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

var DEBUG = os.Getenv("DEBUG")

func Shell(cmd string) (string, error) {
	const ShellToUse = "bash"
	c := exec.Command(ShellToUse, "-c", cmd)
	if DEBUG != "" {
		log.Println(c.Args)
	}
	c.Stderr = os.Stderr
	out, err := c.Output()
	if DEBUG != "" {
		log.Println(string(out))
	}
	return strings.Replace(string(out), "\n", "", -1), err
}

func Shellf(format string, v ...interface{}) (string, error) {
	return Shell(fmt.Sprintf(format, v...))
}
