package git

import (
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
	log.Println(out)
}

func Add(file string) {
	out, err := pkg.Shell("git add " + file)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(out)
}

func Push(target string) {
	if target == "" {
		target = "--all"
		log.Println("no target specified; use default target:", target)
	}
	out, err := pkg.Shell("git push origin " + target)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(out)
}

func Tag(tag string, message string) {
	out, err := pkg.Shellf("git tag --annotate %q --message %q", tag, message)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(out)
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

func GetLatestTag(prefix string, suffix string) (*string, error) {
	cmd := "git tag"
	if prefix != "" {
		cmd += "| grep -i -e " + prefix
	}
	if suffix != "" {
		cmd += "| grep -i -e " + suffix
	}
	cmd += "| sort -rn | head -1"
	out, err := pkg.Shell(cmd)
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
	log.Println(out)
}
