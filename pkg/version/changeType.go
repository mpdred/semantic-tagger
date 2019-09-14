package version

type ChangeType int

const (
	BREAKING = iota
	FEATURE
	PATCH
)

func (cType *ChangeType) String() string {
	return []string{"(breaking)", "(feature)", "(patch)"}[*cType]
}
