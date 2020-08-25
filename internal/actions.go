package internal

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	"semtag/pkg/output"
	"semtag/pkg/terminal"
	"semtag/pkg/version"
	"semtag/pkg/versionControl"
)

// TagGit creates and pushes a git tag
func TagGit(tag *versionControl.Tag, pushChanges bool) error {
	if isTagged := versionControl.IsAlreadyTagged(tag.Name); isTagged != false {
		return fmt.Errorf("the current commit has already been tagged with tag %q", tag.Name)
	}

	if pushChanges {
		if err := tag.Create(); err != nil {
			return err
		}
		if err := tag.Push(); err != nil {
			return err
		}
	}
	output.Logger().WithFields(logrus.Fields{
		"tag":       tag.Name,
		"tagPushed": pushChanges,
	}).Info("the git commit has been tagged")
	return nil
}

// TagFile updates a substring in a file based on a pattern
func TagFile(ver version.Version, filePath string, versionPattern string, pushChanges bool) error {
	const commitMsgVerBump = "chore(version): "

	f := version.File{
		Path:          filePath,
		VersionFormat: versionPattern,
		Version:       ver.String(),
	}
	newContents := f.ReplaceSubstring()

	f.Write(newContents)
	if err := versionControl.Add(filePath); err != nil {
		return err
	}
	if pushChanges {
		if err := versionControl.Commit(commitMsgVerBump + ver.String()); err != nil {
			return err
		}
		if err := versionControl.Push(""); err != nil {
			return err
		}
	}
	output.Logger().WithFields(logrus.Fields{
		"file":              f.Path,
		"fileChangesPushed": pushChanges,
	}).Info("the file has been updated")
	return nil
}

// HasRelevantChanges checks if there have been any changes in the current commit for a list of paths
func HasRelevantChanges(relevantPaths []string) (bool, error) {
	changes, err := GetChangedFiles()
	if err != nil {
		return false, err
	}

	logFields := logrus.Fields{
		"changesFound":  changes,
		"relevantPaths": relevantPaths,
	}

	for _, appPath := range relevantPaths {
		if strings.Contains(changes, appPath) {
			output.Logger().
				WithFields(logFields).
				WithField("relevantChangesFound", appPath).
				Info("found at least one relevant change in the current commit")
			return true, nil
		}
	}
	output.Logger().WithFields(logFields).
		Warn("no relevant changes found in this commit")
	return false, nil
}

// GetChangedFiles checks the HEAD commit for changes and return the changed file names
func GetChangedFiles() (string, error) {
	commit, err := versionControl.GetHash()
	if err != nil {
		return "", err
	}

	out, err := terminal.Shellf("git diff %[1]s~1..%[1]s --name-only", commit)
	if err != nil {
		return "", err
	}

	return out, nil
}
