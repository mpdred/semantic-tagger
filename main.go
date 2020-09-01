package main

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"

	"semtag/internal"
	"semtag/pkg/changelog"
	"semtag/pkg/output"
	"semtag/pkg/terminal"
	"semtag/pkg/version"
	"semtag/pkg/versionControl"
)

var (
	ErrNotPushMode = errors.New("push to Git skipped: use the `-push` flag to push changes")
)

var GitRepo versionControl.VersionControl = &versionControl.GitRepository{}

func main() {
	args := internal.CliArgs{}
	args.ParseFlags()

	v := setVersion(args)

	// print the version to stdout; execute as the last command so that it can be grepped by simple shell scripts
	defer fmt.Print(v.String())

	if args.Push {
		if err := versionControl.TrySetGitCredentialsBasicAuth(); err != nil {
			output.Logger().Debug(err)
		}
	}

	if args.ExecuteCommand != "" {
		for _, val := range v.AsList() {
			_, err := terminal.Shellf(args.ExecuteCommand, val)
			if err != nil {
				output.Logger().Fatal(err)
			}
		}
	}

	if args.ShouldTagGit {
		hasRelevantChanges, err := versionControl.HasRelevantChanges(args.RelevantPaths)
		if err != nil {
			output.Logger().Fatal(err)
		}

		if hasRelevantChanges {
			tag := &versionControl.Tag{
				Name: v.String(),
			}
			if err := TagGit(tag, args.Push); err != nil {
				output.Logger().Fatal(err)
			}
			if !args.Push {
				output.Logger().Warn(ErrNotPushMode)
			}
		}
	}

	shouldTagInFile := len(args.FileName) > 0 && len(args.FileVersionPattern) > 0
	if shouldTagInFile {
		if err := TagFile(v, args.FileName, args.FileVersionPattern, args.Push); err != nil {
			output.Logger().Fatal(err)
		}
		if !args.Push {
			output.Logger().Warn(ErrNotPushMode)
		}
	}

	if args.Changelog {
		chLog, err := changelog.NewLog()
		if err != nil {
			output.Logger().Fatal(err)
		}
		chLog.Prefix = args.Prefix
		chLog.Suffix = args.Suffix
		chLog.Regex = args.ChangelogRegex
		if err := chLog.Generate(); err != nil {
			output.Logger().Fatal(err)
		}
	}

}

func setVersion(args internal.CliArgs) version.Version {
	v := version.Version{
		Prefix: args.Prefix,
		Suffix: args.Suffix,
	}

	if args.CustomVersion == "" {
		if err := v.SetVersionFromGit(); err != nil {
			output.Logger().Fatal(err)
		}
	} else {
		if err := v.UseCustomVersion(args.Prefix, args.CustomVersion, args.Suffix); err != nil {
			output.Logger().Fatal(err)
		}
	}

	if err := v.SetIncrementScope(args.VersionScopeAsString); err != nil {
		output.Logger().Fatal(err)
	}
	return v
}

// TagGit creates and pushes a git tag
func TagGit(tag *versionControl.Tag, pushChanges bool) error {
	if pushChanges {
		if err := tag.Create(); err != nil {
			return err
		}
		if err := tag.Push(); err != nil {
			return err
		}
	}
	output.Logger().WithFields(logrus.Fields{
		"tag":       tag.Name,
		"tagPushed": pushChanges,
	}).Info("the git commit has been tagged")
	return nil
}

// TagFile updates a substring in a file based on a pattern
func TagFile(ver version.Version, filePath string, versionPattern string, pushChanges bool) error {
	const commitMsgVerBump = "chore(version): "

	f := version.File{
		Path:          filePath,
		VersionFormat: versionPattern,
		Version:       ver.String(),
	}
	newContents, err := f.ReplaceSubstring()
	if err != nil {
		return err
	}

	if err := f.Write(newContents); err != nil {
		return err
	}
	if err := GitRepo.Add(filePath); err != nil {
		return err
	}
	if pushChanges {
		if err := GitRepo.Commit(commitMsgVerBump + ver.String()); err != nil {
			return err
		}
		if err := GitRepo.Push(""); err != nil {
			return err
		}
	}
	output.Logger().WithFields(logrus.Fields{
		"file":              f.Path,
		"fileChangesPushed": pushChanges,
	}).Info("the file has been updated")
	return nil
}
