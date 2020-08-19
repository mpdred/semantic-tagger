package versionControl

import (
	"log"

	"semtag/pkg/output"
	"semtag/pkg/terminal"
)

func TrySetGitCredentialsBasicAuth() {
	gitUsername, err := terminal.GetEnv("GIT_USERNAME")
	if err != nil {
		return
	}
	gitPassword, err := terminal.GetEnv("GIT_PASSWORD")
	if err != nil {
		return
	}
	_, err = terminal.Shellf(`git config credential.helper '!f() { sleep 1; echo "username=%v"; echo "password=%v"; }; f'`, gitUsername, gitPassword)
	if err != nil {
		output.Logger().Fatal(err)
	}
	log.Printf("add to git credential.helper: %s, $GIT_PASSWORD\n", gitUsername)
}
