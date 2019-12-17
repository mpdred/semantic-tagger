package version

type ChangeType int

const (
	MAJOR = iota
	MINOR
	PATCH
)

func (cType *ChangeType) String() string {
	return []string{"(major)", "(minor)", "(patch)"}[*cType]
}
