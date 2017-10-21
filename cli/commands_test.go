package cli

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

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

func TestInitCmd(t *testing.T) {
	d, err := ioutil.TempDir("", "gnormInitTest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(d)
	if err := initFunc(d); err != nil {
		t.Fatalf("error running initfunc: %v", err)
	}
	cfgFile := filepath.Join(d, "gnorm.toml")
	b, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		t.Errorf("error reading gnorm.toml: %v", err)
	} else {
		if diff := cmp.Diff(sample, string(b)); diff != "" {
			t.Fatal("gnorm.toml differs from expected data:\n" + diff)
		}
	}
	if _, err := os.Stat(filepath.Join(d, "static")); err != nil {
		t.Errorf("error finding static dir: %v", err)
	}
	fs, err := ioutil.ReadDir(filepath.Join(d, "templates"))
	if err != nil {
		t.Fatal(err)
	}
	if len(fs) != 3 {
		t.Fatalf("expected 3 files in templates dir, but got %d", len(fs))
	}
	names := map[string]bool{}
	for _, f := range fs {
		names[f.Name()] = true
	}
	if !names["table.gotmpl"] {
		t.Errorf("missing table template")
	}
	if !names["schema.gotmpl"] {
		t.Errorf("missing schema template")
	}
	if !names["enum.gotmpl"] {
		t.Errorf("missing enum template")
	}
}
