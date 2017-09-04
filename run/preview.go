package run

import (
	"bytes"
	"encoding/csv"
	"strings"
	"text/template"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"gnorm.org/gnorm/environ"
	yaml "gopkg.in/yaml.v2"
)

var previewTpl = template.Must(
	template.New("preview").
		Funcs(map[string]interface{}{"makeTable": makeTable}).
		Parse(`
{{- range .Schemas }}{{$schema := .DBName -}}
Schema: {{.Name}}({{.DBName}})
{{range .Enums}}
Enum: {{.Name}}({{$schema}}.{{.DBName}})
{{makeTable .Values "{{.Name}}|{{.DBName}}|{{.Value}}" "Name" "DBName" "Value" }}
{{end -}}
{{range .Tables}}
Table: {{.Name}}({{$schema}}.{{.DBName}})
{{makeTable .Columns "{{.Name}}|{{.DBName}}|{{.Type}}|{{.DBType}}|{{.IsArray}}|{{.Length}}|{{.UserDefined}}|{{.Nullable}}|{{.HasDefault}}" "Name" "DBName" "Type" "DBType" "IsArray" "Length" "UserDefined" "Nullable" "HasDefault"}}
{{end -}}
{{end -}}
`))

// Preview displays the database info that woudl be passed to your template
// based on your configuration.
func Preview(env environ.Values, cfg *Config, useYaml bool) error {
	info, err := getDBInfo(env, cfg)
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
	return previewTpl.Execute(env.Stdout, info)
}

// makeTable makes a nice-looking textual table from the given data using the
// given template as the rendering for each line.  Columns in the template
// should be separated by a pipe '|'.  Column titles are prepended to the table
// if they exist.
//
// For example where here people is a slice of structs with a Name and Age fields:
//    ```
//    makeTable(people, "{{.Name}}|{{.Age}}", "Name", "Age")
//
//    +----------+-----+
//    |   Name   | Age |
//    +----------+-----+
//    | Bob      |  30 |
//    | Samantha |   3 |
//    +----------+-----+
//    ```
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
	table.SetAutoFormatHeaders(false)
	table.Render()
	return output.String(), nil
}
