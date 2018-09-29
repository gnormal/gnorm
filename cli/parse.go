package cli // import "gnorm.org/gnorm/cli"

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"

	"gnorm.org/gnorm/database"
	"gnorm.org/gnorm/database/drivers/mysql"
	"gnorm.org/gnorm/database/drivers/postgres"
	"gnorm.org/gnorm/environ"
	"gnorm.org/gnorm/run"
	"gnorm.org/gnorm/run/data"
)

func parseFile(env environ.Values, file string) (*run.Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, errors.WithMessage(err, "can't open config file")
	}
	defer f.Close()
	return Parse(env, f)
}

// Parse reads the configuration file and returns a gnorm config value.
func Parse(env environ.Values, r io.Reader) (*run.Config, error) {
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
	if c.OutputDir == "" {
		c.OutputDir = "."
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
		ConfigData: data.ConfigData{
			ConnStr:          c.ConnStr,
			DBType:           c.DBType,
			Schemas:          c.Schemas,
			NullableTypeMap:  c.NullableTypeMap,
			TypeMap:          c.TypeMap,
			PostRun:          c.PostRun,
			ExcludeTables:    exclude,
			IncludeTables:    include,
			OutputDir:        c.OutputDir,
			StaticDir:        c.StaticDir,
			PluginDirs:       c.PluginDirs,
			NoOverwriteGlobs: c.NoOverwriteGlobs,
		},
		Params: c.Params,
	}
	d, err := getDriver(strings.ToLower(c.DBType))
	if err != nil {
		return nil, err
	}
	cfg.Driver = d

	environ.FuncMap["plugin"] = environ.Plugin(c.PluginDirs)

	t, err := template.New("NameConversion").Funcs(environ.FuncMap).Parse(c.NameConversion)
	if err != nil {
		return nil, errors.WithMessage(err, "error parsing NameConversion template")
	}
	cfg.NameConversion = t

	if len(c.TemplateEngine.CommandLine) != 0 {
		for _, s := range c.TemplateEngine.CommandLine {
			t, err = template.New("EngineCLI").Funcs(environ.FuncMap).Parse(s)
			if err != nil {
				return nil, errors.WithMessage(err, "error parsing TemplateEngine CLI template")
			}
			cfg.TemplateEngine.CommandLine = append(cfg.TemplateEngine.CommandLine, t)
		}
		cfg.TemplateEngine.UseStdin = c.TemplateEngine.UseStdin
		cfg.TemplateEngine.UseStdout = c.TemplateEngine.UseStdout
	}

	useEngine := len(c.TemplateEngine.CommandLine) != 0
	cfg.SchemaPaths, err = parseOutputTargets(c.SchemaPaths, useEngine)
	if err != nil {
		return nil, errors.WithMessage(err, "error parsing SchemaPaths")
	}

	cfg.TablePaths, err = parseOutputTargets(c.TablePaths, useEngine)
	if err != nil {
		return nil, errors.WithMessage(err, "error parsing TablePaths")
	}

	cfg.EnumPaths, err = parseOutputTargets(c.EnumPaths, useEngine)
	if err != nil {
		return nil, errors.WithMessage(err, "error parsing EnumPaths")
	}

	if len(cfg.EnumPaths) == 0 && len(cfg.TablePaths) == 0 && len(cfg.SchemaPaths) == 0 {
		return nil, errors.New("no output paths defined, so no output will be generated")
	}

	cfg.ConnStr = os.Expand(c.ConnStr, func(s string) string {
		return env.Env[s]
	})
	return cfg, nil
}

func getDriver(name string) (database.Driver, error) {
	switch name {
	case "postgres":
		return postgres.PG{}, nil
	case "mysql":
		return mysql.MySQL{}, nil
	default:
		return nil, errors.Errorf("unknown database type: %v", name)
	}
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

func parseOutputTargets(vals map[string]string, usePath bool) ([]run.OutputTarget, error) {
	out := make([]run.OutputTarget, 0, len(vals))
	for fnTempl, contTempl := range vals {
		fn, err := template.New("filename").Funcs(environ.FuncMap).Parse(fnTempl)
		if err != nil {
			return nil, errors.WithMessage(err, "error parsing filename template")
		}
		if usePath {
			// use path means we're using an external template engine, so don't try to
			// parse the template.
			if _, err := os.Stat(contTempl); err != nil {
				return nil, errors.WithMessage(err, "error checking contents template")
			}
			out = append(out, run.OutputTarget{Filename: fn, ContentsPath: contTempl})
			continue
		}
		b, err := ioutil.ReadFile(contTempl)
		if err != nil {
			return nil, errors.WithMessage(err, "error reading contents template")
		}
		cont, err := template.New(contTempl).Funcs(environ.FuncMap).Parse(string(b))
		if err != nil {
			return nil, errors.WithMessage(err, "error parsing contents template")
		}
		out = append(out, run.OutputTarget{Filename: fn, Contents: cont})
	}
	return out, nil
}
