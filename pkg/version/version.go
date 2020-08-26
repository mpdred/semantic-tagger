package version

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"

	"semtag/pkg/output"
	"semtag/pkg/versionControl"
)

const (
	defaultVersion       = "0.1.0"
	semanticTaggingRegex = `[0-9]*\.[0-9]*\.[0-9]*`
)

var (
	ErrParseVersion      = errors.New("unable to parse version")
	ErrParseVersionMajor = errors.New("unable to parse version scope: major")
	ErrParseVersionMinor = errors.New("unable to parse version scope: minor")
	ErrParseVersionPatch = errors.New("unable to parse version scope: patch")
	ErrIncrementVersion  = errors.New("bad parameter for increment")
)

type Version struct {
	Prefix string
	Major  int
	Minor  int
	Patch  int
	Suffix string

	Hash  string
	Scope Scope
}

// UseCustomVersion will set the version number from user input
func (v *Version) UseCustomVersion(prefix string, customVersion string, suffix string) error {
	v.Suffix = suffix
	v.Prefix = prefix
	if err := v.Parse(customVersion); err != nil {
		return err
	}
	output.Logger().WithFields(logrus.Fields{
		"versionInputPrefix": prefix,
		"versionInputNumber": customVersion,
		"versionInputSuffix": suffix,
		"version":            v.String(),
	}).Info("loaded version number provided by user input")
	return nil
}

// SetVersionFromGit retrieves the latest version number based on existing git tags
func (v *Version) SetVersionFromGit() error {
	if err := versionControl.Fetch(); err != nil {
		return err
	}

	var latest string
	tag, err := versionControl.GetLatestTag(v.Prefix, semanticTaggingRegex, v.Suffix)
	if err != nil {
		latest = defaultVersion
		output.Logger().WithFields(logrus.Fields{
			"defaultVersion": defaultVersion,
			"err":            err,
		}).Warn("no previous version found, using the default version number")
	} else {
		latest = tag
		output.Logger().WithFields(logrus.Fields{
			"latestTag": latest,
		}).Debug("use the version number from git")
	}

	if err := v.Parse(latest); err != nil {
		return err
	}

	v.Hash, err = versionControl.GetHash()
	if err != nil {
		return err
	}
	output.Logger().WithFields(logrus.Fields{
		"version": v.String(),
	}).Info("got version number from git tags")
	return nil
}

// Validate if the version number adheres to semantic versioning as defined at https://semver.org/
func (v *Version) Validate(version string) error {
	expectedRegex := "^" + semanticTaggingRegex + "$"
	re := regexp.MustCompile(expectedRegex)
	allStrings := re.FindAllString(version, -1)
	if len(allStrings) == 1 {
		output.Logger().WithFields(logrus.Fields{
			"version":              version,
			"versionRegexExpected": expectedRegex,
		}).Debug("the version format is valid")
		return nil
	}
	return errors.New(fmt.Sprintf("%v: version=%q, expected regex=%q", ErrParseVersion, version, expectedRegex))
}

func (v *Version) String() string {
	s := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	s = v.appendPrefix([]string{s})[0]
	s = v.appendSuffix([]string{s})[0]
	return s
}

// Parse the version without a prefix and/or suffix
func (v *Version) Parse(version string) error {
	newV := v.removePrefixAndSuffix(version)

	if err := v.Validate(newV); err != nil {
		return err
	}

	if err := v.load(newV); err != nil {
		return err
	}

	output.Logger().WithFields(logrus.Fields{
		"version":          version,
		"versionObjParsed": fmt.Sprintf("%#v", newV),
	}).Debug("parse version number")
	return nil
}

// removePrefixAndSuffix from the version string to get only the version number
func (v *Version) removePrefixAndSuffix(version string) string {
	cleanVersion := version
	if v.Prefix != "" {
		cleanVersion = strings.Replace(cleanVersion, v.Prefix, "", 1)
	}
	if v.Suffix != "" {
		cleanVersion = strings.Replace(cleanVersion, v.Suffix, "", 1)
	}
	output.Logger().WithFields(logrus.Fields{
		"versionFull":   version,
		"versionNumber": cleanVersion,
	}).Debug("removed prefix and suffix from version string")
	return cleanVersion
}

