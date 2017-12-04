package run

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"text/template"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"gnorm.org/gnorm/environ"
	"gnorm.org/gnorm/run/data"
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
Table: {{.Name}}({{$schema}}.{{.DBName}}){{if ne .Comment ""}}; {{.Comment}}{{end}}
{{makeTable .Columns "{{.Name}}|{{.DBName}}|{{.Type}}|{{.DBType}}|{{.IsArray}}|{{.IsPrimaryKey}}|{{.IsFK}}|{{.HasFKRef}}|{{.Length}}|{{.UserDefined}}|{{.Nullable}}|{{.HasDefault}}|{{.Comment}}" "Name" "DBName" "Type" "DBType" "IsArray" "IsPrimaryKey" "IsFK" "HasFKRef" "Length" "UserDefined" "Nullable" "HasDefault" "Comment" -}}
Indexes:
{{makeTable .Indexes "{{.Name}}|{{.DBName}}|{{join .Columns.Names \", \"}}" "Name" "DBName" "Columns"}}
{{end -}}
{{end -}}
`))

// PreviewFormat defines the types of output that Preview can return.
type PreviewFormat int

const (
	// PreviewTabular shows the data in textual tables.
	PreviewTabular PreviewFormat = iota
	// PreviewYAML shows the data in YAML.
	PreviewYAML
	// PreviewJSON shows the data in JSON.
	PreviewJSON
	// PreviewTypes just prints out the column types used by the DB.
	PreviewTypes
)

// Preview displays the database info that would be passed to your template
// based on your configuration.
func Preview(env environ.Values, cfg *Config, format PreviewFormat) error {
	info, err := cfg.Driver.Parse(env.Log, cfg.ConnStr, cfg.Schemas, makeFilter(cfg.IncludeTables, cfg.ExcludeTables))
	if err != nil {
		return err
	}
	data, err := makeData(env.Log, info, cfg)
	if err != nil {
		return err
	}
	switch format {
	case PreviewTypes:
		displayTypes(env, data)
		return nil
	case PreviewYAML:
		b, err := yaml.Marshal(data)
		if err != nil {
			return errors.WithMessage(err, "couldn't convert data to yaml")
		}
		_, err = env.Stdout.Write(b)
		return err
	case PreviewJSON:
		b, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return errors.WithMessage(err, "couldn't convert data to json")
		}
		_, err = env.Stdout.Write(b)
		return err
	case PreviewTabular:
		return previewTpl.Execute(env.Stdout, data)
	default:
		return errors.Errorf("Unsupported format: %v", format)
	}
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
	t, err := template.New("table").
		Funcs(map[string]interface{}{"join": strings.Join}).
		Parse("{{range .}}" + templateStr + "\n{{end}}")
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

func displayTypes(env environ.Values, info *data.DBData) {
	var nullCols []*data.Column
	var cols []*data.Column
	lookUp := make(map[string]bool)
	nullLookUp := make(map[string]bool)
	for _, v := range info.Schemas {
		for _, t := range v.Tables {
			for _, c := range t.Columns {
				if c.Nullable {
					if !nullLookUp[c.DBType] {
						nullCols = append(nullCols, c)
						nullLookUp[c.DBType] = true
					}
				} else {
					if !lookUp[c.DBType] {
						cols = append(cols, c)
						lookUp[c.DBType] = true
					}
				}
			}
		}
	}
	sort.SliceStable(cols, func(i, j int) bool {
		return cols[i].DBType < cols[j].DBType
	})
	sort.SliceStable(nullCols, func(i, j int) bool {
		return nullCols[i].DBType < nullCols[j].DBType
	})

	fmt.Fprintln(env.Stdout, "[TypeMap]")
	for _, c := range cols {
		fmt.Fprintf(env.Stdout, "%q = %q\n", c.DBType, c.Type)
	}
	fmt.Fprintln(env.Stdout)

	fmt.Fprintln(env.Stdout, "[NullableTypeMap]")
	for _, c := range nullCols {
		fmt.Fprintf(env.Stdout, "%q = %q\n", c.DBType, c.Type)
	}
}
