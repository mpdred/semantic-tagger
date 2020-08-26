package versionControl

type VersionControl interface {
	// Commit the current changes
	Commit(msg string) error

	// Add a file to the git index
	Add(file string) error

	// Push the local commits to remote
	Push(target string) error

	// Tag the current commit with an annotated git tag
	TagCommit(tag string, message string) error

	// DescribeLong gives an object a human readable name based on an available ref
	DescribeLong() (string, error)

	// GetLatestTag returns the latest tag that adheres to the semantic versioning regex and that has the provided prefix and suffix
	GetLatestTag(prefix, baseRegex, suffix string) (string, error)

	/*
	   IsAlreadyTagged checks if the current commit has been tagged with the current version number
	   Rules:
	   	- if there is no tag for the commit then return false
	   	- if the current commit has a different tag than the current version then return false
	*/
	IsAlreadyTagged(ver string) bool

	// GetTagsHead retrieves all the tags for the HEAD commit
	GetTagsHead() (string, error)

	// GetHash returns the git has for the HEAD commit
	GetHash() (string, error)

	// GetLatestCommitLogs returns the latest n commit logs
	GetLatestCommitLogs(count int) (string, error)

	// Fetch downloads the objects and refs from the remote
	Fetch() error
}
