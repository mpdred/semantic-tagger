package versionControl

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"semtag/pkg/output"
	"semtag/pkg/terminal"
)

const (
	EnvVarGitUsername = "GIT_USERNAME"
	EnvVarGitPassword = "GIT_PASSWORD"
)

// TrySetGitCredentialsBasicAuth attempts to set the git config username and password if the environment variables are found (EnvVarGitUsername, EnvVarGitPassword)
func TrySetGitCredentialsBasicAuth() error {
	gitUsername, err := terminal.GetEnv(EnvVarGitUsername)
	if err != nil {
		return nil
	}
	gitPassword, err := terminal.GetEnv(EnvVarGitPassword)
	if err != nil {
		return nil
	}
	if _, err := terminal.Shellf(`git config credential.helper '!f() { sleep 1; echo "username=%v"; echo "password=%v"; }; f'`, gitUsername, gitPassword); err != nil {
		return fmt.Errorf("failed setting the username and password to git credential.helper: %#v", err)
	}
	output.Logger().WithFields(logrus.Fields{
		"gitUsername": gitUsername,
		"gitPassword": "***",
	}).Info("added username and password to git credential.helper")
	return nil
}
