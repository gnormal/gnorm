{{- $short := (shortname .Name "err" "res" "sqlstr" "db" "XOLog") -}}
{{- $table := (schema .Schema .Table.TableName) -}}

// {{ .Name }}Table is the database name for the table.
const {{ .Name }}Table = "{{$table}}"

{{ if .Comment -}}
// {{ .Comment }}
{{- else -}}
// {{ .Name }} represents a row from '{{ $table }}'.
{{- end }}
type {{ .Name }} struct {
{{- range .Fields }}
	{{ .Name }} {{ retype .Type }} `json:"{{ .Col.ColumnName }}"` // {{ .Col.ColumnName }}
{{- end }}
}

{{$name := .Name}}

// Constants defining each column in the table.
const (
{{- range .Fields }}
	{{$name}}{{.Name}}Field = "{{ .Col.ColumnName }}"
{{- end -}}
)


{{/* For each field, we generate a where clause of the appropriate type. */}}

// WhereClauses for every type in {{.Name}}.
var (
{{- range .Fields }}
	{{$name}}{{.Name}}Where {{.Type | typeName | titleCase }}Field = "{{ .Col.ColumnName }}"
{{- end -}}
)

// QueryOne{{ .Name }} retrieves a row from '{{ $table }}' as a {{ .Name }}.
func QueryOne{{ .Name }}(db XODB, where WhereClause, order OrderBy) (*{{ .Name }}, error) {
	const origsqlstr = `SELECT ` +
		`{{ colnames .Fields }} ` +
		`FROM {{ $table }} WHERE (`

	idx := 1
	sqlstr := origsqlstr + where.String(&idx) + ") " + order.String() + " LIMIT 1"

	{{ $short }} := &{{ .Name }}{}
	err := db.QueryRow(sqlstr, where.Values()...).Scan({{ fieldnames .Fields (print "&" $short) }})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return {{ $short }}, nil
}

// Query{{ .Name }} retrieves rows from '{{ $table }}' as a slice of {{ .Name }}.
func Query{{ .Name }}(db XODB, where WhereClause, order OrderBy) ([]*{{ .Name }}, error) {
	const origsqlstr = `SELECT ` +
		`{{ colnames .Fields }} ` +
		`FROM {{ $table }} WHERE (`

	idx := 1
	sqlstr := origsqlstr + where.String(&idx) + ") " + order.String()

	var vals []*{{ .Name }}
	q, err := db.Query(sqlstr, where.Values()...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	for q.Next() {
		{{ $short }} := {{ .Name }}{}

		err = q.Scan({{ fieldnames .Fields (print "&" $short) }})
		if err != nil {
			return nil, errors.WithStack(err)
		}

		vals = append(vals, &{{ $short }})
	}
	return vals, nil
}

// All{{ .Name }} retrieves all rows from '{{ $table }}' as a slice of {{ .Name }}.
func All{{ .Name }}(db XODB, order OrderBy) ([]*{{ .Name }}, error) {
	const origsqlstr = `SELECT ` +
		`{{ colnames .Fields }} ` +
		`FROM {{ $table }}`

	sqlstr := origsqlstr + order.String()

	var vals []*{{ .Name }}
	q, err := db.Query(sqlstr)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	for q.Next() {
		{{ $short }} := {{ .Name }}{}

		err = q.Scan({{ fieldnames .Fields (print "&" $short) }})
		if err != nil {
			return nil, errors.WithStack(err)
		}

		vals = append(vals, &{{ $short }})
	}
	return vals, nil
}

