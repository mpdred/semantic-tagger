package versionControl

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/sirupsen/logrus"

	"github.com/go-git/go-git/v5"

	"semtag/pkg/output"
	"semtag/pkg/terminal"
)

type GitRepository struct {
}

var repository *git.Repository

func (g *GitRepository) Open(path string) error {
	r, err := git.PlainOpen(path)
	output.Logger().WithField("path", path).Info("open git repo")
	if err != nil {
		return err
	}
	repository = r
	return nil
}

func (g *GitRepository) Fetch() error {
	opts := &git.FetchOptions{
		RemoteName: "origin",
	}

	err := repository.Fetch(opts)
	if err != nil {
		if !strings.Contains(err.Error(), "already up-to-date") {
			return fmt.Errorf("unable to sync with remote %q: %v", opts.RemoteName, err)
		}
	}

	output.Logger().WithField("details", err).Debug("successfully fetched changes from remote")
	return nil
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

func (g *GitRepository) GetTags() (*object.TagIter, error) {
	if err := g.Fetch(); err != nil {
		return nil, err
	}

	tags, err := repository.TagObjects()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve tags: %v", err)
	}
	return tags, nil
}

func (g *GitRepository) GetLatestTag(prefix, baseRegex, suffix string) (string, error) {
	regex := g.getVersionRegex(prefix, baseRegex, suffix)

	tags, err := g.GetTags()
	if err != nil {
		return "", err
	}

	err = tags.ForEach(func(t *object.Tag) error {
		fmt.Println(t)
		return nil
	})

	cmd := fmt.Sprintf("git tag --sort=v:refname | grep -E %s | tail -1", regex)
	out, err := terminal.Shell(cmd)
	if err != nil || out == "" {
		return "", fmt.Errorf("unable to get the latest tag: %v", err)
	}
	return out, nil
}

func GetLatestTagFromRepository(repository *git.Repository) (string, error) {
	tagRefs, err := repository.Tags()
	if err != nil {
		return "", err
	}

	var latestTagCommit *object.Commit
	var latestTagName string
	err = tagRefs.ForEach(func(tagRef *plumbing.Reference) error {
		revision := plumbing.Revision(tagRef.Name().String())
		tagCommitHash, err := repository.ResolveRevision(revision)
		if err != nil {
			return err
		}

		commit, err := repository.CommitObject(*tagCommitHash)
		if err != nil {
			return err
		}

		if latestTagCommit == nil {
			latestTagCommit = commit
			latestTagName = tagRef.Name().String()
		}

		if commit.Committer.When.After(latestTagCommit.Committer.When) {
			latestTagCommit = commit
			latestTagName = tagRef.Name().String()
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	return latestTagName, nil
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
	headRef, err := repository.Head()
	if err != nil {
		return "", fmt.Errorf("unable to get the hash of the HEAD commit: %v", err)
	}
	headSha := headRef.Hash().String()
	output.Logger().WithFields(logrus.Fields{
		"commit":     "HEAD",
		"commitHash": headSha,
	}).Debug("found the hash of the commit")
	return headSha, nil
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
