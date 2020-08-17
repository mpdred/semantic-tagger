package internal

import "semtag/pkg/version"

type CliArgs struct {
	Prefix        string
	Suffix        string
	CustomVersion string

	Push         bool
	VersionScope version.Scope

	ShouldTagGit   bool
	FilePath       string
	FileVerPattern string

	ExecuteCommand string
}
