package git

import (
	"fmt"
	"log"
	"os"

	"semtag/pkg"
)

func SetGitConfig() {
	fmt.Println("set git config")
	trySetGitConfigUserAndEmail()
	trySetGitCredentialsSshKey()
	trySetGitCredentialsBasicAuth()
}

func trySetGitConfigUserAndEmail() {
	gitEmail, isGitEmail := os.LookupEnv("GIT_EMAIL")
	gitUsername, isGitUser := os.LookupEnv("GIT_USERNAME")
	if !(isGitEmail && isGitUser) {
		if pkg.DEBUG != "" {
			fmt.Println("skip setting git email and username: at least one environment variable is missing ['GIT_EMAIL, GIT_USERNAME']")
		}
		return
	}
	_, err := pkg.Shellf(`
	set -euox pipefail
	which git
	git config --global user.email %v
	git config --global user.name %v
`, gitEmail, gitUsername)
	if err != nil {
		log.Fatal(err)
	}
}

func trySetGitCredentialsSshKey() {
	host, isHost := os.LookupEnv("GIT_HOSTNAME")
	sshKey, isSshKey := os.LookupEnv("GIT_SSH_KEY_PRIVATE")
	if !(isHost && isSshKey) {
		if pkg.DEBUG != "" {
			fmt.Println("skip setting git to work on SSH: at least one environment variable is missing ['GIT_HOSTNAME, GIT_SSH_KEY_PRIVATE']")
		}
		return
	}
	_, err := pkg.Shellf(`
	set -euox pipefail
	which git
	which ssh-agent

	// echo "%s" | tr -d '\r' | ssh-add - > /dev/null
	test -d ~/.ssh || (mkdir -p ~/.ssh && chmod 700 ~/.ssh)
	set +x
	echo "%s" | tr -d '\r' > ~/.ssh/id_rsa
	set -x
	eval $(ssh-agent -s)
    ssh-add ~/.ssh/id_rsa
	ssh-keyscan %v >> ~/.ssh/known_hosts
`, sshKey, host)
	if err != nil {
		log.Fatal(err)
	}

	_, err = pkg.Shellf(`
	set -euox pipefail
	which ssh-keyscan

	ssh-keyscan %v >> ~/.ssh/known_hosts
	chmod 644 ~/.ssh/known_hosts
	git config --global url.ssh://git@%s/.insteadOf https://%s/
`, host, host, host)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("changed git global url: HTTPS -> SSH")
}

func trySetGitCredentialsBasicAuth() {
	gitUsername, isGitUser := os.LookupEnv("GIT_USERNAME")
	gitPassword, isGitPw := os.LookupEnv("GIT_PASSWORD")
	if !(isGitUser && isGitPw) {
		return
	}
	_, err := pkg.Shellf(`git config credential.helper '!f() { sleep 1; echo "username=%v"; echo "password=%v"; }; f'`, gitUsername, gitPassword)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("add to git credential.helper: %s, $GIT_PASSWORD\n", gitUsername)
}
