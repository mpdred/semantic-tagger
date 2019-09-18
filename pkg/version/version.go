package version

import (
	"fmt"
	"log"
	"strconv"
	"strings"

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
	tag, err := git.GetLatestTag()
	if err != nil {
		fmt.Println(err)
		v.Parse("0.1.0")
	} else {
		latest := strings.Split(*tag, "-")[0]
		latest = strings.Replace(latest, "\n", "", -1)
		latest = strings.Replace(latest, "v", "", 1)
		v.Parse(latest)
	}
	return v
}

func (v *Version) Parse(version string) {
	var err error
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
	out, err := git.GetBuildNumber()
	if err != nil {
		log.Fatal(err)
	}
	v.Build, err = strconv.Atoi(*out)
	if err != nil {
		log.Fatal(err)
	}
}

func (v *Version) String() string {
	s := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.Prefix != "" {
		s = v.Prefix + s
	}
	if v.Suffix != "" {
		s += v.Suffix
	}
	return s
}

func (v *Version) AsList() []string {
	var list []string
	list = append(list, strings.Replace(git.DescribeLong(), "v", "", 1))
	list = append(list, fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch))
	list = append(list, fmt.Sprintf("%d.%d", v.Major, v.Minor))
	list = append(list, fmt.Sprint(v.Major))

	if v.Prefix != "" {
		var list2 []string
		for _, s := range list {
			s = v.Prefix + "-" + s
			list2 = append(list2, s)
		}
		list = list2
	}
	if v.Suffix != "" {
		var list2 []string
		for _, s := range list {
			s += "-" + v.Suffix
			list2 = append(list2, s)
		}
		list = list2
	}
	return list
}

func (v *Version) IncrementAuto() *ChangeType {
	out, err := git.GetLastCommits(-1)
	if err != nil {
		log.Fatal(err)
	}
	commitMsg, err := git.GetLastCommitNames(-1)
	if err != nil {
		log.Fatal(err)
	}
	if strings.Contains(*commitMsg, "ver inc") {
		log.Println("version increment skipped: version has already been incremented")
		return nil
	}
	for _, cType := range []ChangeType{BREAKING, FEATURE, PATCH} {
		if strings.Contains(*out, cType.String()) {
			v.Increment(cType)
			return &cType
		}
	}
	log.Fatal("version increment skipped: semver keywords missing from latest commit")
	return nil
}

func (v *Version) Increment(cType ChangeType) {
	switch cType {
	case BREAKING:
		v.Major += 1
		v.Minor = 0
		v.Patch = 0
	case FEATURE:
		v.Minor += 1
		v.Patch = 0
	case PATCH:
		v.Patch += 1
	default:
		log.Fatal("bad parameter for increment:", string(cType))
	}
}
