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

func TestForeignKeyRefs(t *testing.T) {
	t.Parallel()

	c := &Config{
		NameConversion: template.Must(template.New("").Funcs(environ.FuncMap).Parse(`{{.}}`)),
		ConfigData: data.ConfigData{
			TypeMap: map[string]string{
				"int": "INTEGER",
			},
		},
	}

	info := &database.Info{
		Schemas: []*database.Schema{{
			Name: "public",
			Tables: []*database.Table{
				{
					Name: "tbl_1",
					Columns: []*database.Column{
						{
							Name:         "col_1",
							Type:         "int",
							IsForeignKey: true,
							ForeignKey: &database.ForeignKey{
								Name:              "fk_1",
								SchemaName:        "public",
								TableName:         "tbl_1",
								ColumnName:        "col_1",
								ForeignTableName:  "tbl_2",
								ForeignColumnName: "col_1",
							},
						},
						{
							Name:         "col_2",
							Type:         "int",
							IsForeignKey: true,
							ForeignKey: &database.ForeignKey{
								Name:              "fk_1",
								SchemaName:        "public",
								TableName:         "tbl_1",
								ColumnName:        "col_2",
								ForeignTableName:  "tbl_2",
								ForeignColumnName: "col_2",
							},
						},
					},
				},
				{
					Name: "tbl_2",
					Columns: []*database.Column{
						{
							Name: "col_1",
							Type: "int",
						},
						{
							Name: "col_2",
							Type: "int",
						},
						{
							Name:         "col_3",
							Type:         "int",
							IsForeignKey: true,
							ForeignKey: &database.ForeignKey{
								Name:              "fk_2",
								SchemaName:        "public",
								TableName:         "tbl_2",
								ColumnName:        "col_3",
								ForeignTableName:  "tbl_3",
								ForeignColumnName: "col_1",
							},
						},
					},
				},
				{
					Name: "tbl_3",
					Columns: []*database.Column{
						{
							Name: "col_1",
							Type: "int",
						},
					},
				},
			},
		}},
	}

	log := log.New(&bytes.Buffer{}, "", 0)
	data, err := makeData(log, info, c)
	if err != nil {
		t.Fatalf("unexpected err: %s", err)
	}

	tbl1FKs := data.Schemas[0].TablesByName["tbl_1"].ForeignKeys
	if l := len(tbl1FKs); l != 1 {
		t.Fatalf("incorrect number of foreign keys; expected %d, got %d", 1, l)
	}
	if name := tbl1FKs[0].Name; name != "fk_1" {
		t.Fatalf("incorrect foreign key name; expected %s, got %s", "fk_1", name)
	}
	if l := len(tbl1FKs[0].FKColumns); l != 2 {
		t.Fatalf("too many foreign key columns; expected %d, got %d", 2, l)
	}
	if name := tbl1FKs[0].RefTable.Name; name != "tbl_2" {
		t.Fatalf("incorrect foreign key table; expected %s, got %s", "tbl_2", name)
	}
	if name := tbl1FKs[0].FKColumns[0].Column.Name; name != "col_1" {
		t.Fatalf("incorrect foreign key column name; expected %s, got %s", "col_1", name)
	}
	if name := tbl1FKs[0].FKColumns[1].Column.Name; name != "col_2" {
		t.Fatalf("incorrect foreign key column name; expected %s, got %s", "col_2", name)
	}

	tbl2FKRefs := data.Schemas[0].TablesByName["tbl_2"].ForeignKeyRefs
	if l := len(tbl2FKRefs); l != 1 {
		t.Fatalf("incorrect number of foreign key refs; expected %d, got %d", 1, l)
	}
	if name := tbl2FKRefs[0].Name; name != "fk_1" {
		t.Fatalf("incorrect foreign key ref; expected %s, got %s", "fk_1", name)
	}
	if name := tbl2FKRefs[0].Table.Name; name != "tbl_1" {
		t.Fatalf("incorrect foreign key ref table; expected %s, got %s", "tbl_1", name)
	}
	if l := len(tbl2FKRefs[0].FKColumns); l != 2 {
		t.Fatalf("incorrect number of foreign key ref columns; expected %d, got %d", 2, l)
	}
	if name := tbl2FKRefs[0].FKColumns[0].Column.Name; name != "col_1" {
		t.Fatalf("incorrect foreign key ref column; expected %s, got %s", "col_1", name)
	}
	if name := tbl2FKRefs[0].FKColumns[1].Column.Name; name != "col_2" {
		t.Fatalf("incorrect foreign key ref column; expected %s, got %s", "col_2", name)
	}

	tbl2FKs := data.Schemas[0].TablesByName["tbl_2"].ForeignKeys
	if l := len(tbl2FKs); l != 1 {
		t.Fatalf("incorrect number of foreign keys; expected %d, got %d", 1, l)
	}
	if name := tbl2FKs[0].Name; name != "fk_2" {
		t.Fatalf("incorrect foreign key name; expected %s, got %s", "fk_2", name)
	}
	if l := len(tbl2FKs[0].FKColumns); l != 1 {
		t.Fatalf("too many foreign key columns; expected %d, got %d", 1, l)
	}
	if name := tbl2FKs[0].RefTable.Name; name != "tbl_3" {
		t.Fatalf("incorrect foreign key table; expected %s, got %s", "tbl_3", name)
	}
	if name := tbl2FKs[0].FKColumns[0].Column.Name; name != "col_3" {
		t.Fatalf("incorrect foreign key column name; expected %s, got %s", "col_3", name)
	}

	// See: https://github.com/gnormal/gnorm/issues/123
	tbl3FKRefs := data.Schemas[0].TablesByName["tbl_3"].ForeignKeyRefs
	if l := len(tbl3FKRefs); l != 1 {
		t.Fatalf("incorrect number of foreign keys refs; expected %d, got %d", 1, l)
	}
	if name := tbl3FKRefs[0].Name; name != "fk_2" {
		t.Fatalf("incorrect foreign key ref; expected %s, got %s", "fk_2", name)
	}
	if name := tbl3FKRefs[0].Table.Name; name != "tbl_2" {
		t.Fatalf("incorrect foreign key ref table; expected %s, got %s", "tbl_2", name)
	}
	if l := len(tbl3FKRefs[0].FKColumns); l != 1 {
		t.Fatalf("incorrect number of foreign key ref columns; expected %d, got %d", 1, l)
	}
	if name := tbl3FKRefs[0].FKColumns[0].Column.Name; name != "col_3" {
		t.Fatalf("incorrect foreign key ref column; expected %s, got %s", "col_3", name)
	}
}
