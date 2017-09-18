package run

import (
	"bytes"
	"log"
	"testing"
	"text/template"

	"gnorm.org/gnorm/database"
	"gnorm.org/gnorm/environ"
	"gnorm.org/gnorm/run/data"
)

func TestMakeData(t *testing.T) {
	c := &Config{
		NameConversion: template.Must(template.New("").Funcs(environ.FuncMap).Parse(`{{print "abc " .}}`)),
		ConfigData: data.ConfigData{
			NullableTypeMap: map[string]string{
				"*int": "*INTEGER",
			},
			TypeMap: map[string]string{
				"int": "INTEGER",
			},
		},
	}

	info := &database.Info{
		Schemas: []*database.Schema{{
			Name: "schema",
			Tables: []*database.Table{{
				Name: "table",
				Columns: []*database.Column{{
					Name: "col1",
					Type: "int",
				}, {
					Name:     "col2",
					Type:     "*int",
					Nullable: true,
				}, {
					Name: "col3",
					Type: "string",
				}, {
					Name:     "col4",
					Type:     "*string",
					Nullable: true,
				}},
			}},
			Enums: []*database.Enum{{
				Name: "enum",
				Values: []*database.EnumValue{{
					Name: "enumvalue",
				}},
			}},
		}},
	}

	buf := &bytes.Buffer{}
	l := log.New(buf, "", 0)
	data, err := makeData(l, info, c)
	if err != nil {
		t.Fatal("unexpected error from convertNames", err)
	}
	expected := "abc " + info.Schemas[0].Name
	got := data.Schemas[0].Name
	if got != expected {
		t.Errorf("schema name expected %q but got %q", expected, got)
	}

	expected = "abc " + info.Schemas[0].Tables[0].Name
	got = data.Schemas[0].Tables[0].Name
	if got != expected {
		t.Errorf("table name expected %q but got %q", expected, got)
	}

	expected = "abc " + info.Schemas[0].Tables[0].Columns[0].Name
	got = data.Schemas[0].Tables[0].Columns[0].Name
	if got != expected {
		t.Errorf("column0 name expected %q but got %q", expected, got)
	}

	expected = c.TypeMap[info.Schemas[0].Tables[0].Columns[0].Type]
	got = data.Schemas[0].Tables[0].Columns[0].Type
	if got != expected {
		t.Errorf("column0 type expected %q but got %q", expected, got)
	}

	expected = "abc " + info.Schemas[0].Tables[0].Columns[1].Name
	got = data.Schemas[0].Tables[0].Columns[1].Name
	if got != expected {
		t.Errorf("column1 name expected %q but got %q", expected, got)
	}

	expected = c.NullableTypeMap[info.Schemas[0].Tables[0].Columns[1].Type]
	got = data.Schemas[0].Tables[0].Columns[1].Type
	if got != expected {
		t.Errorf("column1 type expected %q but got %q", expected, got)
	}

	got = data.Schemas[0].Tables[0].Columns[2].Type
	if got != "" {
		t.Errorf("column2 type expected to be empty but got %q", got)
	}

	got = data.Schemas[0].Tables[0].Columns[3].Type
	if got != "" {
		t.Errorf("column3 type expected to be empty but got %q", got)
	}

	expected = "abc " + info.Schemas[0].Enums[0].Name
	got = data.Schemas[0].Enums[0].Name
	if got != expected {
		t.Errorf("enum name expected %q but got %q", expected, got)
	}

	expected = "abc " + info.Schemas[0].Enums[0].Values[0].Name
	got = data.Schemas[0].Enums[0].Values[0].Name
	if got != expected {
		t.Errorf("enum value name expected %q but got %q", expected, got)
	}
}
