package version

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"semtag/pkg"
	"semtag/pkg/git"
	"semtag/pkg/output"
)

const (
	defaultVersion     string = "0.1.0"
	EmptyVersionNumber        = ""
)

var (
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

	UseGit bool
	Hash   string
}

// func (v *Version) GetVersion(args internal.CliArgs) *Version {
// 	if args.CustomVersion == EmptyVersionNumber {
// 		readFromGit(args.Prefix, args.Suffix)
// 	} else {
// 		useVersionProvidedByUser(args.CustomVersion, args.Prefix, args.Suffix)
// 	}
// 	output.Debug(v.String())
// 	return v
// }

func (v *Version) UseVersionProvidedByUser(prefix string, customVersion string, suffix string) {
	v.Suffix = suffix
	v.Prefix = prefix
	v.Parse(customVersion)
}

func (v *Version) GetLatestFromGit() *Version {
	v.UseGit = true
	git.Fetch()
	var latest string

	tag, err := git.GetLatestTag(v.Prefix, v.Suffix)
	if err != nil {
		output.Debug(err)
	} else {
		latest = *tag
	}

	if latest == "" {
		latest = defaultVersion
	}

	v.Parse(latest)
	return v
}

func (v *Version) Parse(version string) {
	var err error
	v.Validate(version)

	vSplit := strings.Split(version, ".")
	v.Major, err = strconv.Atoi(vSplit[0])
	if err != nil {
		log.Fatal(pkg.NewErrorDetails(ErrParseVersionMajor, err))
	}
	v.Minor, err = strconv.Atoi(vSplit[1])
	if err != nil {
		log.Fatal(pkg.NewErrorDetails(ErrParseVersionMinor, err))
	}
	v.Patch, err = strconv.Atoi(vSplit[2])
	if err != nil {
		log.Fatal(pkg.NewErrorDetails(ErrParseVersionPatch, err))
	}

	if v.UseGit {
		v.Hash = git.GetHashShort()
	}
}

func (v *Version) Validate(version string) bool {
	re := regexp.MustCompile("[0-9]+.[0-9]+.[0-9]+")
	allStrings := re.FindAllString(version, -1)
	if len(allStrings) == 1 {
		return true
	}
	return false
}

func (v *Version) String() string {
	s := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	s = v.appendPrefix([]string{s})[0]
	s = v.appendSuffix([]string{s})[0]
	return s
}

func (v *Version) AsList() []string {
	var list []string
	if v.UseGit {
		list = append(list, fmt.Sprintf("%d.%d.%d-%s", v.Major, v.Minor, v.Patch, v.Hash))
	}
	list = append(list, fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch))
	list = append(list, fmt.Sprintf("%d.%d", v.Major, v.Minor))
	list = append(list, fmt.Sprint(v.Major))

	list = v.appendPrefix(list)
	return v.appendSuffix(list)
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

func (v *Version) IncrementAuto(scopeAsString string) {
	out, err := git.GetLastCommits(-1)
	if err != nil {
		log.Panic(err)
	}
	s := Scope{PATCH}
	if strings.ToLower(scopeAsString) == "major" ||
		(strings.ToLower(scopeAsString) == "auto" && strings.Contains(out, "BREAKING CHANGE")) {
		s.Id = MAJOR
	} else {
		if strings.ToLower(scopeAsString) == "minor" ||
			(strings.ToLower(scopeAsString) == "auto" && (strings.Contains(out, "feat:") || strings.Contains(out, "feat("))) {
			s.Id = MINOR
		}
	}
	output.Debug("increment version number:", s.String())
	v.Increment(s)
}

func (v *Version) Increment(s Scope) {
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
		log.Panic(ErrIncrementVersion, s)
	}
}
