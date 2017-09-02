package run

import (
	"bytes"
	"log"
	"testing"
	"text/template"

	"gnorm.org/gnorm/database"
	"gnorm.org/gnorm/environ"
)

func TestConvertNames(t *testing.T) {
	c := &Config{
		NameConversion: template.Must(template.New("").Funcs(environ.FuncMap).Parse(`{{print "abc " .}}`)),
		NullableTypeMap: map[string]string{
			"*int": "*INTEGER",
		},
		TypeMap: map[string]string{
			"int": "INTEGER",
		},
	}

	info := &database.Info{
		Schemas: []*database.Schema{{
			DBName: "schema",
			Tables: []*database.Table{{
				DBSchema: "schema",
				DBName:   "table",
				Columns: []*database.Column{{
					DBName: "col1",
					DBType: "int",
				}, {
					DBName:   "col2",
					DBType:   "*int",
					Nullable: true,
				}, {
					DBName: "col3",
					DBType: "string",
				}, {
					DBName:   "col4",
					DBType:   "*string",
					Nullable: true,
				}},
			}},
			Enums: []*database.Enum{{
				DBSchema: "schema",
				DBName:   "enum",
				Values: []*database.EnumValue{{
					DBName: "enumvalue",
				}},
			}},
		}},
	}

	buf := &bytes.Buffer{}
	l := log.New(buf, "", 0)
	err := convertNames(l, info, c)
	if err != nil {
		t.Fatal("unexpected error from convertNames", err)
	}
	expected := "abc " + info.Schemas[0].DBName
	got := info.Schemas[0].Name
	if got != expected {
		t.Errorf("schema name expected %q but got %q", expected, got)
	}

	expected = "abc " + info.Schemas[0].Tables[0].DBName
	got = info.Schemas[0].Tables[0].Name
	if got != expected {
		t.Errorf("table name expected %q but got %q", expected, got)
	}

	expected = info.Schemas[0].Name
	got = info.Schemas[0].Tables[0].Schema
	if got != expected {
		t.Errorf("table schema expected %q but got %q", expected, got)
	}

	expected = info.Schemas[0].DBName
	got = info.Schemas[0].Tables[0].DBSchema
	if got != expected {
		t.Errorf("table db schema expected %q but got %q", expected, got)
	}

	expected = "abc " + info.Schemas[0].Tables[0].Columns[0].DBName
	got = info.Schemas[0].Tables[0].Columns[0].Name
	if got != expected {
		t.Errorf("column0 name expected %q but got %q", expected, got)
	}

	expected = c.TypeMap[info.Schemas[0].Tables[0].Columns[0].DBType]
	got = info.Schemas[0].Tables[0].Columns[0].Type
	if got != expected {
		t.Errorf("column0 type expected %q but got %q", expected, got)
	}

	expected = "abc " + info.Schemas[0].Tables[0].Columns[1].DBName
	got = info.Schemas[0].Tables[0].Columns[1].Name
	if got != expected {
		t.Errorf("column1 name expected %q but got %q", expected, got)
	}

	expected = c.NullableTypeMap[info.Schemas[0].Tables[0].Columns[1].DBType]
	got = info.Schemas[0].Tables[0].Columns[1].Type
	if got != expected {
		t.Errorf("column1 type expected %q but got %q", expected, got)
	}

	got = info.Schemas[0].Tables[0].Columns[2].Type
	if got != "" {
		t.Errorf("column2 type expected to be empty but got %q", got)
	}

	got = info.Schemas[0].Tables[0].Columns[3].Type
	if got != "" {
		t.Errorf("column3 type expected to be empty but got %q", got)
	}

	expected = "abc " + info.Schemas[0].Enums[0].DBName
	got = info.Schemas[0].Enums[0].Name
	if got != expected {
		t.Errorf("enum name expected %q but got %q", expected, got)
	}

	expected = info.Schemas[0].Name
	got = info.Schemas[0].Enums[0].Schema
	if got != expected {
		t.Errorf("enum schema expected %q but got %q", expected, got)
	}

	expected = "abc " + info.Schemas[0].Enums[0].Values[0].DBName
	got = info.Schemas[0].Enums[0].Values[0].Name
	if got != expected {
		t.Errorf("enum value name expected %q but got %q", expected, got)
	}
}
