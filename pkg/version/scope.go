package version

import (
	"errors"
	"strings"

	"semtag/pkg"
	"semtag/pkg/output"
)

const EmptyScope string = "undefined"

var ErrParseScopeName = errors.New("scope name can't be parsed")
var ErrParseScopeId = errors.New("scope id can't be parsed")

type Scope struct {
	Id int
}

const (
	EMPTY = iota
	AUTO
	MAJOR
	MINOR
	PATCH
)

func (s *Scope) Parse(scopeToParse string) {
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
		output.Logger().Fatal(pkg.NewErrorDetails(ErrParseScopeName, scopeToParse))
	}
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
	case EMPTY:
		return EmptyScope
	default:
		output.Logger().Fatal(pkg.NewErrorDetails(ErrParseScopeId, s.Id))
		return EmptyScope
	}
}
