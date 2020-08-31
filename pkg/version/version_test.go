package version

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"semtag/pkg/versionControl"
)

func Test_Increment(t *testing.T) {
	// arrange
	GitRepo = &versionControl.GitRepositoryMock{}
	tables := []struct {
		ver   string
		scope Scope
		want  string
	}{
		{"1.2.3", Scope{MAJOR}, "2.0.0"},
		{"1.2.3", Scope{MINOR}, "1.3.0"},
		{"1.2.3", Scope{PATCH}, "1.2.4"},
		{"1.2.3", Scope{NONE}, "1.2.3"},
	}
	assertCorrectMessage := func(t *testing.T, got, want string, changeType Scope) {
		t.Helper()
		if got != want {
			t.Errorf("got %q want %q %s", got, want, changeType.String())
		}
	}

	// act
	for _, tb := range tables {
		t.Run(fmt.Sprintf("ver=%q, scope=%s, want=%q", tb.ver, tb.scope.String(), tb.want), func(t *testing.T) {
			ver := Version{}

			if err := ver.load(tb.ver); err != nil {
				t.Error(err)
			}
			if err := ver.Increment(tb.scope); err != nil {
				t.Error(err)
			}

			// assert
			got := ver.String()
			assertCorrectMessage(t, got, tb.want, tb.scope)
		})
	}
}

func Test_AsList(t *testing.T) {
	// arrange
	GitRepo = &versionControl.GitRepositoryMock{}
	tables := []struct {
		gitSha       string
		expectedList []string
	}{
		{"a3ed223b", []string{"0.2.1-ga3ed223b", "0.2.1", "0.2", "0"}},
	}

	ver := Version{}
	if err := ver.load("0.2.1"); err != nil {
		t.Error(err)
	}

	// act
	for _, tb := range tables {
		ver.Hash = tb.gitSha
		actualList := ver.AsList()

		// assert
		// compare the two slices
		expLen := len(tb.expectedList)
		actLen := len(actualList)
		if expLen != actLen {
			t.Errorf("expected length of %d but found %d", expLen, actLen)
		}
		if !reflect.DeepEqual(tb.expectedList, actualList) {
			t.Errorf("expected %v but found %v", tb.expectedList, actualList)
		}
	}
}

func Test_String(t *testing.T) {
	// arrange
	GitRepo = &versionControl.GitRepositoryMock{}

	tables := []struct {
		prefix string
		suffix string

		want string
	}{
		{"v", "", "v0.1.0"},
		{"", "-api", "0.1.0-api"},
		{"v", "-api", "v0.1.0-api"},
	}
	assertCorrectMessage := func(t *testing.T, got, want string) {
		t.Helper()
		if got != want {
			t.Errorf("got %q want %q ", got, want)
		}
	}

	// act
	for _, tb := range tables {
		t.Run(fmt.Sprintf("prefix=%q, suffix=%q, want=%q", tb.prefix, tb.suffix, tb.want), func(t *testing.T) {
			vNumber := "0.1.0"
			ver := Version{}
			ver.Prefix = tb.prefix
			ver.Suffix = tb.suffix
			if err := ver.load(vNumber); err != nil {
				t.Error(err)
			}

			// assert
			if tb.want != ver.String() {
				t.Errorf("expected %q but found %q", tb.want, ver.String())
			}

			got := ver.String()
			assertCorrectMessage(t, got, tb.want)
		})
	}
}

func Test_Load(t *testing.T) {
	// arrange
	GitRepo = &versionControl.GitRepositoryMock{}

	tables := []struct {
		raw string

		want Version
	}{
		{"1.2.3", Version{Major: 1, Minor: 2, Patch: 3}},
		{"11.2.3", Version{Major: 11, Minor: 2, Patch: 3}},
		{"1.22.3", Version{Major: 1, Minor: 22, Patch: 3}},
		{"1.2.33", Version{Major: 1, Minor: 2, Patch: 33}},
	}
	assertCorrectVersion := func(t *testing.T, got, want Version) {
		t.Helper()
		if got.String() != want.String() {
			t.Errorf("got %q want %q ", got, want)
		}
	}

	// act
	for _, tb := range tables {
		t.Run(fmt.Sprintf("raw=%q, want=%s", tb.raw, tb.want.String()), func(t *testing.T) {
			v := Version{}
			if err := v.load(tb.raw); err != nil {
				t.Error(err)
			}
			assertCorrectVersion(t, v, tb.want)
		})
	}
}

