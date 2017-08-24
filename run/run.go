package run // import "gnorm.org/gnorm/run"
import (
	"bytes"
	"encoding/csv"
	"os"
	"strings"
	"text/template"

	"github.com/olekukonko/tablewriter"
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
func Preview(env environ.Values, cfg Config, useYaml, verbose bool) error {
	info, err := dbInfo(env, cfg)
	if err != nil {
		return err
	}
	if useYaml {
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
{{- range .Schemas }}{{$schema := .Name -}}
Schema: {{.Name}}
{{range .Tables}}
Table: {{$schema}}.{{.Name}}
{{makeTable .Columns "{{.Name}}|{{.Type}}" "Column" "Type"}}
{{end -}}
{{end -}}
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

// makeTable makes a nice-looking textual table from the given data using the
// given template as the rendering for each line.  Columns in the template
// should be separated by a pipe '|'.  Column titles are prepended to the table
// if they exist.
func makeTable(data interface{}, templateStr string, columnTitles ...string) (string, error) {
	t, err := template.New("table").Parse("{{range .}}" + templateStr + "\n{{end}}")
	if err != nil {
		return "", errors.WithMessage(err, "failed to parse table template")
	}
	buf := &bytes.Buffer{}
	hasHeader := len(columnTitles) > 0
	if hasHeader {
		// this can't fail so we drop the error
		_, _ = buf.WriteString(strings.Join(columnTitles, "|") + "\n")
	}
	if err := t.Execute(buf, data); err != nil {
		return "", errors.WithMessage(err, "failed to run table template")
	}
	r := csv.NewReader(buf)
	r.Comma = '|'
	output := &bytes.Buffer{}
	table, err := tablewriter.NewCSVReader(output, r, hasHeader)
	if err != nil {
		return "", errors.WithMessage(err, "failed to render from pipe delimited bytes")
	}
	table.Render()
	return output.String(), nil
}
