package main

import (
	"errors"
	"fmt"

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

func main() {
	args := internal.CliArgs{}
	args.ParseFlags()

	v := setVersion(args)

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
		hasRelevantChanges, err := internal.HasRelevantChanges(args.RelevantPaths)
		if err != nil {
			output.Logger().Fatal(err)
		}

		if hasRelevantChanges {
			tag := &versionControl.Tag{
				Name: v.String(),
			}
			if err := internal.TagGit(tag, args.Push); err != nil {
				output.Logger().Fatal(err)
			}
			if !args.Push {
				output.Logger().Warn(ErrNotPushMode)
			}
		}
	}

	shouldTagInFile := len(args.FileName) > 0 && len(args.FileVersionPattern) > 0
	if shouldTagInFile {
		if err := internal.TagFile(v, args.FileName, args.FileVersionPattern, args.Push); err != nil {
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

	// print the version to stdout; execute as the last command so that it can be grepped by simple shell scripts
	fmt.Print(v.String())
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

	shouldIncrementVersion := args.VersionScope.String() != ""
	if shouldIncrementVersion {
		if err := v.SetScope(args.VersionScope.String()); err != nil {
			output.Logger().Fatal(err)
		}
	}
	return v
}
