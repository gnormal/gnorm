package database // import "gnorm.org/gnorm/database"

import (
	"database/sql"

	// import here for now
	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"gnorm.org/gnorm/database/pg"
	"gnorm.org/gnorm/environ"
)

//go:generate go get github.com/xoxo-go/xoxo
//go:generate xoxo pgsql://$DB_USER:$DB_PASSWORD@$DB_HOST/$DB_NAME?sslmode=$DB_SSL_MODE --schema information_schema -o pg --template-path ./templates

type SchemaInfo struct {
	Schemas []Schema
}

type Schema struct {
	Name   string
	Tables []Table
	Enums  []Enum
}

type Enum struct {
	Name   string
	Values []EnumValue
}

type EnumValue struct {
	Name  string
	Value int
}

type Table struct {
	Name    string
	Columns []Column
}

type Column struct {
	Name        string
	Type        string
	OrigType    string
	UserDefined bool
	Nullable    bool
	Orig        pg.Column //`yaml:"-"`
}

func Parse(typeMap, nullableTypeMap map[string]string, env environ.Values, conn string, schemaNames []string) (*SchemaInfo, error) {
	env.Log.Println("connecting to", "postgres", "with DSN", conn)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	sch := make([]sql.NullString, len(schemaNames))
	for x := range schemaNames {
		sch[x] = sql.NullString{String: schemaNames[x], Valid: true}
	}
	env.Log.Println("querying table schemas for", schemaNames)
	tables, err := pg.QueryTable(db, pg.TableTableSchemaWhere.In(sch), pg.UnOrdered)
	if err != nil {
		return nil, err
	}

	schemas := make(map[string]map[string][]Column, len(schemaNames))
	for _, name := range schemaNames {
		schemas[name] = map[string][]Column{}
	}

	for _, t := range tables {
		s, ok := schemas[t.TableSchema.String]
		if !ok {
			env.Log.Printf("Should be impossible: table %q references unknown schema %q", t.TableName.String, t.TableSchema.String)
			continue
		}
		s[t.TableName.String] = nil
	}

	env.Log.Println("querying columns")
	columns, err := pg.QueryColumn(db, pg.ColumnTableSchemaWhere.In(sch), pg.UnOrdered)
	if err != nil {
		return nil, err
	}

	env.Log.Printf("parsing %v columns", len(columns))
	for _, c := range columns {
		schema, ok := schemas[c.TableSchema.String]
		if !ok {
			env.Log.Printf("Should be impossible: column %q references unknown schema %q", c.ColumnName.String, c.TableSchema.String)
			continue
		}
		_, ok = schema[c.TableName.String]
		if !ok {
			env.Log.Printf("Should be impossible: column %q references unknown table %q in schema %q", c.ColumnName.String, c.TableName.String, c.TableSchema.String)
			continue
		}
		origtyp := c.DataType.String
		var userDefined bool
		if c.DataType.String == "USER-DEFINED" {
			userDefined = true
			origtyp = c.UdtName.String
		}

		var gotype string
		if !userDefined {
			if c.IsNullable {
				gotype = nullableTypeMap[origtyp]
			} else {
				gotype = typeMap[origtyp]
			}
		}
		col := Column{
			Name:        c.ColumnName.String,
			OrigType:    origtyp,
			Type:        gotype,
			UserDefined: userDefined,
			Nullable:    bool(c.IsNullable),
			Orig:        *c,
		}
		schema[c.TableName.String] = append(schema[c.TableName.String], col)
	}

	res := &SchemaInfo{Schemas: make([]Schema, 0, len(schemas))}
	for name, tables := range schemas {
		s := Schema{Name: name}
		for tname, columns := range tables {
			s.Tables = append(s.Tables, Table{Name: tname, Columns: columns})
		}
		res.Schemas = append(res.Schemas, s)
	}

	return res, nil
}

/*

SELECT
DISTINCT column_name AS enum_name
FROM information_schema.columns
WHERE data_type = 'enum' AND table_schema = `public`


SELECT
e.enumlabel,
e.enumsortorder
FROM pg_type t
JOIN ONLY pg_namespace n ON n.oid = t.typnamespace
LEFT JOIN pg_enum e ON t.oid = e.enumtypid
WHERE n.nspname = 'public' AND t.typname = 'user_role';

*/
