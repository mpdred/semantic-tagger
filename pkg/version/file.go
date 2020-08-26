package version

import (
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"

	"semtag/pkg/output"
)

var ErrNoMatchFoundVersionFormat = errors.New("no match found for version format")

type File struct {
	Path          string
	VersionFormat string
	Version       string
}

// Read data from file
func (f *File) Read() ([]byte, error) {
	dat, err := ioutil.ReadFile(f.Path)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to read file %q: %v", f.Path, err)
	}
	output.Logger().WithFields(logrus.Fields{
		"filePath": f.Path,
	}).Debug("read from file successfully")
	return dat, nil
}

// Write data to file
func (f *File) Write(data string) error {
	newContents := []byte(data)
	err := ioutil.WriteFile(f.Path, newContents, 0)
	if err != nil {
		return fmt.Errorf("failed to write data=%q to file=%q: %v", data, f.Path, err)
	}
	output.Logger().WithFields(logrus.Fields{
		"filePath": f.Path,
		"fileData": data,
	}).Debug("write to file successfully")
	return nil
}

// ReplaceSubstring in file
func (f *File) ReplaceSubstring() (string, error) {
	toFind := strings.Replace(f.VersionFormat, "%s", ".*", 1)
	re := regexp.MustCompile(fmt.Sprintf("%s", toFind))
	dat, err := f.Read()
	if err != nil {
		return "", err
	}
	match := re.FindStringSubmatch(string(dat))
	if len(match) != 1 {
		return "", fmt.Errorf("%v: file=%q, versionFormat=%q", ErrNoMatchFoundVersionFormat, f.Path, f.VersionFormat)
	}

	newVersionLine := fmt.Sprintf(f.VersionFormat, f.Version)
	newContents := strings.Replace(string(dat), match[0], newVersionLine, -1)

	output.Logger().WithFields(logrus.Fields{
		"filePath":           f.Path,
		"fileNewVersionLine": newVersionLine,
	}).Info("string replaced in file")
	return newContents, nil
}
