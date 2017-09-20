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

func TestPreviewYAML(t *testing.T) {
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
	if err := Preview(env, cfg, PreviewYAML); err != nil {
		t.Fatal(err)
	}
	v := out.String()
	if v != expectYaml {
		t.Errorf(cmp.Diff(expectYaml, v))
	}
}

var expectJSON = `
{
  "Schemas": [
    {
      "Name": "abc schema",
      "DBName": "schema",
      "Tables": [
        {
          "Name": "abc table",
          "DBName": "table",
          "Columns": [
            {
              "Name": "abc col1",
              "DBName": "col1",
              "Type": "INTEGER",
              "DBType": "int",
              "IsArray": false,
              "Length": 0,
              "UserDefined": false,
              "Nullable": false,
              "HasDefault": false
            },
            {
              "Name": "abc col2",
              "DBName": "col2",
              "Type": "*INTEGER",
              "DBType": "*int",
              "IsArray": false,
              "Length": 0,
              "UserDefined": false,
              "Nullable": true,
              "HasDefault": false
            },
            {
              "Name": "abc col3",
              "DBName": "col3",
              "Type": "",
              "DBType": "string",
              "IsArray": false,
              "Length": 0,
              "UserDefined": false,
              "Nullable": false,
              "HasDefault": false
            },
            {
              "Name": "abc col4",
              "DBName": "col4",
              "Type": "",
              "DBType": "*string",
              "IsArray": false,
              "Length": 0,
              "UserDefined": false,
              "Nullable": true,
              "HasDefault": false
            }
          ]
        }
      ],
      "Enums": [
        {
          "Name": "abc enum",
          "DBName": "enum",
          "Values": [
            {
              "Name": "abc enumvalue",
              "DBName": "enumvalue",
              "Value": 0
            }
          ]
        }
      ]
    }
  ]
}`[1:]

func TestPreviewJSON(t *testing.T) {
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
	// with json
	if err := Preview(env, cfg, PreviewJSON); err != nil {
		t.Fatal(err)
	}
	v := out.String()
	if v != expectJSON {
		t.Errorf(cmp.Diff(expectJSON, v))
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
	if err := Preview(env, cfg, PreviewTabular); err != nil {
		t.Fatal(err)
	}

	v := out.String()
	if v != expectTabular {
		t.Errorf("tabular format differs from expected: %s", cmp.Diff(expectTabular, v))
	}
}

var typesOut = `
[TypeMap]
"int" = "INTEGER"
"string" = ""

[NullableTypeMap]
"*int" = "*INTEGER"
"*string" = ""
`[1:]

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
	if err := Preview(env, cfg, PreviewTypes); err != nil {
		t.Fatal(err)
	}
	v := out.String()
	if v != typesOut {
		t.Errorf("expected %s got %s", typesOut, v)
	}
}
