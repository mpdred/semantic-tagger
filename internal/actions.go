package internal

import (
	"log"
	"strings"

	"semtag/pkg/git"
	"semtag/pkg/output"
	"semtag/pkg/version"
)

func GetVersion(args CliArgs) version.Version {
	var v version.Version
	if args.CustomVersion != "" {
		v = getVersionFromUser(args.CustomVersion)
	} else {
		v = getVersionFromGit(args.Prefix, args.Suffix)
	}
	output.Debug(v.String())
	return v
}

func getVersionFromUser(customVersion string) version.Version {
	var v version.Version
	v.Parse(customVersion)
	return v
}

func getVersionFromGit(prefix string, suffix string) version.Version {
	var v version.Version
	v.Suffix = suffix
	v.Prefix = prefix
	v = *v.GetLatest()
	return v
}

func TagGit(ver version.Version, push bool) {
	tag := &git.TagObj{
		Name: ver.String(),
	}
	tag.SetMessage()
	output.Debug("git tag:", tag)
	if push {
		tag.Push()
	}
}

func TagFile(ver version.Version, filePath string, versionPattern string, push bool) {
	const commitMsgVerBump = "chore(version): "
	f := version.File{
		Path:          filePath,
		VersionFormat: versionPattern,
		Version:       ver.String(),
	}
	out, err := git.GetLastCommitNames(-1)
	if err != nil {
		log.Fatal(err)
	}
	if strings.Contains(*out, commitMsgVerBump) {
		output.Info("skip version increment: already incremented")
		return
	}
	newContents := f.ReplaceSubstring()

	output.Debug("tag file:", f, "\n", *newContents)
	f.Write(newContents)
	git.Add(filePath)
	msg, err := git.GetLastCommitNames(-1)
	if err != nil {
		log.Fatal(err)
	}
	git.Commit(commitMsgVerBump + *msg)
	if push {
		git.Push("")
	}
}
