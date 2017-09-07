package mysql // import "gnorm.org/gnorm/database/drivers/mysql"

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"

	"gnorm.org/gnorm/database"
	"gnorm.org/gnorm/database/drivers/mysql/gnorm/columns"
	"gnorm.org/gnorm/database/drivers/mysql/gnorm/tables"
)

// Parse reads the postgres schemas for the given schemas and converts them into
// database.Info structs.
func Parse(log *log.Logger, conn string, schemaNames []string, filterTables func(schema, table string) bool) (*database.Info, error) {
	log.Println("connecting to mysql with DSN", conn)
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	log.Println("querying table schemas for", schemaNames)
	tables, err := tables.Query(db, tables.TableSchemaCol.In(schemaNames))
	if err != nil {
		return nil, err
	}

	schemas := make(map[string]map[string]database.Columns, len(schemaNames))
	for _, name := range schemaNames {
		schemas[name] = map[string]database.Columns{}
	}

	for _, t := range tables {
		if !filterTables(t.TableSchema, t.TableName) {
			continue
		}
		s, ok := schemas[t.TableSchema]
		if !ok {
			log.Printf("Should be impossible: table %q references unknown schema %q", t.TableName, t.TableSchema)
			continue
		}
		s[t.TableName] = nil
	}

	columns, err := columns.Query(db, columns.TableSchemaCol.In(schemaNames))
	if err != nil {
		return nil, err
	}

	enums := map[string][]*database.Enum{}

	for _, c := range columns {
		if !filterTables(c.TableSchema, c.TableName) {
			continue
		}
		schema, ok := schemas[c.TableSchema]
		if !ok {
			log.Printf("Should be impossible: column %q references unknown schema %q", c.ColumnName, c.TableSchema)
			continue
		}
		_, ok = schema[c.TableName]
		if !ok {
			log.Printf("Should be impossible: column %q references unknown table %q in schema %q", c.ColumnName, c.TableName, c.TableSchema)
			continue
		}
		col, enum, err := toDBColumn(c, log)
		if err != nil {
			return nil, err
		}
		schema[c.TableName] = append(schema[c.TableName], col)
		if enum != nil {
			enum.DBTable = c.TableName
			enum.DBSchema = c.TableSchema
			enums[enum.DBSchema] = append(enums[enum.DBSchema], enum)
		}
	}

	res := &database.Info{Schemas: make([]*database.Schema, 0, len(schemas))}
	for _, schema := range schemaNames {
		tables := schemas[schema]
		s := &database.Schema{
			DBName: schema,
			Enums:  enums[schema],
		}
		for tname, columns := range tables {
			s.Tables = append(s.Tables, &database.Table{DBName: tname, DBSchema: schema, Columns: columns})
		}
		res.Schemas = append(res.Schemas, s)
	}

	return res, nil
}

func toDBColumn(c *columns.Row, log *log.Logger) (*database.Column, *database.Enum, error) {
	col := &database.Column{
		DBName:     c.ColumnName,
		Nullable:   c.IsNullable == "YES",
		HasDefault: c.ColumnDefault.String != "",
		DBType:     c.DataType,
		Orig:       *c,
	}

	// MySQL always specifies length even if it's not a part of the type. We
	// only really care if it's a part of the type, so check if the size is part
	// of the column_type.
	if strings.HasSuffix(c.ColumnType, fmt.Sprintf("(%v)", c.CharacterMaximumLength.Int64)) {
		col.Length = int(c.CharacterMaximumLength.Int64)
	}

	if col.Type != "enum" {
		return col, nil, nil
	}
	// in mysql, enums are specific to a column in a table, so all their data is
	// contained in the column they're used by.

	// column type should be enum('foo', 'bar')
	if len(c.ColumnType) < 5 {
		return nil, nil, errors.New("unexpected column type: " + c.ColumnType)
	}

	// we'll call the enum the same as the column name.
	// the function above will set the table name etc
	enum := &database.Enum{
		DBName: col.DBName,
	}
	// strip off the enum and parens
	s := c.ColumnType[5 : len(c.ColumnType)-1]
	vals := strings.Split(s, ",")
	enum.Values = make([]*database.EnumValue, len(vals))
	for x := range vals {
		enum.Values[x] = &database.EnumValue{
			// strip off the quotes
			DBName: vals[x][1 : len(vals[x])-1],
			// enum values start at 1 in mysql
			Value: x + 1,
		}
	}

	return col, enum, nil
}
