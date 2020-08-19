package internal

import (
	"errors"
	"flag"
	"fmt"

	"semtag/pkg/changelog"
	"semtag/pkg/output"
	"semtag/pkg/version"
)

const (
	EmptyStringFlag     = ""
	DefaultRelevantPath = "."
)

var (
	ErrOnlyOneDeployFileFlagsSet = errors.New("to update a file please use both the '-file' and the '-file-pattern' flags")
)

type CliArgs struct {
	Prefix        string
	Suffix        string
	CustomVersion string
	VersionScope  version.Scope

	RelevantPaths relevantPaths

	Push           bool
	ShouldTagGit   bool
	ExecuteCommand string

	Changelog      bool
	ChangelogRegex string

	FilePath           string
	FileVersionPattern string
}

type relevantPaths []string

func (i *relevantPaths) String() string {
	return DefaultRelevantPath
}

func (i *relevantPaths) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func (args *CliArgs) ParseFlags() {

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

	flag.StringVar(&args.ChangelogRegex, "changelog-regex", changelog.ChangelogDefaultRegex, "if set, generate the changelog only for specific tags")
	flag.BoolVar(&args.Changelog, "changelog", false, fmt.Sprintf(`if set, generate a full changelog for the repository. In order to have correct hyperlinks you will need to provide two environment variables for your web-based git repository: %s for the URL of the commits and %s for the URL of the tags
	e.g.:
	input: GIT_COMMIT_URL="https://gitlab.com/my_org/my_group/my_repository/-/commit/" GIT_TAG_URL="https://gitlab.com/my_org/my_group/my_repository/-/tags/" ./semtag -generate-changelog
	output: a full repository changelog in a Markdown file that shows the commit name(s) included in each tag
`, changelog.EnvVarGitCommitUrl, changelog.EnvVarGitTagUrl))

	flag.StringVar(&args.FilePath, "file", "", `a file that contains the version number (e.g. setup.py)`)
	flag.StringVar(&args.FileVersionPattern, "file-version-pattern", "", `the pattern expected for the file version
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

	flag.Var(&args.RelevantPaths, "path", `if set, create a git tag only if changes are detected in the provided path(s)
	e.g.:
	input: ./semtag -path="src" -path="lib/" -path="Dockerfile"
`)

	flag.Parse()

	if (args.FilePath == EmptyStringFlag) != (args.FileVersionPattern == EmptyStringFlag) {
		output.Logger().Fatalln(ErrOnlyOneDeployFileFlagsSet)
	}
	if len(args.RelevantPaths) == 0 {
		args.RelevantPaths = relevantPaths{DefaultRelevantPath}
	}

	args.VersionScope.Parse(versionScope)
}
