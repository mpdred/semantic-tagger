package changelog

import (
	"fmt"

	"semtag/pkg/terminal"
)

const (
	EnvVarGitCommitUrl = "GIT_COMMIT_URL"
	EnvVarGitTagUrl    = "GIT_TAG_URL"

	DefaultRegexFormat = `^%s[0-9]+\.[0-9]+\.[0-9]+%s$`
)

/*
Log represents the changelog of the repository

WARNING: THIS FEATURE IS EXPERIMENTAL. THE BEHAVIOR MAY CHANGE OR MAY BE REMOVED WITHOUT NOTICE
*/
type Log struct {
	Prefix string
	Suffix string

	// Regex for Git tags, so as to include only the relevant Git tags in the changelog
	Regex string
	// File name for the changelog
	File file

	// urlCommit is used to generating hyperlinks
	urlCommit string
	// urlTag is used to generating hyperlinks
	urlTag string
}

func NewLog() (Log, error) {
	log := Log{}

	log.File = file{}
	log.setFileName()
	if err := log.setCommitUrl(); err != nil {
		return Log{}, err
	}
	if err := log.setTagUrl(); err != nil {
		return Log{}, err
	}
	return log, nil
}

// Generate the changelog
func (l *Log) Generate() error {
	l.setRegex()
	if err := l.File.Generate(l.File.name, l.Regex, l.urlCommit, l.urlTag); err != nil {
		return err
	}
	return nil
}

func (l *Log) setRegex() {
	if l.Regex == "" {
		l.Regex = DefaultRegexFormat
	}
	l.Regex = fmt.Sprintf(l.Regex, l.Prefix, l.Suffix)
}

func (l *Log) setFileName() {
	if l.File.name != "" {
		return
	}
	l.File.name = DefaultChangelogFile
}

func (l *Log) setTagUrl() error {
	url, err := terminal.GetEnv(EnvVarGitTagUrl)
	if err != nil {
		return err
	}
	l.urlTag = url
	return nil
}

func (l *Log) setCommitUrl() error {
	url, err := terminal.GetEnv(EnvVarGitCommitUrl)
	if err != nil {
		return err
	}
	l.urlCommit = url
	return nil
}
