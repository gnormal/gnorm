package cli // import "gnorm.org/gnorm/cli"

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/BurntSushi/toml"
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

// ParseAndRun parses the environment and runs the command.
func ParseAndRun(env environ.Values) int {
	// the error from the executed command
	var code int
	rootCmd := &cobra.Command{
		Use:   "gnorm",
		Short: "GNORM is Not an ORM, it's a db schema->code generator",
		Long: `
A flexible code generator that turns your DB schema into
runnable code.  See full docs at https://gnorm.org`[1:],
	}
	rootCmd.SetArgs(env.Args)
	rootCmd.SetOutput(env.Stderr)

	rootCmd.AddCommand(previewCmd(env, &code))
	rootCmd.AddCommand(genCmd(env, &code))
	rootCmd.AddCommand(versionCmd(env))
	if err := rootCmd.Execute(); err != nil {
		// cobra outputs the error itself.
		return 2
	}
	return code
}

func previewCmd(env environ.Values, code *int) *cobra.Command {
	var cfgFile string
	var useYaml bool
	var verbose bool
	preview := &cobra.Command{
		Use:   "preview",
		Short: "Preview the data that will be sent to your templates",
		Long: `
		Reads your gnorm.toml file and connects to your database, translating the schema
		just as it would be during a full run.  It is then printed out in an
		easy-to-read format.`[1:],
		Run: func(cmd *cobra.Command, args []string) {
			env.InitLog(verbose)
			cfg, err := parseFile(env, cfgFile)
			if err != nil {
				fmt.Fprintln(env.Stderr, err)
				*code = 2
				return
			}
			if err := run.Preview(env, cfg, useYaml); err != nil {
				fmt.Fprintln(env.Stderr, err)
				*code = 1
			}
		},
	}
	preview.Flags().StringVar(&cfgFile, "config", "gnorm.toml", "relative path to gnorm config file")
	preview.Flags().BoolVar(&useYaml, "yaml", false, "show output in yaml instead of tabular")
	preview.Flags().BoolVar(&verbose, "verbose", false, "show debugging output")
	return preview
}

func genCmd(env environ.Values, code *int) *cobra.Command {
	var cfgFile string
	var verbose bool
	gen := &cobra.Command{
		Use:   "gen",
		Short: "Generate code from DB schema",
		Long: `
Reads your gnorm.toml file and connects to your database, translating the schema
into in-memory objects.  Then reads your templates and writes files to disk
based on those templates.`[1:],
		Run: func(cmd *cobra.Command, args []string) {
			env.InitLog(verbose)
			cfg, err := parseFile(env, cfgFile)
			if err != nil {
				fmt.Fprintln(env.Stderr, err)
				*code = 2
				return
			}
			if err := run.Generate(env, cfg); err != nil {
				fmt.Fprintln(env.Stderr, err)
				*code = 1
			}
		},
	}
	gen.Flags().StringVar(&cfgFile, "config", "gnorm.toml", "relative path to gnorm config file")
	gen.Flags().BoolVar(&verbose, "verbose", false, "show debugging output")
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
func parseFile(env environ.Values, file string) (*run.Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, errors.WithMessage(err, "can't open config file")
	}
	defer f.Close()
	return parse(env, f)
}

// parse reads the configuration file and returns a gnorm config value.
func parse(env environ.Values, r io.Reader) (*run.Config, error) {
	c := Config{}
	m, err := toml.DecodeReader(r, &c)
	if err != nil {
		return nil, errors.WithMessage(err, "error parsing config file")
	}
	undec := m.Undecoded()
	if len(undec) > 0 {
		log.Println("Warning: unknown values present in config file:", undec)
	}

	if len(c.Schemas) == 0 {
		return nil, errors.New("no schemas specified in config")
	}

	if c.NameConversion == "" {
		return nil, errors.New("no NameConversion specified in config")
	}
	if len(c.ExcludeTables) > 0 && len(c.IncludeTables) > 0 {
		return nil, errors.New("both include tables and exclude tables")
	}

	include, err := parseTables(c.IncludeTables, c.Schemas)
	if err != nil {
		return nil, err
	}

	exclude, err := parseTables(c.ExcludeTables, c.Schemas)
	if err != nil {
		return nil, err
	}

	cfg := &run.Config{
		ConnStr:         c.ConnStr,
		Schemas:         c.Schemas,
		NullableTypeMap: c.NullableTypeMap,
		TypeMap:         c.TypeMap,
		TemplateDir:     c.TemplateDir,
		PostRun:         c.PostRun,
		ExcludeTables:   exclude,
		IncludeTables:   include,
	}

	switch strings.ToLower(c.DBType) {
	case "":
		return nil, errors.New("no DBType specificed")
	case "postgres":
		cfg.DBType = run.Postgres
	case "mysql":
		cfg.DBType = run.Mysql
	default:
		return nil, errors.Errorf("unsupported dbtype %q", c.DBType)
	}

	t, err := template.New("NameConversion").Funcs(environ.FuncMap).Parse(c.NameConversion)
	if err != nil {
		return nil, errors.WithMessage(err, "error parsing NameConversion template")
	}
	cfg.NameConversion = t

	if c.SchemaPath != "" {
		t, err := template.New("SchemaPath").Funcs(environ.FuncMap).Parse(c.SchemaPath)
		if err != nil {
			return nil, errors.WithMessage(err, "error parsing SchemaPath template")
		}
		cfg.SchemaPath = t
	}

	if c.TablePath != "" {
		t, err := template.New("TablePath").Funcs(environ.FuncMap).Parse(c.TablePath)
		if err != nil {
			return nil, errors.WithMessage(err, "error parsing SchemaPath template")
		}
		cfg.TablePath = t
	}

	if c.EnumPath != "" {
		t, err := template.New("EnumPath").Funcs(environ.FuncMap).Parse(c.EnumPath)
		if err != nil {
			return nil, errors.WithMessage(err, "error parsing SchemaPath template")
		}
		cfg.EnumPath = t
	}

	expand := func(s string) string {
		return env.Env[s]
	}

	cfg.ConnStr = os.Expand(c.ConnStr, expand)
	return cfg, nil
}

// parseTables takes a list of tablenames in "<schema.>table" format and spits
// out a map of schema to list of tables.  Tables with no schema apply to all
// schemas.  Tables with a schema apply to only that schema.  Tables that
// specify a schema not in the list of schemas given are an error.
func parseTables(tables, schemas []string) (map[string][]string, error) {
	out := make(map[string][]string, len(schemas))
	for _, s := range schemas {
		out[s] = nil
	}
	for _, t := range tables {
		vals := strings.Split(t, ".")
		switch len(vals) {
		case 1:
			// just the table name, so it goes for all schemas
			for schema := range out {
				out[schema] = append(out[schema], t)
			}
		case 2:
			// schema and table
			list, ok := out[vals[0]]
			if !ok {
				return nil, errors.Errorf("%q specified for tables but schema %q not in schema list", t, vals[0])
			}
			out[vals[0]] = append(list, vals[1])
		default:
			// too many periods... bad format
			return nil, errors.Errorf(`badly formatted table: %q, should be just "table" or "table.schema"`, t)
		}
	}

	return out, nil
}
