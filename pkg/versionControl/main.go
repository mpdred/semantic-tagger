package versionControl

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"

	"semtag/pkg/output"
	"semtag/pkg/terminal"
)

const (
	EmptyGitCommit string = ""
)

var (
	r = Open()

	ErrNotAGitRepository = errors.New("unable to find the .git folder")
	ErrGetLastGitCommits = errors.New("failed to retrieve commits")
)

func Open() *git.Repository {
	r, err := git.PlainOpen(git.GitDirName)
	if err != nil {
		log.Panic(ErrNotAGitRepository)
	}
	return r
}

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

func DescribeLong() string {
	out, err := terminal.Shell("git describe --tags --long --dirty --always")
	if err != nil {
		return GetCommitObject().Hash.String()
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
		output.Debug(err)
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

func GetCommitObject() *object.Commit {
	ref, err := r.Head()
	if err != nil {
		log.Panic(err)
	}
	out, err := r.CommitObject(ref.Hash())
	if err != nil {
		log.Panic(err)
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
	output.Debug("git fetch")
	err := r.Fetch(&git.FetchOptions{
		Progress: nil,
	})
	if err != nil {
		if strings.Contains(err.Error(), "up-to-date") {
			output.Debug(err)
			return
		}
		log.Panic(err)
	}
}
