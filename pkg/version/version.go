package version

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"semtag/pkg"
	"semtag/pkg/git"
)

type Version struct {
	Major  int
	Minor  int
	Patch  int
	Build  int
	Suffix string
	Prefix string
}

func (v *Version) GetLatest() *Version {
	const UserDefinedVersion string = "VERSION"
	userDefinedVersion, isFound := os.LookupEnv(UserDefinedVersion)

	var latest string
	if isFound {
		if len(userDefinedVersion) == 0 {
			log.Fatalf("Environment variable `%s` is defined but its value is empty!", UserDefinedVersion)
		}
		latest = userDefinedVersion
	} else {
		tag, err := git.GetLatestTag(v.Prefix, v.Suffix)
		if err != nil {
			log.Println(err)
		} else {
			latest = *tag
		}
	}

	if latest == "" {
		latest = "0.1.0"
	}

	v.Parse(latest)
	return v
}

func (v *Version) Parse(version string) {
	var err error
	re := regexp.MustCompile("[0-9]+.[0-9]+.[0-9]+")
	version = re.FindAllString(version, -1)[0]

	vSplit := strings.Split(version, ".")
	v.Major, err = strconv.Atoi(vSplit[0])
	if err != nil {
		log.Fatal(err)
	}
	v.Minor, err = strconv.Atoi(vSplit[1])
	if err != nil {
		log.Fatal(err)
	}
	v.Patch, err = strconv.Atoi(vSplit[2])
	if err != nil {
		log.Fatal(err)
	}

	buildNumber, err := git.GetBuildNumber()
	if err != nil {
		log.Fatal(err)
	}
	v.Build, err = strconv.Atoi(*buildNumber)
	if err != nil {
		log.Fatal(err)
	}
}

func (v *Version) String() string {
	s := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	s = v.appendPrefix([]string{s})[0]
	s = v.appendSuffix([]string{s})[0]
	return s
}

func (v *Version) AsList(gitDescribe string) []string {
	var list []string
	list = append(list, strings.Replace(gitDescribe, "v", "", 1))
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

func (v *Version) IncrementAuto() {
	out, err := git.GetLastCommits(-1)
	if err != nil {
		log.Fatal(err)
	}
	changeType := ChangeType(PATCH)
	if strings.Contains(*out, "BREAKING CHANGE") {
		changeType = ChangeType(MAJOR)
	} else {
		if strings.Contains(*out, "feat:") || strings.Contains(*out, "feat(") {
			changeType = ChangeType(MINOR)
		}
	}
	if pkg.DEBUG != "" {
		log.Println("increment version number:", changeType.String())
	}
	v.Increment(changeType)
}

func (v *Version) Increment(cType ChangeType) {
	switch cType {
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
		log.Fatal("bad parameter for increment:", string(cType))
	}
}
