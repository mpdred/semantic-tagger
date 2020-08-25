package version

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrParseScopeName = errors.New("scope name can't be parsed")
	ErrParseScopeId   = errors.New("scope id can't be parsed")
)

type Scope struct {
	Id int
}

const (
	AUTO = iota
	MAJOR
	MINOR
	PATCH
)

// Parse a string that contains a scope and set the Scope's id
func (s *Scope) Parse(scopeToParse string) error {
	scopeToParse = strings.ToLower(scopeToParse)
	switch scopeToParse {
	case "":
	case "auto":
		s.Id = AUTO
	case "major":
		s.Id = MAJOR
	case "minor":
		s.Id = MINOR
	case "patch":
		s.Id = PATCH
	default:
		return errors.New(fmt.Sprintf("%v: %s", ErrParseScopeName, scopeToParse))
	}
	return nil
}

func (s *Scope) String() string {
	switch s.Id {
	case MAJOR:
		return "major"
	case MINOR:
		return "minor"
	case PATCH:
		return "patch"
	case AUTO:
		return "auto"
	default:
		return "auto"
	}
}
