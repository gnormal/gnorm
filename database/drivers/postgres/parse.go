package postgres // import "gnorm.org/gnorm/database/drivers/postgres"

import (
	"database/sql"
	"log"
	"strconv"
	"strings"

	// register postgres driver
	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"gnorm.org/gnorm/database"
	"gnorm.org/gnorm/database/drivers/postgres/pg"
)

//go:generate go get github.com/xoxo-go/xoxo
//go:generate xoxo pgsql://$DB_USER:$DB_PASSWORD@$DB_HOST/$DB_NAME?sslmode=$DB_SSL_MODE --schema information_schema -o pg --template-path ./templates

func Parse(typeMap, nullableTypeMap map[string]string, log *log.Logger, conn string, schemaNames []string) (*database.Info, error) {
	log.Println("connecting to postgres with DSN", conn)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	sch := make([]sql.NullString, len(schemaNames))
	for x := range schemaNames {
		sch[x] = sql.NullString{String: schemaNames[x], Valid: true}
	}
	log.Println("querying table schemas for", schemaNames)
	tables, err := pg.QueryTable(db, pg.TableTableSchemaWhere.In(sch), pg.UnOrdered)
	if err != nil {
		return nil, err
	}

	schemas := make(map[string]map[string][]database.Column, len(schemaNames))
	for _, name := range schemaNames {
		schemas[name] = map[string][]database.Column{}
	}

	for _, t := range tables {
		s, ok := schemas[t.TableSchema.String]
		if !ok {
			log.Printf("Should be impossible: table %q references unknown schema %q", t.TableName.String, t.TableSchema.String)
			continue
		}
		s[t.TableName.String] = nil
	}

	columns, err := pg.QueryColumn(db, pg.ColumnTableSchemaWhere.In(sch), pg.UnOrdered)
	if err != nil {
		return nil, err
	}

	for _, c := range columns {
		schema, ok := schemas[c.TableSchema.String]
		if !ok {
			log.Printf("Should be impossible: column %q references unknown schema %q", c.ColumnName.String, c.TableSchema.String)
			continue
		}
		_, ok = schema[c.TableName.String]
		if !ok {
			log.Printf("Should be impossible: column %q references unknown table %q in schema %q", c.ColumnName.String, c.TableName.String, c.TableSchema.String)
			continue
		}

		col := toDBColumn(c, typeMap, nullableTypeMap, log)
		schema[c.TableName.String] = append(schema[c.TableName.String], col)
	}

	res := &database.Info{Schemas: make([]database.Schema, 0, len(schemas))}
	for name, tables := range schemas {
		s := database.Schema{Name: name}
		for tname, columns := range tables {
			s.Tables = append(s.Tables, database.Table{Name: tname, Columns: columns})
		}
		res.Schemas = append(res.Schemas, s)
	}

	return res, nil
}

func toDBColumn(c *pg.Column, typeMap, nullableTypeMap map[string]string, log *log.Logger) database.Column {
	col := database.Column{
		Name:     c.ColumnName.String,
		Nullable: bool(c.IsNullable),
		Orig:     *c,
	}

	typ := c.DataType.String
	switch typ {
	case "ARRAY":
		col.IsArray = true
		// when it's an array, postges prepends an underscore to the standard
		// name.
		typ = c.UdtName.String[1:]

	case "USER-DEFINED":
		col.UserDefined = true
		typ = c.UdtName.String
	}

	length, newtyp, err := calculateLength(typ)
	switch {
	case err != nil:
		log.Println(err)
	case length > 0:
		col.Length = length
		typ = newtyp
	}
	col.DBType = typ

	if col.Nullable {
		t, ok := nullableTypeMap[typ]
		if !ok {
			log.Println("unmapped nullable type:", typ)
		} else {
			col.Type = t
		}
	} else {
		t, ok := typeMap[typ]
		if !ok {
			log.Println("unmapped non-nullable type:", typ)
		} else {
			col.Type = t
		}
	}
	return col
}

func calculateLength(typ string) (length int, newtyp string, err error) {
	idx := strings.Index(typ, "[")
	if idx == -1 {
		// no length
		return 0, "", nil
	}
	end := strings.LastIndex(typ, "]")
	// we expect the length of the type to be the end of the name.
	if end == len(typ)-1 {
		lstr := typ[idx+1 : end]
		l, err := strconv.Atoi(lstr)
		if err != nil {
			return 0, "", err
		}
		return l, typ[:idx], nil
	}
	// something wonky with the brackets
	return 0, "", errors.New("unknown bracket format in type name")
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
