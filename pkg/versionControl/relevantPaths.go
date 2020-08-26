package versionControl

import (
	"strings"

	"github.com/sirupsen/logrus"

	"semtag/pkg/output"
	"semtag/pkg/terminal"
)

const (
	DefaultRelevantPath = "."
)

var g VersionControl = &GitRepository{}

type RelevantPaths []string

func (i *RelevantPaths) String() string {
	return DefaultRelevantPath
}

func (i *RelevantPaths) Set(value string) error {
	*i = append(*i, value)
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
		"RelevantPaths": relevantPaths,
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
	commit, err := g.GetHash()
	if err != nil {
		return "", err
	}

	out, err := terminal.Shellf("git diff %[1]s~1..%[1]s --name-only", commit)
	if err != nil {
		return "", err
	}

	return out, nil
}
