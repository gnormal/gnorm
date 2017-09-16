package run

import (
	"bytes"
	"log"
	"testing"
	"text/template"

	"github.com/google/go-cmp/cmp"
	"gnorm.org/gnorm/database"
	"gnorm.org/gnorm/environ"
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
	}, nil
}

const expectYaml = `schemas:
- name: abc schema
  dbname: schema
  tables:
  - schema: abc schema
    dbschema: schema
    name: abc table
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
  - schema: abc schema
    dbschema: schema
    table: ""
    dbtable: ""
    name: abc enum
    dbname: enum
    values:
    - name: abc enumvalue
      dbname: enumvalue
      value: 0
`

const expectNoYaml = `Schema: abc schema(schema)

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

func TestPreview(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer
	env := environ.Values{
		Stdout: &out,
		Log:    log.New(&errOut, "test: ", log.Lshortfile),
	}

	cfg := &Config{
		NameConversion: template.Must(template.New("").Funcs(environ.FuncMap).Parse(`{{print "abc " .}}`)),
		NullableTypeMap: map[string]string{
			"*int": "*INTEGER",
		},
		TypeMap: map[string]string{
			"int": "INTEGER",
		},
		Driver: dummyDriver{},
	}
	// with yaml
	if err := Preview(env, cfg, true); err != nil {
		t.Fatal(err)
	}
	v := out.String()
	if v != expectYaml {
		t.Errorf("expected %s got %s", expectYaml, v)
	}

	//without yaml
	out.Reset()
	errOut.Reset()
	if err := Preview(env, cfg, false); err != nil {
		t.Fatal(err)
	}
	v = out.String()
	if v != expectNoYaml {
		t.Errorf("expected %s got %s", expectNoYaml, v)
	}
}
