package internal

import (
	"strings"

	"semtag/pkg/output"
	"semtag/pkg/terminal"
	"semtag/pkg/version"
	"semtag/pkg/versionControl"
)

func TagGit(tag *versionControl.TagObj, push bool) {
	if versionControl.IsAlreadyTagged(tag.Name) {
		output.Logger().Fatal("The current commit has already been tagged with ", tag.Name)
	}
	output.Logger().Debug("git tag:", tag)
	if push {
		tag.Create()
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

	output.Logger().Debug("tag file:", f, "\n", *newContents)
	f.Write(newContents)
	versionControl.Add(filePath)
	if push {
		versionControl.Commit(commitMsgVerBump + ver.String())
		versionControl.Push("")
	}
}

func HasRelevantChanges(relevantPaths []string) bool {
	output.Logger().Debug("check for code changes")
	changes := GetChangedFiles()
	for _, appPath := range relevantPaths {
		if strings.Contains(changes, appPath) {
			output.Logger().Debug("found relevant changes in the current commit")
			return true
		}
	}
	output.Logger().Info("no relevant changes found in the current commit")
	return false
}

func GetChangedFiles() string {
	output.Logger().Debug("get changed files for this commit")
	commit := versionControl.GetHash()
	out, err := terminal.Shellf("git diff %[1]s~1..%[1]s --name-only", commit)
	if err != nil {
		output.Logger().Fatalln(err)
	}
	return out
}
