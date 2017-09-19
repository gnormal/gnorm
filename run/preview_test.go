package run

import (
	"bytes"
	"log"
	"testing"
	"text/template"

	"github.com/google/go-cmp/cmp"
	"gnorm.org/gnorm/database"
	"gnorm.org/gnorm/environ"
	"gnorm.org/gnorm/run/data"
)

func TestMakeTable(t *testing.T) {
	people := []struct {
		Name string
		Age  int
	}{
		{
			Name: "Bob",
			Age:  30,
		},
		{
			Name: "Samantha",
			Age:  3,
		},
	}

	s, err := makeTable(people, "{{.Name}}|{{.Age}}", "Name", "Age")
	if err != nil {
		t.Fatal(err)
	}
	expected := `
+----------+-----+
|   Name   | Age |
+----------+-----+
| Bob      |  30 |
| Samantha |   3 |
+----------+-----+
`[1:]
	if s != expected {
		t.Fatal(cmp.Diff(s, expected))
	}
}

type dummyDriver struct{}

func (dummyDriver) Parse(log *log.Logger, conn string, schemaNames []string, filterTables func(schema, table string) bool) (*database.Info, error) {
	return &database.Info{
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
	}, nil
}

const expectYaml = `schemas:
- name: abc schema
  dbname: schema
  tables:
  - name: abc table
    dbname: table
    columns:
    - name: abc col1
      dbname: col1
      type: INTEGER
      dbtype: int
      isarray: false
      length: 0
      userdefined: false
      nullable: false
      hasdefault: false
    - name: abc col2
      dbname: col2
      type: '*INTEGER'
      dbtype: '*int'
      isarray: false
      length: 0
      userdefined: false
      nullable: true
      hasdefault: false
    - name: abc col3
      dbname: col3
      type: ""
      dbtype: string
      isarray: false
      length: 0
      userdefined: false
      nullable: false
      hasdefault: false
    - name: abc col4
      dbname: col4
      type: ""
      dbtype: '*string'
      isarray: false
      length: 0
      userdefined: false
      nullable: true
      hasdefault: false
  enums:
  - name: abc enum
    dbname: enum
    values:
    - name: abc enumvalue
      dbname: enumvalue
      value: 0
`

const expectTabular = `Schema: abc schema(schema)

Enum: abc enum(schema.enum)
+---------------+-----------+-------+
|     Name      |  DBName   | Value |
+---------------+-----------+-------+
| abc enumvalue | enumvalue |     0 |
+---------------+-----------+-------+


Table: abc table(schema.table)
+----------+--------+----------+---------+---------+--------+-------------+----------+------------+
|   Name   | DBName |   Type   | DBType  | IsArray | Length | UserDefined | Nullable | HasDefault |
+----------+--------+----------+---------+---------+--------+-------------+----------+------------+
| abc col1 | col1   | INTEGER  | int     | false   |      0 | false       | false    | false      |
| abc col2 | col2   | *INTEGER | *int    | false   |      0 | false       | true     | false      |
| abc col3 | col3   |          | string  | false   |      0 | false       | false    | false      |
| abc col4 | col4   |          | *string | false   |      0 | false       | true     | false      |
+----------+--------+----------+---------+---------+--------+-------------+----------+------------+

`

func TestPreviewYaml(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer
	env := environ.Values{
		Stdout: &out,
		Log:    log.New(&errOut, "test: ", log.Lshortfile),
	}

	cfg := &Config{
		NameConversion: template.Must(template.New("").Funcs(environ.FuncMap).Parse(`{{print "abc " .}}`)),
		ConfigData: data.ConfigData{
			NullableTypeMap: map[string]string{
				"*int": "*INTEGER",
			},
			TypeMap: map[string]string{
				"int": "INTEGER",
			},
		},
		Driver: dummyDriver{},
	}
	// with yaml
	if err := Preview(env, cfg, "yaml"); err != nil {
		t.Fatal(err)
	}
	v := out.String()
	if v != expectYaml {
		t.Errorf(cmp.Diff(expectYaml, v))
	}
}

func TestPreviewTabular(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer
	env := environ.Values{
		Stdout: &out,
		Log:    log.New(&errOut, "test: ", log.Lshortfile),
	}

	cfg := &Config{
		NameConversion: template.Must(template.New("").Funcs(environ.FuncMap).Parse(`{{print "abc " .}}`)),
		ConfigData: data.ConfigData{
			NullableTypeMap: map[string]string{
				"*int": "*INTEGER",
			},
			TypeMap: map[string]string{
				"int": "INTEGER",
			},
		},
		Driver: dummyDriver{},
	}

	// tabular
	if err := Preview(env, cfg, "tabular"); err != nil {
		t.Fatal(err)
	}

	v := out.String()
	if v != expectTabular {
		t.Errorf("tabular format differs from expected: %s", cmp.Diff(expectTabular, v))
	}
}

const typesOut = `+---------------+----------------+
| ORIGINAL TYPE | CONVERTED TYPE |
+---------------+----------------+
| *int          | *INTEGER       |
+---------------+----------------+
| *string       |                |
+---------------+----------------+
| int           | INTEGER        |
+---------------+----------------+
| string        |                |
+---------------+----------------+
`

func TestPreviewTypes(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer
	env := environ.Values{
		Stdout: &out,
		Log:    log.New(&errOut, "test: ", log.Lshortfile),
	}

	cfg := &Config{
		NameConversion: template.Must(template.New("").Funcs(environ.FuncMap).Parse(`{{print "abc " .}}`)),
		ConfigData: data.ConfigData{
			NullableTypeMap: map[string]string{
				"*int": "*INTEGER",
			},
			TypeMap: map[string]string{
				"int": "INTEGER",
			},
		},
		Driver: dummyDriver{},
	}
	if err := Preview(env, cfg, "types"); err != nil {
		t.Fatal(err)
	}
	v := out.String()
	if v != typesOut {
		t.Errorf("expected %s got %s", typesOut, v)
	}
}
