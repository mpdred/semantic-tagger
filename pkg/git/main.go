package git

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"semtag/pkg"
)

func Commit(msg string) {
	out, err := pkg.Shellf("git commit -m %q ", msg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out)
}

func Add(file string) {
	out, err := pkg.Shell("git add " + file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out)
}

func Push(target string) {
	gitUsername, isGitUser := os.LookupEnv("GIT_USERNAME")
	gitPassword, isGitPw := os.LookupEnv("GIT_PASSWORD")
	if isGitUser && isGitPw {
		_, err := pkg.Shellf(`git config credential.helper '!f() { sleep 1; echo "username=%v"; echo "password=%v"; }; f'`, gitUsername, gitPassword)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("add to git credential.helper: %s, $GIT_PASSWORD\n", gitUsername)
	}
	if target == "" {
		target = "--all"
		log.Println("no target specified; use default target:", target)
	}
	out, err := pkg.Shell("git push origin " + target)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out)
}

func Tag(tag string, message string) {
	out, err := pkg.Shellf("git tag --annotate %q --message %q", tag, message)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out)
}

func RevParse() string {
	out, err := pkg.Shell("git rev-parse --short HEAD")
	if err != nil {
		log.Fatal(err)
	}
	out = strings.Replace(out, "\n", "", -1)
	return out
}

func DescribeLong() string {
	out, err := pkg.Shell("git describe --tags --long --dirty --always")
	if err != nil {
		return RevParse()
	}
	out = strings.Replace(out, "\n", "", -1)
	return out
}

func GetLatestTag() (*string, error) {
	out, err := pkg.Shell(`git describe --tags $(git rev-list --tags --max-count=1)`)
	return &out, err
}

func GetBuildNumber() (*string, error) {
	out, err := pkg.Shell("git rev-list --count --first-parent HEAD")
	return &out, err
}

func GetLastCommits(count int) (*string, error) {
	out, err := pkg.Shellf("git log %d", count)
	return &out, err
}

func GetLastCommitNames(count int) (*string, error) {
	out, err := pkg.Shell("git log --pretty=format:%s " + strconv.Itoa(count))
	return &out, err
}

func Fetch() {
	out, err := pkg.Shell("git fetch origin")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out)
}
