package versionControl

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"semtag/pkg/output"
	"semtag/pkg/terminal"
)

const (
	EmptyGitCommit string = ""
)

var (
	ErrGetLastGitCommits = errors.New("failed to retrieve commits")
)

func Commit(msg string) {
	out, err := terminal.Shellf("git commit -m %q ", msg)
	if err != nil {
		output.Logger().Fatal(err)
	}
	output.Logger().Debug(out)
}

func Add(file string) {
	out, err := terminal.Shell("git add " + file)
	if err != nil {
		output.Logger().Fatal(err)
	}
	output.Logger().Debug(out)
}

func Push(target string) {
	if target == "" {
		target = "--all"
		output.Logger().Debug("no target specified; use default target:", target)
	}
	out, err := terminal.Shell("git push origin " + target)
	if err != nil {
		output.Logger().Fatal(err)
	}
	output.Logger().Info(out)
}

func Tag(tag string, message string) {
	out, err := terminal.Shellf("git tag --annotate %q --message %q", tag, message)
	if err != nil {
		output.Logger().Fatal(err)
	}
	output.Logger().Debug(out)
}

func DescribeLong() string {
	out, err := terminal.Shell("git describe --tags --long --dirty --always")
	if err != nil {
		return GetHash()
	}
	out = strings.Replace(out, "\n", "", -1)
	return out
}

func GetLatestTag(prefix string, suffix string) (*string, error) {
	regex := `[0-9]*\.[0-9]*\.[0-9]*`
	if prefix != "" {
		regex = "^" + prefix + regex
	} else {
		regex = "^" + regex
	}
	if suffix != "" {
		regex += suffix + "$"
	} else {
		regex += "$"
	}
	cmd := fmt.Sprintf("git tag --sort=v:refname | grep -e %q | tail -1", regex)

	out, err := terminal.Shell(cmd)
	return &out, err
}

func IsAlreadyTagged(ver string) bool {
	tags, err := GetTagsForCurrentCommit()
	if err != nil {
		output.Logger().Debug(err)
		return false
	}
	re := regexp.MustCompile("^" + ver + "$")
	result := re.FindAllString(*tags, -1)
	if len(result) > 0 {
		return true
	}
	return false
}

func GetTagsForCurrentCommit() (*string, error) {
	out, err := terminal.Shell("git tag --points-at HEAD")
	return &out, err
}

func GetHash() string {
	out, err := terminal.Shell("git rev-parse HEAD")
	if err != nil {
		output.Logger().Panic(err)
	}
	return out
}

func GetLastCommits(count int) (string, error) {
	out, err := terminal.Shellf("git log %d", count)
	if err != nil {
		return EmptyGitCommit, ErrGetLastGitCommits
	}
	return out, nil
}

func Fetch() {
	out, _ := terminal.GetEnv("SEMTAG_NOFETCH")
	if out != "" {
		return
	}
	_, err := terminal.Shell("git fetch")
	if err != nil {
		output.Logger().Panic(err)
	}
}
