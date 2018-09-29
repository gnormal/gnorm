package run

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"text/template"

	"gnorm.org/gnorm/environ"
)

func TestMain(m *testing.M) {
	if os.Getenv("GNORM_RUNHELPER") == "" {
		os.Exit(m.Run())
	}
	testEngine()
}

func TestAtomicGenerate(t *testing.T) {
	target := OutputTarget{
		Filename: template.Must(template.New("").Parse("{{.}}")),
		// the contents tempalte will fail to execute because the contents will
		// not have a .Name field.
		Contents: template.Must(template.New("").Parse("{{.Name}}")),
	}
	env := environ.Values{
		Log: log.New(ioutil.Discard, "", 0),
	}
	filename := "testfile.out"
	original := []byte("goodbye world")
	err := ioutil.WriteFile(filename, original, 0600)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(filename)
	contents := "hello world"
	err = genFile(env, filename, contents, target, nil, nil, ".", templateEngine{})
	if err == nil {
		t.Fatal("Unexpected nil error generating contents. Should have failed.")
	}
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, original) {
		t.Fatalf("Expected file to be unchanged, but was different.  Expected: %q, got: %q", original, b)
	}
}

func TestCopyStaticFiles(t *testing.T) {
	originPaths := []string{
		"base/base.md",
		"base/level_one/level_one.md",
		"base/level_two/level_two.md",
	}
	source := "testdata"
	dest := "static_asset"

	err := copyStaticFiles(environ.Values{}, source, dest)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := os.RemoveAll(dest)
		if err != nil {
			t.Fatal(err)
		}
	}()

	// make sure the structure is preserved
	var newPaths []string
	filepath.Walk(dest, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		r, err := filepath.Rel(dest, path)
		if err != nil {
			return err
		}
		newPaths = append(newPaths, r)
		return nil
	})

	sort.Strings(originPaths)
	sort.Strings(newPaths)

	if !reflect.DeepEqual(originPaths, newPaths) {
		t.Errorf("expected %v to equal %v", newPaths, originPaths)
	}
}

func TestNoOverwriteGlobs(t *testing.T) {
	target := OutputTarget{
		Filename: template.Must(template.New("").Parse("{{.}}")),
		Contents: template.Must(template.New("").Parse("{{.}}")),
	}
	env := environ.Values{
		Log: log.New(ioutil.Discard, "", 0),
	}

	filename := "testfile.out"

	t.Run("file exists and matches glob", func(t *testing.T) {
		original := []byte("goodbye world")
		err := ioutil.WriteFile(filename, original, 0600)
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(filename)

		err = genFile(env, filename, "hello world", target, []string{"*.out"}, nil, ".", templateEngine{})
		if err != nil {
			t.Fatalf("Unexpected error generating contents: %s", err)
		}

		b, err := ioutil.ReadFile(filename)
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(b, original) {
			t.Fatalf("Expected file to be unchanged, but was different.  Expected: %q, got: %q", original, b)
		}

		t.Run("does not match glob", func(t *testing.T) {
			content := "hello world"
			err = genFile(env, filename, content, target, []string{"bob"}, nil, ".", templateEngine{})
			if err != nil {
				t.Fatalf("Unexpected error generating contents: %s", err)
			}

			b, err := ioutil.ReadFile(filename)
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(b, []byte(content)) {
				t.Fatalf("Expected file to contain content, but did not.  Expected: %q, got: %q", content, b)
			}
		})
	})

	t.Run("file does not exist but matches glob", func(t *testing.T) {
		if _, err := os.Stat(filename); err == nil {
			t.Fatalf("File should not exist, but does.")
		}

		content := "hello world"
		err := genFile(env, filename, content, target, []string{"*.out"}, nil, ".", templateEngine{})
		if err != nil {
			t.Fatalf("Unexpected error generating contents: %s", err)
		}
		defer os.Remove(filename)

		b, err := ioutil.ReadFile(filename)
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(b, []byte(content)) {
			t.Fatalf("Expected file to contain content, but did not.  Expected: %q, got: %q", content, b)
		}
	})
}

