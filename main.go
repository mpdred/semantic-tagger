package main

import (
	"flag"
	"fmt"

	"semtag/internal"
	"semtag/pkg/git"
	"semtag/pkg/output"
)

func parseFlags() internal.CliArgs {
	args := internal.CliArgs{}

	flag.StringVar(&args.Prefix, "prefix", "", `if set, append the prefix to the version number (e.g. "api-0.1.0")`)
	flag.StringVar(&args.Suffix, "suffix", "", `if set, append the suffix to the version number (e.g. "0.1.0-rc")`)
	flag.StringVar(&args.CustomVersion, "version", "", `if set, use the user-provided version`)

	flag.StringVar(&args.VersionScope, "increment", "", "if set, increment the version scope: auto | major | minor | patch")
	flag.BoolVar(&args.Push, "push", false, "if set, push the created/updated object(s)")

	flag.BoolVar(&args.ShouldTagGit, "git-tag", false, "if set, create an annotated tag")

	flag.StringVar(&args.FilePath, "file", "", `a file that contains the version number (e.g. "setup.py")`)
	flag.StringVar(&args.FileVerPattern, "file-version-pattern", "%s", `the pattern expected for the file version (e.g. "version='%s',")`)

	flag.Parse()

	return args
}

func main() {
	args := parseFlags()

	v := internal.GetVersion(args)

	shouldIncrementVersion := args.VersionScope != ""
	if !shouldIncrementVersion {
		fmt.Print(v.String())
	} else {
		v.IncrementAuto(args.VersionScope)
		fmt.Print(v.String())
	}

	git.TrySetGitCredentialsBasicAuth()

	notPushModeMessage := "push to Git skipped: use the `-push` flag to push changes"

	if args.ShouldTagGit {
		internal.TagGit(v, args.Push)
		if !args.Push {
			output.Debug(notPushModeMessage)
		}
	}

	shouldTagInFile := len(args.FilePath) > 0 && len(args.FileVerPattern) > 0
	if shouldTagInFile {
		internal.TagFile(v, args.FilePath, args.FileVerPattern, args.Push)
		if !args.Push {
			output.Debug(notPushModeMessage)
		}
	}
}
