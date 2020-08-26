package terminal

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func Test_GetEnv(t *testing.T) {
	// arrange
	tables := []struct {
		key    string
		exists bool

		want string
	}{
		{"foo", true, "bar"},
		{"foo", false, ErrEnvVarNotFound.Error()},
	}
	assertCorrectMessage := func(t *testing.T, got, want string) {
		t.Helper()
		if !strings.Contains(got, want) {
			t.Errorf("got %q want %q", got, want)
		}
	}

	// act
	for _, tb := range tables {
		t.Run(fmt.Sprintf("Test GetEnv: key=%q, exists=%t, want=%q", tb.key, tb.exists, tb.want), func(t *testing.T) {
			if tb.exists {
				if err := os.Setenv(tb.key, tb.want); err != nil {
					t.Error(err)
				}
			} else {
				if err := os.Unsetenv(tb.key); err != nil {
					t.Error(err)
				}
			}
			out, err := GetEnv(tb.key)

			// assert
			var got string
			got = out
			if err != nil {
				got = err.Error()
			}

			assertCorrectMessage(t, got, tb.want)
		})
	}
}

func Test_ShellRaw(t *testing.T) {
	// arrange
	tables := []struct {
		cmd string

		want string
	}{
		{"echo 'hello kitty'", "hello kitty"},
		{"echo -e 'hello kitty'\n", "hello kitty\n"},
		{"echo 'hello kitty", ErrShellCommand.Error()},
		{"aaaa", ErrShellCommand.Error()},
	}
	assertCorrectMessage := func(t *testing.T, got, want string) {
		t.Helper()
		if !strings.Contains(got, want) {
			t.Errorf("got %q want %q", got, want)
		}
	}

	// act
	for _, tb := range tables {
		t.Run("Test ShellRaw", func(t *testing.T) {
			out, err := ShellRaw(tb.cmd)

			// assert
			var got string
			got = out
			if err != nil {
				got = err.Error()
			}

			assertCorrectMessage(t, got, tb.want)
		})
	}
}
func Test_Shell(t *testing.T) {
	// arrange
	tables := []struct {
		cmd string

		want string
	}{
		{"echo -e 'hello kitty'\n", "hello kitty\n"},
		{"echo 'hello kitty", ErrShellCommand.Error()},
		{"aaaa", ErrShellCommand.Error()},
	}
	assertCorrectMessage := func(t *testing.T, got, want string) {
		t.Helper()
		if !strings.Contains(got, want) {
			t.Errorf("got %q want %q", got, want)
		}
	}

	// act
	for _, tb := range tables {
		t.Run("Test Shel", func(t *testing.T) {
			out, err := ShellRaw(tb.cmd)

			// assert
			var got string
			got = out
			if err != nil {
				got = err.Error()
			}

			assertCorrectMessage(t, got, tb.want)
		})
	}
}