func testEngine() {
	file := os.Getenv("GNORM_ARGSFILE")
	args := strings.Join(os.Args[1:], "\n")
	if err := ioutil.WriteFile(file, []byte(args), 0600); err != nil {
		panic(err)
	}
	if stdin := os.Getenv("GNORM_STDINFILE"); stdin != "" {
		f, err := os.Create(stdin)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		if _, err := io.Copy(f, os.Stdin); err != nil {
			panic(err)
		}
	}
	if stdout := os.Getenv("GNORM_STDOUT"); stdout != "" {
		if _, err := os.Stdout.WriteString(stdout); err != nil {
			panic(err)
		}
	}
	if datafile := os.Getenv("GNORM_COPYARG1FILE"); datafile != "" {
		b, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile(datafile, b, 0600)
		if err != nil {
			panic(err)
		}
	}
}

func TestTemplateEngine(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	argsfile := filepath.Join(dir, "argsfile")
	datafile := filepath.Join(dir, "datafile")
	env := map[string]string{
		"GNORM_RUNHELPER":    "1",
		"GNORM_ARGSFILE":     argsfile,
		"GNORM_COPYARG1FILE": datafile,
	}
	data := map[string]interface{}{
		"name": "Bob",
		"age":  22,
	}
	cmd := []string{os.Args[0], "{{.Data}}", "{{.Output}}", "{{.Template}}"}

	var templates []*template.Template
	for _, c := range cmd {
		tp, err := template.New("test").Funcs(environ.FuncMap).Parse(c)
		if err != nil {
			t.Fatal(err)
		}
		templates = append(templates, tp)
	}

	err = runExternalEngine(env, "outputpath", "templatepath", data, templateEngine{CommandLine: templates})
	if err != nil {
		t.Fatal(err)
	}
	b, err := ioutil.ReadFile(argsfile)
	if err != nil {
		t.Fatal(err)
	}
	args := strings.Split(string(b), "\n")
	if len(args) != 3 {
		t.Fatalf("expected 3 args, but got %q", args)
	}
	if args[1] != "outputpath" {
		t.Errorf("expected output path to be second arg, but was %q", args[1])
	}
	if args[2] != "templatepath" {
		t.Errorf("expected template path to be third arg, but was %q", args[2])
	}

	b, err = ioutil.ReadFile(datafile)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, expected) {
		t.Fatalf("expected data:\n%s\n\nactual data:\n%s", expected, b)
	}
}

func TestTemplateEngineStdinStdout(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	argsfile := filepath.Join(dir, "argsfile")
	datafile := filepath.Join(dir, "datafile")
	outputfile := filepath.Join(dir, "output")
	output := "SOMETHING TO OUTPUT"
	env := map[string]string{
		"GNORM_RUNHELPER": "1",
		"GNORM_ARGSFILE":  argsfile,
		"GNORM_STDINFILE": datafile,
		"GNORM_STDOUT":    output,
	}
	data := map[string]interface{}{
		"name": "Bob",
		"age":  22,
	}
	cmd := []string{os.Args[0], "{{.Output}}", "{{.Template}}"}

	var templates []*template.Template
	for _, c := range cmd {
		tp, err := template.New("test").Funcs(environ.FuncMap).Parse(c)
		if err != nil {
			t.Fatal(err)
		}
		templates = append(templates, tp)
	}

	err = runExternalEngine(env, outputfile, "templatepath", data, templateEngine{CommandLine: templates, UseStdin: true, UseStdout: true})
	if err != nil {
		t.Fatal(err)
	}
	b, err := ioutil.ReadFile(argsfile)
	if err != nil {
		t.Fatal(err)
	}
	args := strings.Split(string(b), "\n")
	if len(args) != 2 {
		t.Fatalf("expected 2 args, but got %q", args)
	}
	if args[0] != outputfile {
		t.Errorf("expected output path to be second arg, but was %q", args[1])
	}
	if args[1] != "templatepath" {
		t.Errorf("expected template path to be third arg, but was %q", args[2])
	}

	b, err = ioutil.ReadFile(datafile)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, expected) {
		t.Fatalf("expected data:\n%s\n\nactual data:\n%s", expected, b)
	}

	b, err = ioutil.ReadFile(outputfile)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != output {
		t.Errorf("expected to have written stdout to file %q, but got %q", output, b)
	}
}
