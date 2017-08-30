package run

import (
	"bytes"
	"encoding/csv"
	"os"
	"strings"
	"text/template"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"gnorm.org/gnorm/database"
	"gnorm.org/gnorm/database/drivers/postgres"
	"gnorm.org/gnorm/environ"
	yaml "gopkg.in/yaml.v2"
)

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
{{range .Enums}}
Enum: {{$schema}}.{{.Name}}
{{makeTable .Values "{{.Name}}|{{.Value}}" "Name" "Value" }}
{{end -}}
{{range .Tables}}
Table: {{$schema}}.{{.Name}}
{{makeTable .Columns "{{.Name}}|{{.Type}}|{{.DBType}}|{{.IsArray}}|{{.Length}}|{{.UserDefined}}|{{.Nullable}}|{{.HasDefault}}" "Column" "Type" "DBType" "IsArray" "Length" "UserDefined" "Nullable" "HasDefault"}}
{{end -}}
{{end -}}
`))
	return t.Execute(env.Stdout, info)
}

func dbInfo(env environ.Values, cfg Config) (*database.Info, error) {
	expand := func(s string) string {
		return env.Env[s]
	}
	conn := os.Expand(cfg.ConnStr, expand)
	return postgres.Parse(cfg.TypeMap, cfg.NullableTypeMap, env.Log, conn, cfg.Schemas)

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
