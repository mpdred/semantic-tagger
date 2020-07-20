package internal

type CliArgs struct {
	Prefix        string
	Suffix        string
	CustomVersion string

	Push         bool
	VersionScope string

	ShouldTagGit   bool
	FilePath       string
	FileVerPattern string

	ExecuteCommand string
}
