package internal

import (
	"semtag/pkg/git"
	"semtag/pkg/output"
	"semtag/pkg/version"
)

func GetVersion(args CliArgs) version.Version {
	var v version.Version
	if args.CustomVersion != "" {
		v = getVersionFromUser(args.CustomVersion, args.Prefix, args.Suffix)
	} else {
		v = getVersionFromGit(args.Prefix, args.Suffix)
	}
	output.Debug(v.String())
	return v
}

func getVersionFromUser(customVersion string, prefix string, suffix string) version.Version {
	var v version.Version
	v.Suffix = suffix
	v.Prefix = prefix
	v.Parse(customVersion)
	return v
}

func getVersionFromGit(prefix string, suffix string) version.Version {
	git.Fetch()
	var v version.Version
	v.UseGit = true
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
	newContents := f.ReplaceSubstring()

	output.Debug("tag file:", f, "\n", *newContents)
	f.Write(newContents)
	git.Add(filePath)
	if push {
		git.Commit(commitMsgVerBump + ver.String())
		git.Push("")
	}
}
