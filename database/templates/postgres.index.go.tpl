{{- $short := (shortname .Type.Name "err" "sqlstr" "db" "q" "res" "XOLog" .Fields) -}}
{{- $table := (schema .Schema .Type.Table.TableName) -}}
// {{ .FuncName }} retrieves a row from '{{ $table }}' as a {{ .Type.Name }}.
//
// Generated from index '{{ .Index.IndexName }}'.
func {{ .FuncName }}(db XODB{{ goparamlist .Fields true true }}) ({{ if not .Index.IsUnique }}[]{{ end }}*{{ .Type.Name }}, error) {
	const sqlstr = `SELECT ` +
		`{{ colnames .Type.Fields }} ` +
		`FROM {{ $table }} ` +
		`WHERE {{ colnamesquery .Fields " AND " }}`

{{- if .Index.IsUnique }}
	{{ $short }} := {{ .Type.Name }}{}

	err := db.QueryRow(sqlstr{{ goparamlist .Fields true false }}).Scan({{ fieldnames .Type.Fields (print "&" $short) }})
	if err != nil {
		return nil, err
	}

	return &{{ $short }}, nil
{{- else }}
	q, err := db.Query(sqlstr{{ goparamlist .Fields true false }})
	if err != nil {
		return nil, err
	}
	defer q.Close()

	res := []*{{ .Type.Name }}{}
	for q.Next() {
		{{ $short }} := {{ .Type.Name }}{}

		err := q.Scan({{ fieldnames .Type.Fields (print "&" $short) }})
		if err != nil {
			return nil, err
		}

		res = append(res, &{{ $short }})
	}

	return res, nil
{{- end }}
}

