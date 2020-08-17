package versionControl

import (
	"log"
	"os"

	"semtag/pkg/output"
	"semtag/pkg/terminal"
)

func TryConfigureIdentity() {
	configureUserEmail()
	configureUserName()
}

func configureUserName() {
	_, err := terminal.Shell("git config --list | grep 'user.name'")
	if err == nil {
		return
	}

	gitUser, err := terminal.GetEnv("CI_GIT_USERNAME")
	if err != nil {
		output.Debug(err)
		return
	}

	_, err = terminal.Shellf(`git config --global user.name "%s"`, gitUser)
	if err != nil {
		log.Fatal(err)
	}
}

func configureUserEmail() {
	_, err := terminal.Shell("git config --list | grep 'user.email'")
	if err == nil {
		return
	}

	gitEmail, err := terminal.GetEnv("CI_GIT_EMAIL")
	if err != nil {
		output.Debug(err)
		return
	}
	_, err = terminal.Shellf(`git config --global user.email "%s"`, gitEmail)
	if err != nil {
		log.Fatal(err)
	}
}

func TrySetGitCredentialsBasicAuth() {
	gitUsername, isGitUser := os.LookupEnv("GIT_USERNAME")
	gitPassword, isGitPw := os.LookupEnv("GIT_PASSWORD")
	if !(isGitUser && isGitPw) {
		return
	}
	_, err := terminal.Shellf(`git config credential.helper '!f() { sleep 1; echo "username=%v"; echo "password=%v"; }; f'`, gitUsername, gitPassword)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("add to git credential.helper: %s, $GIT_PASSWORD\n", gitUsername)
}
