package version

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrParseScopeName = errors.New("scope name can't be parsed")
)

type Scope struct {
	Id int
}

const (
	AUTO = iota
	MAJOR
	MINOR
	PATCH
	NONE
)

// Parse a string that contains a scope and set the Scope's id
func (s *Scope) Parse(scopeToParse string) error {
	scopeToParse = strings.ToLower(scopeToParse)
	switch scopeToParse {
	case "auto":
		s.Id = AUTO
	case "major":
		s.Id = MAJOR
	case "minor":
		s.Id = MINOR
	case "patch":
		s.Id = PATCH
	case "none":
	case "":
		s.Id = NONE
	default:
		return fmt.Errorf("%v: %s", ErrParseScopeName, scopeToParse)
	}
	return nil
}

func (s Scope) String() string {
	switch s.Id {
	case MAJOR:
		return "major"
	case MINOR:
		return "minor"
	case PATCH:
		return "patch"
	case AUTO:
		return "auto"
	case NONE:
		return "none"
	}
	return "none"
}
