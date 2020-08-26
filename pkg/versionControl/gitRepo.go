package versionControl

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"

	"semtag/pkg/output"
	"semtag/pkg/terminal"
)

type GitRepository struct {
}

func (g *GitRepository) Commit(msg string) error {
	_, err := terminal.Shellf("git commit -m %q ", msg)
	if err != nil {
		return fmt.Errorf("unable to commit changes to git: %v", err)
	}
	output.Logger().WithField("commitMessage", msg).Info("changes committed to git")
	return nil
}

func (g *GitRepository) Add(file string) error {
	_, err := terminal.Shell("git add " + file)
	if err != nil {
		return fmt.Errorf("unable to add file %q to git: %v", file, err)
	}
	output.Logger().WithFields(logrus.Fields{
		"fileName": file,
	}).Info("file added to git")
	return nil
}

func (g *GitRepository) Push(target string) error {
	if target == "" {
		target = "--all"
		output.Logger().WithField("target", target).Debug("no target specified, using default target")
	}
	_, err := terminal.Shell("git push origin " + target)
	if err != nil {
		return fmt.Errorf("unable to push target %q to git: %v", target, err)
	}
	output.Logger().WithField("target", target).Info("git push has been successful")
	return nil
}

func (g *GitRepository) TagCommit(tag string, message string) error {
	_, err := terminal.Shellf("git tag --annotate %q --message %q", tag, message)
	if err != nil {
		return fmt.Errorf("unable to push tag %q (%s): %v", tag, message, err)
	}
	output.Logger().WithFields(logrus.Fields{
		"tagName":    tag,
		"tagMessage": message,
	}).Info("tag pushed successfully")
	return nil
}

func (g *GitRepository) DescribeLong() (string, error) {
	out, err := terminal.Shell("git describe --tags --long --dirty --always")
	if err != nil {
		output.Logger().WithField("err", err).Warn("git describe failed; falling back to using the hash")
		return g.GetHash()
	}
	out = strings.Replace(out, "\n", "", -1)
	return out, nil
}

func (g *GitRepository) GetLatestTag(prefix, baseRegex, suffix string) (string, error) {
	regex := g.getVersionRegex(prefix, baseRegex, suffix)

	cmd := fmt.Sprintf("git tag --sort=v:refname | grep -e %q | tail -1", regex)
	out, err := terminal.Shell(cmd)
	if err != nil {
		return "", fmt.Errorf("unable to get the latest tag: %v", err)
	}
	return out, nil
}

func (g *GitRepository) getVersionRegex(prefix, baseRegex, suffix string) string {
	return "^" + prefix + baseRegex + suffix + "$"
}

func (g *GitRepository) IsAlreadyTagged(ver string) bool {
	tags, err := g.GetTagsHead()

	re := regexp.MustCompile("^" + ver + "$")
	result := re.FindAllString(tags, -1)

	isTagged := len(result) > 0 && err != nil
	output.Logger().WithFields(logrus.Fields{
		"version":                ver,
		"versionIsAlreadyTagged": isTagged,
		"err":                    err,
	}).Debug("checked if the current commit is already tagged with the current version number")
	return isTagged
}

func (g *GitRepository) GetTagsHead() (string, error) {
	const commit = "HEAD"
	out, err := terminal.Shell("git tag --points-at " + commit)
	if err != nil {
		return "", fmt.Errorf("unable to get tags for %q: %v", commit, err)
	}
	output.Logger().WithFields(logrus.Fields{
		"commit":     commit,
		"commitTags": out,
	}).Debug("found tags for the commit name")
	return out, nil
}

func (g *GitRepository) GetHash() (string, error) {
	const commit = "HEAD"
	out, err := terminal.Shell("git rev-parse " + commit)
	if err != nil {
		return "", fmt.Errorf("unable to get the hash for %q: %v", commit, err)
	}
	output.Logger().WithFields(logrus.Fields{
		"commit":     commit,
		"commitHash": out,
	}).Debug("found the hash of the commit")
	return out, nil
}

func (g *GitRepository) GetLatestCommitLogs(count int) (string, error) {
	out, err := terminal.Shellf("git log %d", count)
	if err != nil {
		return "", fmt.Errorf("unable to retrieve the last %d commit logs: %v", count, err)
	}
	output.Logger().WithFields(logrus.Fields{
		"commitCount": count,
		"commitLogs":  out,
	}).Debug("retrieved the latest n commit logs")
	return out, nil
}

func (g *GitRepository) Fetch() error {
	_, err := terminal.Shell("git fetch --prune --prune-tags --tags")
	if err != nil {
		return fmt.Errorf("unable to sync with remote: %v", err)
	}
	output.Logger().Debug("successfully fetched changes from remote")
	return nil
}
