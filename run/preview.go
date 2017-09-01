package run

import (
	"text/template"

	"github.com/pkg/errors"
	"gnorm.org/gnorm/database/drivers/postgres"
	"gnorm.org/gnorm/environ"
	yaml "gopkg.in/yaml.v2"
)

// Preview displays the database info that woudl be passed to your template
// based on your configuration.
func Preview(env environ.Values, cfg *Config, useYaml bool) error {
	info, err := postgres.Parse(env.Log, cfg.ConnStr, cfg.Schemas)
	if err != nil {
		return err
	}
	if err := convertNames(env.Log, info, cfg); err != nil {
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
	t := template.Must(template.New("summary").Funcs(environ.FuncMap).Parse(`
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
	return t.Execute(env.Stdout, info)
}
