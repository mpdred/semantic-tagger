package internal

type relevantPaths []string

func (i *relevantPaths) String() string {
	return DefaultRelevantPath
}

func (i *relevantPaths) Set(value string) error {
	*i = append(*i, value)
	return nil
}
