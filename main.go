package main

import (
	"errors"
	"fmt"
	"log"

	"semtag/internal"
	"semtag/pkg/git"
	"semtag/pkg/output"
	"semtag/pkg/terminal"
	"semtag/pkg/version"
)

var (
	ErrNotPushMode = errors.New("push to Git skipped: use the `-push` flag to push changes")
)

func main() {
	args := internal.CliArgs{}
	args.ParseFlags()

	v := version.Version{}
	if args.CustomVersion == internal.EmptyStringFlag {
		v.GetLatestFromGit()
	} else {
		v.UseVersionProvidedByUser(args.Prefix, args.CustomVersion, args.Suffix)
	}

	shouldIncrementVersion := args.VersionScope.String() != version.EmptyScope
	if !shouldIncrementVersion {
		fmt.Print(v.String())
	} else {
		v.IncrementAuto(args.VersionScope.String())
		fmt.Print(v.String())
	}

	git.TryConfigureIdentity()
	git.TrySetGitCredentialsBasicAuth()

	if args.ExecuteCommand != "" {
		for _, val := range v.AsList() {
			_, err := terminal.Shellf(args.ExecuteCommand, val)
			if err != nil {
				log.Fatal(err)
			}
		}

	}

	if args.ShouldTagGit {
		if internal.HasRelevantChanges(args.RelevantPaths) {
			internal.TagGit(v, args.Push)
			if !args.Push {
				output.Debug(ErrNotPushMode)
			}
		}
	}

	shouldTagInFile := len(args.FilePath) > 0 && len(args.FileVersionPattern) > 0
	if shouldTagInFile {
		internal.TagFile(v, args.FilePath, args.FileVersionPattern, args.Push)
		if !args.Push {
			output.Debug(ErrNotPushMode)
		}
	}
}