// load a version string by updating this Version object with the new data. Get major, minor, and patch numbers from a version string
func (v *Version) load(raw string) error {
	vSplit := strings.Split(raw, ".")

	var err error
	v.Major, err = strconv.Atoi(vSplit[0])
	if err != nil {
		return errors.New(fmt.Sprintf("%v: %s", ErrParseVersionMajor, raw))
	}
	v.Minor, err = strconv.Atoi(vSplit[1])
	if err != nil {
		return errors.New(fmt.Sprintf("%v: %s", ErrParseVersionMinor, raw))
	}
	v.Patch, err = strconv.Atoi(vSplit[2])
	if err != nil {
		return errors.New(fmt.Sprintf("%v: %s", ErrParseVersionPatch, raw))
	}
	output.Logger().WithFields(logrus.Fields{
		"major": v.Major,
		"minor": v.Minor,
		"patch": v.Patch,
	}).Debug("split version string into major, minor, and patch numbers")
	return nil
}

/* AsList returns a list of stable and fixed version strings (e.g. 0.1.2 -> [0, 0.1, 0.1.2, 0.1.2-ge83655bc]

Append the version prefix and suffix to the strings (e.g. v0.1.2-api -> [v0-api, v0.1-api, v0.1.2-api, v0.1.2-ge83655bc-api]
*/
func (v *Version) AsList() []string {
	var list []string
	if v.Hash != "" {
		list = append(list, fmt.Sprintf("%d.%d.%d-g%s", v.Major, v.Minor, v.Patch, v.Hash))
	}
	list = append(list, fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch))
	list = append(list, fmt.Sprintf("%d.%d", v.Major, v.Minor))
	list = append(list, fmt.Sprint(v.Major))

	list = v.appendPrefix(list)
	list = v.appendSuffix(list)
	output.Logger().WithField("versionStrings", fmt.Sprintf("%#v", list)).Debug("generated stable and fixed version strings")
	return list
}

func (v *Version) appendSuffix(list []string) []string {
	if v.Suffix != "" {
		var list2 []string
		for _, s := range list {
			s += v.Suffix
			list2 = append(list2, s)
		}
		list = list2
	}
	return list
}

func (v *Version) appendPrefix(list []string) []string {
	if v.Prefix != "" {
		var list2 []string
		for _, s := range list {
			s = v.Prefix + s
			list2 = append(list2, s)
		}
		list = list2
	}
	return list
}

/*
SetScope calculates the version Scope that needs to be incremented:
	- if the user-provided scope if set to MAJOR or MINOR, then use that scope
	- if the user-provided scope is AUTO, try to determine the scope by parsing the commit messages
	- defaults to PATCH if no rule can be applied
*/
func (v *Version) SetScope(scopeAsString string) error {
	out, err := versionControl.GetLatestCommitLogs(-1)
	if err != nil {
		return err
	}

	s := Scope{PATCH}
	if strings.ToLower(scopeAsString) == "major" ||
		(strings.ToLower(scopeAsString) == "auto" &&
			strings.Contains(out, "BREAKING CHANGE")) {
		s.Id = MAJOR
	} else {
		if strings.ToLower(scopeAsString) == "minor" ||
			(strings.ToLower(scopeAsString) == "auto" &&
				(strings.Contains(out, "feat:") || strings.Contains(out, "feat("))) {
			s.Id = MINOR
		}
	}

	v.Scope = s
	if err := v.Increment(s); err != nil {
		return err
	}

	output.Logger().WithFields(logrus.Fields{
		"scopeFromUserInput": scopeAsString,
		"scope":              s.String(),
	}).Info("version scope has been set")

	return nil
}

/*
Increment the version number based on the Scope change
	- a breaking change increments the major number, and resets the feature and patch number to zero (e.g. 4.0.7 -> 5.0.0)
	- a feature increments the minor number, and resets the patch number to zero (e.g. 4.0.7 -> 4.1.0)
	- all other change types increment the patch number (e.g. 4.0.7 -> 4.0.8)
*/
func (v *Version) Increment(s Scope) error {
	versionBeforeIncrement := v.String()
	switch s.Id {
	case MAJOR:
		v.Major += 1
		v.Minor = 0
		v.Patch = 0
	case MINOR:
		v.Minor += 1
		v.Patch = 0
	case PATCH:
		v.Patch += 1
	default:
		return errors.New(fmt.Sprintf("%v: %s", ErrIncrementVersion, s.String()))
	}

	output.Logger().WithFields(logrus.Fields{
		"versionBeforeIncrement": versionBeforeIncrement,
		"versionAfterIncrement":  v.String(),
		"versionScope":           s.String(),
	}).Debug()

	return nil
}
