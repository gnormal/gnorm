package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"gnorm.org/gnorm/environ"
	"gnorm.org/gnorm/run"
)

var (
	version    = "DEV"
	timestamp  = "no timestamp, did you build with make.go?"
	commitHash = "no hash, did you build with make.go?"
)

func previewCmd(env environ.Values) *cobra.Command {
	var cfgFile string
	var verbose bool
	var format string
	preview := &cobra.Command{
		Use:   "preview",
		Short: "Preview the data that will be sent to your templates",
		Long: `
Reads your gnorm.toml file and connects to your database, translating the schema
just as it would be during a full run.  It is then printed out in an
easy-to-read format.  By default it prints out the data in a human-readable
plaintext tabular format.  You may specify a different format using the -format
flag, in which case you can print json, yaml, or types, where types is a list of
all types used by columns in your database.  The latter is useful when setting
up TypeMaps.
`[1:],
		RunE: func(cmd *cobra.Command, args []string) error {
			env.InitLog(verbose)
			pformat := run.PreviewTabular
			switch strings.ToLower(format) {
			case "tabular":
				pformat = run.PreviewTabular
			case "yaml":
				pformat = run.PreviewYAML
			case "json":
				pformat = run.PreviewJSON
			case "types":
				pformat = run.PreviewTypes
			default:
				return codeErr{errors.Errorf("unknown preview format %q", format), 2}
			}
			cfg, err := parseFile(env, cfgFile)
			if err != nil {
				return codeErr{err, 2}
			}
			if err := run.Preview(env, cfg, pformat); err != nil {
				return codeErr{err, 1}
			}
			return nil
		},
		Args: cobra.ExactArgs(0),
	}
	preview.Flags().StringVarP(&cfgFile, "config", "c", "gnorm.toml", "relative path to gnorm config file")
	preview.Flags().StringVarP(&format, "format", "f", "tabular", "Specify output format: tabular, yaml, json, or types")
	preview.Flags().BoolVarP(&verbose, "verbose", "v", false, "show debugging output")
	return preview
}

func genCmd(env environ.Values) *cobra.Command {
	var cfgFile string
	var verbose bool
	gen := &cobra.Command{
		Use:   "gen",
		Short: "Generate code from DB schema",
		Long: `
Reads your gnorm.toml file and connects to your database, translating the schema
into in-memory objects.  Then reads your templates and writes files to disk
based on those templates.`[1:],
		RunE: func(cmd *cobra.Command, args []string) error {
			env.InitLog(verbose)
			cfg, err := parseFile(env, cfgFile)
			if err != nil {
				return codeErr{err, 2}
			}
			if err := run.Generate(env, cfg); err != nil {
				return codeErr{err, 1}
			}
			return nil
		},
		Args: cobra.ExactArgs(0),
	}
	gen.Flags().StringVarP(&cfgFile, "config", "c", "gnorm.toml", "relative path to gnorm config file")
	gen.Flags().BoolVarP(&verbose, "verbose", "v", false, "show debugging output")
	return gen
}

func versionCmd(env environ.Values) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Displays the version of GNORM.",
		Long: `
Shows the build date and commit hash used to build this binary.`[1:],
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(env.Stdout, "version: %s\nbuilt at: %s\ncommit hash: %s", version, timestamp, commitHash)
		},
	}
}

func initCmd(env environ.Values) *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Generates the files needed to run GNORM.",
		Long: `
Creates a default gnorm.toml and the various template files needed to run GNORM.`[1:],
		RunE: func(cmd *cobra.Command, args []string) error {
			return initFunc(".")
		},
		Args: cobra.ExactArgs(0),
	}
}

func docCmd(env environ.Values) *cobra.Command {
	return &cobra.Command{
		Use:   "docs",
		Short: "Runs a local webserver serving gnorm documentation.",
		Long: `
Starts a web server running at localhost:8080 that serves docs for this version
of Gnorm.`[1:],
		RunE: func(cmd *cobra.Command, args []string) error {
			return showDocs(env, cmd, args)
		},
	}
}

func initFunc(dir string) error {
	if err := createFile(filepath.Join(dir, "gnorm.toml"), sample); err != nil {
		return err
	}
	if err := createFile(filepath.Join(dir, "templates/table.gotmpl"), "Table: {{.Table.Name}}\n{{printf \"%#v\" .}}"); err != nil {
		return err
	}
	if err := createFile(filepath.Join(dir, "templates/enum.gotmpl"), "Enum: {{.Enum.Name}}\n{{printf \"%#v\" .}}"); err != nil {
		return err
	}
	return createFile(filepath.Join(dir, "templates/schema.gotmpl"), "Schema: {{.Schema.Name}}\n{{printf \"%#v\" .}}")
}

func createFile(name, contents string) error {
	if dir := filepath.Dir(name); dir != "" {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return codeErr{
				errors.WithMessage(err, fmt.Sprintf("Can't create directory %q", dir)),
				1,
			}
		}
	}

	f, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return codeErr{
			errors.WithMessage(err, fmt.Sprintf("Can't create file %q", name)),
			1,
		}
	}
	defer f.Close()
	_, err = f.WriteString(contents)
	if err == nil {
		return nil
	}
	return codeErr{
		errors.WithMessage(err, fmt.Sprintf("Failed writing data to file %q", name)),
		1,
	}
}

type codeErr struct {
	error
	code int
}

func (c codeErr) Code() int {
	return c.code
}
