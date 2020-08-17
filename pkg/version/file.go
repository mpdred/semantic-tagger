package version

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"semtag/pkg"
)

var ErrNoMatchFoundVersionFormat = errors.New("no match found for version format")

type File struct {
	Path          string
	VersionFormat string
	Version       string
}

func (f *File) Read() *[]byte {
	dat, err := ioutil.ReadFile(f.Path)
	if err != nil {
		log.Fatal(err)
	}
	return &dat
}

func (f *File) Write(data *string) {
	newContents := []byte(*data)
	err := ioutil.WriteFile(f.Path, newContents, 0)
	if err != nil {
		log.Fatal(err)
	}
}

func (f *File) ReplaceSubstring() *string {
	toFind := strings.Replace(f.VersionFormat, "%s", ".*", 1)
	re := regexp.MustCompile(fmt.Sprintf("%s", toFind))
	dat := string(*f.Read())
	match := re.FindStringSubmatch(dat)
	if len(match) != 1 {
		log.Fatal(pkg.NewErrorDetails(ErrNoMatchFoundVersionFormat, "file: "+f.Path, "; version format: "+f.VersionFormat))
	}

	newVersionLine := fmt.Sprintf(f.VersionFormat, f.Version)
	newContents := strings.Replace(dat, match[0], newVersionLine, -1)
	return &newContents
}
