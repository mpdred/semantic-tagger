package git

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"semtag/pkg/output"
	"semtag/pkg/terminal"
)

func Commit(msg string) {
	out, err := terminal.Shellf("git commit -m %q ", msg)
	if err != nil {
		log.Fatal(err)
	}
	output.Debug(out)
}

func Add(file string) {
	out, err := terminal.Shell("git add " + file)
	if err != nil {
		log.Fatal(err)
	}
	output.Debug(out)
}

func Push(target string) {
	if target == "" {
		target = "--all"
		output.Debug("no target specified; use default target:", target)
	}
	out, err := terminal.Shell("git push origin " + target)
	if err != nil {
		log.Fatal(err)
	}
	output.Debug(out)
}

func Tag(tag string, message string) {
	out, err := terminal.Shellf("git tag --annotate %q --message %q", tag, message)
	if err != nil {
		log.Fatal(err)
	}
	output.Debug(out)
}

func GetHashShort() string {
	out, err := terminal.Shell("git rev-parse --short HEAD")
	if err != nil {
		log.Fatal(err)
	}
	out = strings.Replace(out, "\n", "", -1)
	return out
}

func DescribeLong() string {
	out, err := terminal.Shell("git describe --tags --long --dirty --always")
	if err != nil {
		return GetHashShort()
	}
	out = strings.Replace(out, "\n", "", -1)
	return out
}

func GetLatestTag(prefix string, suffix string) (*string, error) {
	regex := `[0-9]*\.[0-9]*\.[0-9]*`
	if prefix != "" {
		regex = "^" + prefix + regex
	}
	if suffix != "" {
		regex += suffix
	}
	regex += "$"
	cmd := fmt.Sprintf("git tag | grep -w -e %q | sort -rn | head -1", regex)

	out, err := terminal.Shell(cmd)
	return &out, err
}

func GetBuildNumber() (*string, error) {
	out, err := terminal.Shell("git rev-list --count --first-parent HEAD")
	return &out, err
}

func GetLastCommits(count int) (*string, error) {
	out, err := terminal.Shellf("git log %d", count)
	return &out, err
}

func GetLastCommitNames(count int) (*string, error) {
	out, err := terminal.Shell("git log --pretty=format:%s " + strconv.Itoa(count))
	return &out, err
}

func Fetch() {
	cmd := "git fetch origin --tags --prune --prune-tags"
	if output.IsDebug() {
		cmd += " --verbose"
	}
	out, err := terminal.Shell(cmd)
	if err != nil {
		log.Fatal(err)
	}
	output.Debug(out)
}
