package versionControl

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"

	"semtag/pkg/output"
	"semtag/pkg/terminal"
)

// Commit the current changes
func Commit(msg string) error {
	_, err := terminal.Shellf("git commit -m %q ", msg)
	if err != nil {
		return fmt.Errorf("unable to commit changes to git: %v", err)
	}
	output.Logger().WithField("commitMessage", msg).Info("changes committed to git")
	return nil
}

// Add a file to the git index
func Add(file string) error {
	_, err := terminal.Shell("git add " + file)
	if err != nil {
		return fmt.Errorf("unable to add file %q to git: %v", file, err)
	}
	output.Logger().WithFields(logrus.Fields{
		"fileName": file,
	}).Info("file added to git")
	return nil
}

// Push the local commits to remote
func Push(target string) error {
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

// Tag the current commit with an annotated git tag
func TagCommit(tag string, message string) error {
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

// DescribeLong gives an object a human readable name based on an available ref
func DescribeLong() (string, error) {
	out, err := terminal.Shell("git describe --tags --long --dirty --always")
	if err != nil {
		output.Logger().WithField("err", err).Warn("git describe failed; falling back to using the hash")
		return GetHash()
	}
	out = strings.Replace(out, "\n", "", -1)
	return out, nil
}

// GetLatestTag returns the latest tag that adheres to the semantic versioning regex and that has the provided prefix and suffix
func GetLatestTag(prefix, baseRegex, suffix string) (string, error) {
	regex := getVersionRegex(prefix, baseRegex, suffix)

	cmd := fmt.Sprintf("git tag --sort=v:refname | grep -e %q | tail -1", regex)
	out, err := terminal.Shell(cmd)
	if err != nil {
		return "", fmt.Errorf("unable to get the latest tag: %v", err)
	}
	return out, nil
}

func getVersionRegex(prefix, baseRegex, suffix string) string {
	return "^" + prefix + baseRegex + suffix + "$"
}

/*
IsAlreadyTagged checks if the current commit has been tagged with the current version number
Rules:
	- if there is no tag for the commit then return false
	- if the current commit has a different tag than the current version then return false
*/
func IsAlreadyTagged(ver string) bool {
	tags, err := GetTagsHead()

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

// GetTagsHead retrieves all the tags for the HEAD commit
func GetTagsHead() (string, error) {
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

// GetHash returns the git has for the HEAD commit
func GetHash() (string, error) {
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

// GetLatestCommitLogs returns the latest n commit logs
func GetLatestCommitLogs(count int) (string, error) {
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

// Fetch downloads the objects and refs from the remote
func Fetch() error {
	_, err := terminal.Shell("git fetch --prune --prune-tags --tags")
	if err != nil {
		return fmt.Errorf("unable to sync with remote: %v", err)
	}
	output.Logger().Debug("successfully fetched changes from remote")
	return nil
}
