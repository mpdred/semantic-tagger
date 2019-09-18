package git

import (
	"fmt"
	"log"
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

func Push() {
	PushTarget("")
}

func PushTarget(target string) {
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

func DescribeLong() string {
	out, err := pkg.Shell("git describe --tags --long --dirty --always")
	if err != nil {
		log.Fatal(err)
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

func Pull() {
	out, err := pkg.Shell("git pull origin")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out)
}
