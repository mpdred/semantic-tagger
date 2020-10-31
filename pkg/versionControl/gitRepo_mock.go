package versionControl

import (
	"errors"
)

type GitRepositoryMock struct {
}

func (g *GitRepositoryMock) Commit(msg string) error {
	return nil
}

func (g *GitRepositoryMock) Add(file string) error {
	return nil
}

func (g *GitRepositoryMock) Push(target string) error {
	return nil
}

func (g *GitRepositoryMock) TagCommit(tag string, message string) error {
	return nil
}

func (g *GitRepositoryMock) DescribeLong() (string, error) {
	return "", nil
}

func (g *GitRepositoryMock) GetLatestTag(prefix, baseRegex, suffix string) (string, error) {
	return "", errors.New("")
}

func (g *GitRepositoryMock) getVersionRegex(prefix, baseRegex, suffix string) string {
	actual := GitRepository{}
	return actual.getVersionRegex(prefix, baseRegex, suffix)
}

func (g *GitRepositoryMock) IsAlreadyTagged(ver string) bool {
	return false
}

func (g *GitRepositoryMock) GetTagsHead() (string, error) {
	return "", nil
}

func (g *GitRepositoryMock) GetHash() (string, error) {
	return "", nil
}

func (g *GitRepositoryMock) GetLatestCommitLogs(count int) (string, error) {
	return "", nil
}

func (g *GitRepositoryMock) Fetch() error {
	return nil
}