func Test_LoadNegative(t *testing.T) {
	// arrange
	GitRepo = &versionControl.GitRepositoryMock{}

	tables := []struct {
		raw string

		want string
	}{
		{"", ErrParseVersionMajor.Error()},
		{"foo", ErrParseVersionMajor.Error()},
		{"a.2.3", ErrParseVersionMajor.Error()},
		{"1.b.3", ErrParseVersionMinor.Error()},
		{"1.2.c", ErrParseVersionPatch.Error()},
	}
	assertCorrectVersion := func(t *testing.T, got, want string) {
		t.Helper()
		if !strings.Contains(got, want) {
			t.Errorf("got %q want %q ", got, want)
		}
	}

	// act
	for _, tb := range tables {
		t.Run(fmt.Sprintf("raw=%q, want=%q", tb.raw, tb.want), func(t *testing.T) {
			v := Version{}

			err := v.load(tb.raw)

			// assert
			if err == nil {
				t.Error(err)
			}
			assertCorrectVersion(t, err.Error(), tb.want)
		})
	}
}

func Test_Validate(t *testing.T) {
	// arrange
	GitRepo = &versionControl.GitRepositoryMock{}

	tables := []struct {
		version string

		wantError bool
	}{
		{"1.2.3", false},
		{"foo", true},
		{"1-2.3", true},
		{"1..2.3", true},
		{"v1.2.3", true},
		{"1.2.3-api", true},
	}
	assertCorrectVersion := func(t *testing.T, got, want bool) {
		t.Helper()
		if got != want {
			t.Errorf("got %t want %t ", got, want)
		}
	}

	// act
	for _, tb := range tables {
		t.Run(fmt.Sprintf("raw=%q, wantError=%t", tb.version, tb.wantError), func(t *testing.T) {
			v := Version{}

			err := v.Validate(tb.version)

			// assert
			got := err != nil
			assertCorrectVersion(t, got, tb.wantError)
		})
	}
}

func Test_Parse(t *testing.T) {
	// arrange
	GitRepo = &versionControl.GitRepositoryMock{}

	tables := []struct {
		version Version

		wantError bool
	}{
		{Version{Major: 1, Minor: 2, Patch: 3}, false},
		{Version{Major: 1, Minor: 2, Patch: 3, Suffix: "-api"}, false},
		{Version{Prefix: "v", Major: 1, Minor: 2, Patch: 3}, false},
	}
	assertCorrectVersion := func(t *testing.T, got, want bool) {
		t.Helper()
		if got != want {
			t.Errorf("got %t want %t ", got, want)
		}
	}

	// act
	for _, tb := range tables {
		t.Run(fmt.Sprintf("raw=%q, wantError=%t", tb.version.String(), tb.wantError), func(t *testing.T) {
			v := Version{
				Prefix: tb.version.Prefix,
				Suffix: tb.version.Suffix,
			}

			err := v.Parse(tb.version.String())

			// assert
			got := err != nil
			assertCorrectVersion(t, got, tb.wantError)
		})
	}
}

func Test_SetScope(t *testing.T) {
	// arrange
	GitRepo = &versionControl.GitRepositoryMock{}

	baseVer := Version{
		Major: 1,
		Minor: 2,
		Patch: 3,
	}
	tables := []struct {
		scope string

		want Version
	}{
		{"major", Version{Major: 2, Minor: 0, Patch: 0}},
		{"minor", Version{Major: 1, Minor: 3, Patch: 0}},
	}

	// act
	for _, tb := range tables {
		t.Run(fmt.Sprintf("version=%q, scope=%q, want=%q", baseVer.String(), tb.scope, tb.want.String()), func(t *testing.T) {

			ver := baseVer
			if err := ver.SetIncrementScope(tb.scope); err != nil {
				t.Error(err)
			}

			// assert
			assertCorrectVersion := func(t *testing.T, got, want string) {
				t.Helper()
				if got != want {
					t.Errorf("got %q want %q ", got, want)
				}
			}
			assertCorrectVersion(t, ver.String(), tb.want.String())
		})
	}
}
