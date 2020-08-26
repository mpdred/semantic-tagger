package internal

import (
	"errors"
	"flag"
	"fmt"

	"github.com/sirupsen/logrus"

	"semtag/pkg/changelog"
	"semtag/pkg/output"
	"semtag/pkg/version"
	"semtag/pkg/versionControl"
)

const (
	binaryName = "semtag"

	flagPrefix    = "prefix"
	flagSuffix    = "suffix"
	flagIncrement = "increment"
	flagVersion   = "version"

	flagPath = "path"

	flagShouldTagGit   = "git-tag"
	flagShouldPush     = "push"
	flagExecuteCommand = "command"

	flagFileName           = "file"
	flagFileVersionPattern = "file-version-pattern"

	flagChangelog      = "changelog"
	flagChangelogRegex = "changelog-regex"
)

var (
	errMissingArgs = errors.New("required arguments not found")
)

type CliArgs struct {
	Prefix        string
	Suffix        string
	CustomVersion string
	VersionScope  version.Scope

	RelevantPaths versionControl.RelevantPaths

	Push           bool
	ShouldTagGit   bool
	ExecuteCommand string

	Changelog      bool
	ChangelogRegex string

	FileName           string
	FileVersionPattern string
}

var tempVersionScope string

func (args *CliArgs) ParseFlags() {
	args.loadAllFlags()
	args.parseAndInit()
	args.guardAgainstInvalidArgs()
}

func (args *CliArgs) parseAndInit() {
	flag.Parse()

	if len(args.RelevantPaths) == 0 {
		args.RelevantPaths = versionControl.RelevantPaths{versionControl.DefaultRelevantPath}
	}
	if err := args.VersionScope.Parse(tempVersionScope); err != nil {
		output.Logger().Fatal(err)
	}

	output.Logger().WithField("args", fmt.Sprintf("%#v", args)).Info("arguments parsed")
}

func (args *CliArgs) loadAllFlags() {
	flag.Var(
		&args.RelevantPaths,
		flagPath,
		fmt.Sprintf(`if set, create a git tag only if changes are detected in the provided path(s)
	e.g.:
	$ ./%s -%[2]s="src" -%[2]s="lib/" -%[2]s="Dockerfile"
`,
			binaryName, flagPath))

	args.loadGenericVersionFlags()
	args.loadBaseActionFlags()
	args.loadChangelogFlags()
	args.loadFileActionFlags()
}

func (args *CliArgs) loadFileActionFlags() {
	flag.StringVar(
		&args.FileName,
		flagFileName,
		"",
		`a file that contains the version number (e.g. setup.py)`)
	flag.StringVar(&args.FileVersionPattern, "file-version-pattern", "", `the pattern expected for the file version
	e.g.:
	$ cat setup.py
		setup(
		  name='my-project',
		  version='3.0.28',
		)
`+fmt.Sprintf(`
	$ ./%s -%s=auto -%s=setup.py -%s="version='%%s',"
	$ cat setup.py
		setup(
		  name='my-project',
		  version='3.1.0',
		)
`,
		binaryName, flagIncrement, flagFileName, flagFileVersionPattern))

}

func (args *CliArgs) loadChangelogFlags() {
	flag.StringVar(
		&args.ChangelogRegex,
		flagChangelogRegex,
		changelog.DefaultRegexFormat,
		"if set, generate the changelog only for specific tags")

	flag.BoolVar(
		&args.Changelog,
		flagChangelog,
		false,
		fmt.Sprintf(`if set, generate a full changelog for the repository. In order to have correct hyperlinks you will need to provide two environment variables for your web-based git repository: %[1]s for the URL of the commits and %[2]s for the URL of the tags
	e.g.:
	$ %[1]s="https://gitlab.com/my_org/my_group/my_repository/-/commit/" %[2]s="https://gitlab.com/my_org/my_group/my_repository/-/tags/" ./%s -%s
	output: a full repository changelog in a file (%s) that shows the commit name(s) included in each tag
`,
			changelog.EnvVarGitCommitUrl, changelog.EnvVarGitTagUrl, binaryName, flagChangelog, changelog.DefaultChangelogFile))

}

func (args *CliArgs) loadBaseActionFlags() {
	flag.BoolVar(
		&args.Push,
		flagShouldPush,
		false,
		"if set, push the created/updated object(s): push the git tag AND/OR add, commit and push the updated file")

	flag.BoolVar(
		&args.ShouldTagGit,
		flagShouldTagGit,
		false,
		"if set, create an annotated tag")

	flag.StringVar(
		&args.ExecuteCommand,
		flagExecuteCommand,
		"",
		`execute a shell command for all version tags: use %s as a placeholder for the version number
	e.g.: version tags: v5, v5.0, v5.0.3, v5.0.3-32b0262`+fmt.Sprintf(`

	$ ./%s -%s='v' -%s="docker tag $MY_IMAGE_NAME $MY_DOCKER_REGISTRY/app:%%s"
		docker tag $MY_IMAGE_NAME $MY_DOCKER_REGISTRY/app:v5
		docker tag $MY_IMAGE_NAME $MY_DOCKER_REGISTRY/app:v5.0
		docker tag $MY_IMAGE_NAME $MY_DOCKER_REGISTRY/app:v5.0.3
		docker tag $MY_IMAGE_NAME $MY_DOCKER_REGISTRY/app:v5.0.3-32b0262
`,
			binaryName, flagPrefix, flagExecuteCommand))

}

func (args *CliArgs) loadGenericVersionFlags() {
	flag.StringVar(
		&args.Prefix,
		flagPrefix,
		"",
		fmt.Sprintf(`if set, append the prefix to the version number
	e.g.:
	$ ./%s -%s='api-'
	api-0.1.0
`,
			binaryName, flagPrefix))

	flag.StringVar(
		&args.Suffix,
		flagSuffix,
		"",
		fmt.Sprintf(`if set, append the suffix to the version number
	e.g.:
	$ ./%s -%s='-rc'
	0.1.0-rc
`,
			binaryName, flagSuffix))

	flag.StringVar(
		&args.CustomVersion,
		flagVersion,
		"",
		`if set, use the provided version`)

	var versionScope string
	flag.StringVar(
		&versionScope,
		flagIncrement,
		"",
		"if set, increment the version scope: [ auto | major | minor | patch ]")
	tempVersionScope = versionScope

}

func (args *CliArgs) guardAgainstInvalidArgs() {
	if (args.FileName == "") != (args.FileVersionPattern == "") {
		output.Logger().WithFields(logrus.Fields{
			"flags": []string{flagFileName, flagFileVersionPattern},
		}).Fatalln(errMissingArgs)
	}
}
