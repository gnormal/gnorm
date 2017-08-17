{{- $type := .Name -}}
{{- $short := (shortname $type "enumVal" "text" "buf" "ok" "src") -}}
{{- $reverseNames := .ReverseConstNames -}}
// {{ $type }} is the '{{ .Enum.EnumName }}' enum type from schema '{{ .Schema  }}'.
type {{ $type }} uint16

const (
    // Unknown{{$type}} defines an invalid {{$type}}. 
    Unknown{{$type}} = {{$type}}(0)
{{ range .Values }}
	// {{ if $reverseNames }}{{ .Name }}{{ $type }}{{ else }}{{ $type }}{{ .Name }}{{ end }} is the '{{ .Val.EnumValue }}' {{ $type }}.
	{{ if $reverseNames }}{{ .Name }}{{ $type }}{{ else }}{{ $type }}{{ .Name }}{{ end }} = {{ $type }}({{ .Val.ConstValue }})
{{ end -}}
)

// String returns the string value of the {{ $type }}.
func ({{ $short }} {{ $type }}) String() string {
	switch {{ $short }} {
{{- range .Values }}
	case {{ if $reverseNames }}{{ .Name }}{{ $type }}{{ else }}{{ $type }}{{ .Name }}{{ end }}:
		return "{{ .Val.EnumValue }}"
{{- end }}
	default:
    	return "Unknown{{$type}}"
    }
}

// MarshalText marshals {{ $type }} into text.
func ({{ $short }} {{ $type }}) MarshalText() ([]byte, error) {
	return []byte({{ $short }}.String()), nil
}

// UnmarshalText unmarshals {{ $type }} from text.
func ({{ $short }} *{{ $type }}) UnmarshalText(text []byte) error {
    val, err := Parse{{$type}}(string(text))
    if err != nil {
        return err
    }
    *{{$short}} = val
    return nil
}

// Parse{{$type}} converts s into a {{$type}} if it is a valid
// stringified value of {{$type}}.
func Parse{{$type}}(s string) ({{$type}}, error) {
	switch s {
{{- range .Values }}
	case "{{ .Val.EnumValue }}":
		return {{ if $reverseNames }}{{ .Name }}{{ $type }}{{ else }}{{ $type }}{{ .Name }}{{ end }}, nil
{{- end }}
	default:
		return Unknown{{$type}}, errors.New("invalid {{ $type }}")
	}
}

// Value satisfies the sql/driver.Valuer interface for {{ $type }}.
func ({{ $short }} {{ $type }}) Value() (driver.Value, error) {
	return {{ $short }}.String(), nil
}

// Scan satisfies the database/sql.Scanner interface for {{ $type }}.
func ({{ $short }} *{{ $type }}) Scan(src interface{}) error {
	buf, ok := src.([]byte)
	if !ok {
	   return errors.New("invalid {{ $type }}")
	}

	return {{ $short }}.UnmarshalText(buf)
}

// {{$type}}Field is a component that returns a WhereClause that contains a
// comparison based on its field and a strongly typed value.
type {{$type}}Field string

// Equals returns a WhereClause for this field.
func (f {{$type}}Field) Equals(v {{$type}}) WhereClause {
	return whereClause{
		field: string(f),
		comp:  compEqual,
		value: v,
	}
}

// GreaterThan returns a WhereClause for this field.
func (f {{$type}}Field) GreaterThan(v {{$type}}) WhereClause {
	return whereClause{
		field: string(f),
		comp:  compGreater,
		value: v,
	}
}

// LessThan returns a WhereClause for this field.
func (f {{$type}}Field) LessThan(v {{$type}}) WhereClause {
	return whereClause{
		field: string(f),
		comp:  compEqual,
		value: v,
	}
}

// GreaterOrEqual returns a WhereClause for this field.
func (f {{$type}}Field) GreaterOrEqual(v {{$type}}) WhereClause {
	return whereClause{
		field: string(f),
		comp:  compGTE,
		value: v,
	}
}

// LessOrEqual returns a WhereClause for this field.
func (f {{$type}}Field) LessOrEqual(v {{$type}}) WhereClause {
	return whereClause{
		field: string(f),
		comp:  compLTE,
		value: v,
	}
}

// NotEqual returns a WhereClause for this field.
func (f {{$type}}Field) NotEqual(v {{$type}}) WhereClause {
	return whereClause{
		field: string(f),
		comp:  compNE,
		value: v,
	}
}

// In returns a WhereClause for this field.
func (f {{$type}}Field) In(vals []{{$type}}) WhereClause {
	values := make([]interface{}, len(vals))
	for x := range vals {
		values[x] = vals[x]
	}
	return inClause{
		field: string(f),
		value: values,
	}
}