package cli

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"gnorm.org/gnorm/environ"
)

func makeEnv() (stderr, stdout *bytes.Buffer, env environ.Values) {
	stderr = &bytes.Buffer{}
	stdout = &bytes.Buffer{}
	return stderr, stdout, environ.Values{
		Stdout: stdout,
		Stderr: stderr,
		Log:    log.New(stdout, "", 0),
		Env:    map[string]string{},
	}
}

func TestPreviewBadFormat(t *testing.T) {
	stderr, stdout, env := makeEnv()
	env.Args = []string{"preview", "--format", "abc123"}
	code := ParseAndRun(env)
	if code != 2 {
		t.Fatalf("Expected bad format to produce code 2, but got code %v", code)
	}
	out := stdout.String()
	err := stderr.String()
	if out != "" {
		t.Fatalf("unexpected non-empty output for error case: %q", out)
	}
	if !strings.Contains(err, `unknown preview format "abc123"`) {
		t.Fatalf("unexpected error message from preview:\n%s", err)
	}
}
