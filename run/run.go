package run // import "gnorm.org/gnorm/run"
import (
	"bytes"
	"os"
	"strings"
	"text/tabwriter"
	"text/template"

	"github.com/pkg/errors"

	yaml "gopkg.in/yaml.v2"

	"gnorm.org/gnorm/database"
	"gnorm.org/gnorm/environ"
)

// Config holds the schema that is expected to exist in the gnorm.toml file.
type Config struct {
	// ConnStr is the connection string for the database.
	ConnStr string

	// Schemas holds the names of schemas to generate code for.
	Schemas []string

	// TemplateDir contains the relative path to the directory where gnorm
	// expects to find templates to render.  The default is the current
	// directory where gnorm is running.
	TemplateDir string

	// TablePath is a relative path for tables to be rendered.  The table
	// template will be rendered with each table in turn. If the path is empty,
	// tables will not be rendered.
	//
	// The table path may be a template, in which case the values .Schema and
	// .Table may be referenced, containing the name of the current schema and
	// table being rendered.  For example, "{{.Schema}}/{{.Table}}/{{.Table}}.go" would render
	// the "public.users" table to ./public/users/users.go.
	TablePath string

	// SchemaPath is a relative path for schemas to be rendered.  The schema
	// template will be rendered with each schema in turn. If the path is empty,
	// schema will not be rendered.
	//
	// The schema path may be a template, in which case the value .Schema may be
	// referenced, containing the name of the current schema being rendered. For
	// example, "schemas/{{.Schema}}/{{.Schema}}.go" would render the "public"
	// schema to ./schemas/public/public.go
	SchemaPath string

	TypeMap map[string]string

	NullableTypeMap map[string]string
}

// Preview displays the database info that woudl be passed to your template
// based on your configuration.
func Preview(env environ.Values, cfg Config, verbose bool) error {
	info, err := dbInfo(env, cfg)
	if err != nil {
		return err
	}
	if verbose {
		b, err := yaml.Marshal(info)
		if err != nil {
			return errors.WithMessage(err, "couldn't convert data to yaml")
		}
		_, err = env.Stdout.Write(b)
		return err
	}
	t := template.Must(template.New("summary").Funcs(map[string]interface{}{
		"makeTable": makeTable,
	}).Parse(`
{{- range .Schemas -}}
{{.Name}}
  {{range .Tables -}}
  {{.Name}}
{{makeTable .Columns "    " "{{.Name}}\t{{.Type}}" "ColName" "Type"}}
  {{end}}
{{end}}
`))
	return t.Execute(env.Stdout, info)
}

func dbInfo(env environ.Values, cfg Config) (*database.SchemaInfo, error) {
	expand := func(s string) string {
		return env.Env[s]
	}
	conn := os.Expand(cfg.ConnStr, expand)
	return database.Parse(cfg.TypeMap, cfg.NullableTypeMap, env, conn, cfg.Schemas)

}

func makeTable(data interface{}, prefix, templateStr string, columnTitles ...string) (string, error) {
	t, err := template.New("table").Parse("{{range .}}" + prefix + templateStr + "\n{{end}}")
	if err != nil {
		return "", errors.WithMessage(err, "failed to parse table template")
	}
	buf := &bytes.Buffer{}
	w := tabwriter.NewWriter(buf, 0, 4, 1, byte(' '), 0)
	if len(columnTitles) > 0 {
		columnTitles[0] = prefix + columnTitles[0]
	}
	_, err = w.Write([]byte(strings.Join(columnTitles, "\t") + "\n"))
	if err != nil {
		return "", err
	}
	if err := t.Execute(w, data); err != nil {
		return "", errors.WithMessage(err, "failed to run table template")
	}
	if err := w.Flush(); err != nil {
		return "", err
	}
	return buf.String(), nil
}
