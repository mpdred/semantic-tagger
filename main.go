package main

import (
	"flag"
	"fmt"
	"log"

	"semtag/internal"
	"semtag/pkg/git"
	"semtag/pkg/output"
	"semtag/pkg/terminal"
	"semtag/pkg/version"
)

func parseFlags() internal.CliArgs {
	args := internal.CliArgs{}

	flag.StringVar(&args.Prefix, "prefix", "", `if set, append the prefix to the version number
		e.g.:
		input: ./semtag -prefix='api-'
		output: api-0.1.0`)
	flag.StringVar(&args.Suffix, "suffix", "", `if set, append the suffix to the version number
		e.g.:
		input: ./semtag -suffix='-rc'
		output: 0.1.0-rc
`)
	flag.StringVar(&args.CustomVersion, "version", "", `if set, use the user-provided version`)

	var versionScope string
	flag.StringVar(&versionScope, "increment", "", "if set, increment the version scope: auto | major | minor | patch")
	flag.BoolVar(&args.Push, "push", false, "if set, push the created/updated object(s): push the git tag; add, commit and push the updated file.")

	flag.BoolVar(&args.ShouldTagGit, "git-tag", false, "if set, create an annotated tag")

	flag.StringVar(&args.FilePath, "file", "", `a file that contains the version number (e.g. setup.py)`)
	flag.StringVar(&args.FileVerPattern, "file-version-pattern", "", `the pattern expected for the file version
	e.g.:
	cat setup.py
		setup(
		  name='my-project',
		  version='3.0.28',
		)

	input: ./semtag -increment auto -file=setup.py -file-version-pattern="version='%s',"
	output:
	cat setup.py
		setup(
		  name='my-project',
		  version='3.1.0',
		)
`)

	flag.StringVar(&args.ExecuteCommand, "command", "", `execute a shell command for all version tags: use %s as a placeholder for the version number
	e.g.:
	version tags: v5, v5.0, v5.0.3, v5.0.3-32b0262
	input: ./semtag -prefix='v' -command="docker tag $MY_IMAGE_NAME $MY_DOCKER_REGISTRY/app:%s" 
	output:
		sh -c docker tag $MY_IMAGE_NAME $MY_DOCKER_REGISTRY/app:v5
		sh -c docker tag $MY_IMAGE_NAME $MY_DOCKER_REGISTRY/app:v5.0
		sh -c docker tag $MY_IMAGE_NAME $MY_DOCKER_REGISTRY/app:v5.0.3
		sh -c docker tag $MY_IMAGE_NAME $MY_DOCKER_REGISTRY/app:v5.0.3-32b0262
`)
	flag.Parse()

	args.VersionScope.Parse(versionScope)

	return args
}

func main() {
	args := parseFlags()

	v := internal.GetVersion(args)

	shouldIncrementVersion := args.VersionScope.String() != version.EmptyScope
	if !shouldIncrementVersion {
		fmt.Print(v.String())
	} else {
		v.IncrementAuto(args.VersionScope.String())
		fmt.Print(v.String())
	}

	git.TrySetGitCredentialsBasicAuth()

	if args.ExecuteCommand != "" {
		for _, val := range v.AsList() {
			_, err := terminal.Shellf(args.ExecuteCommand, val)
			if err == terminal.ErrShellCommand {
				log.Panic(err)
			}
		}

	}

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
