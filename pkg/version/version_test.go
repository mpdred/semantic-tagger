package version

import (
	"reflect"
	"testing"
)

func Test_Increment(t *testing.T) {
	tables := []struct {
		ver   string
		scope Scope
		want  string
	}{
		{"1.2.3", Scope{MAJOR}, "2.0.0"},
		{"1.2.3", Scope{MINOR}, "1.3.0"},
		{"1.2.3", Scope{PATCH}, "1.2.4"},
	}
	assertCorrectMessage := func(t *testing.T, got, want string, changeType Scope) {
		t.Helper()
		if got != want {
			t.Errorf("got %q want %q %s", got, want, changeType.String())
		}
	}
	for _, tb := range tables {
		t.Run("Test Increment", func(t *testing.T) {
			ver := Version{}
			ver.Load(tb.ver)
			ver.Increment(tb.scope)
			got := ver.String()
			assertCorrectMessage(t, got, tb.want, tb.scope)
		})
	}
}

func Test_AsList(t *testing.T) {
	tables := []struct {
		gitSha       string
		expectedList []string
	}{
		{"g3ed223b", []string{"0.2.1-g3ed223b", "0.2.1", "0.2", "0"}},
	}

	ver := Version{}
	ver.Load("0.2.1")
	for _, tb := range tables {
		ver.Hash = tb.gitSha
		actualList := ver.AsList()

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
	tables := []struct {
		prefix         string
		suffix         string
		expectedOutput string
	}{
		{"v", "", "v0.1.0"},
		{"", "-api", "0.1.0-api"},
		{"v", "-api", "v0.1.0-api"},
	}
	for _, tb := range tables {
		vNumber := "0.1.0"
		ver := Version{}
		ver.Prefix = tb.prefix
		ver.Suffix = tb.suffix
		ver.Load(vNumber)
		if tb.expectedOutput != ver.String() {
			t.Errorf("expected %v but found %v", tb.expectedOutput, ver.String())
		}
	}
}
