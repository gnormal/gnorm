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
	"gnorm.org/gnorm/database/drivers/mysql/pg"
)

// Parse reads the postgres schemas for the given schemas and converts them into
// database.Info structs.
func Parse(log *log.Logger, conn string, schemaNames []string, filterTables func(schema, table string) bool) (*database.Info, error) {
	log.Println("connecting to mysql with DSN", conn)
	db, err := sql.Open("mysql", conn)
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

	schemas := make(map[string]map[string][]*database.Column, len(schemaNames))
	for _, name := range schemaNames {
		schemas[name] = map[string][]*database.Column{}
	}

	for _, t := range tables {
		if !filterTables(t.TableSchema.String, t.TableName.String) {
			continue
		}
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

	enums := map[string][]*database.Enum{}

	for _, c := range columns {
		if !filterTables(c.TableSchema.String, c.TableName.String) {
			continue
		}
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
		col, enum, err := toDBColumn(c, log)
		if err != nil {
			return nil, err
		}
		schema[c.TableName.String] = append(schema[c.TableName.String], col)
		if enum != nil {
			enum.DBTable = c.TableName.String
			enum.DBSchema = c.TableSchema.String
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

func toDBColumn(c *pg.Column, log *log.Logger) (*database.Column, *database.Enum, error) {
	col := &database.Column{
		DBName:     c.ColumnName.String,
		Nullable:   bool(c.IsNullable),
		HasDefault: c.ColumnDefault.String != "",
		DBType:     c.DataType.String,
		Orig:       *c,
	}

	// MySQL always specifies length even if it's not a part of the type. We
	// only really care if it's a part of the type, so check if the size is part
	// of the column_type.
	if strings.HasSuffix(c.ColumnType.String, fmt.Sprintf("(%v)", c.CharacterMaximumLength.Int64)) {
		col.Length = int(c.CharacterMaximumLength.Int64)
	}

	if col.Type != "enum" {
		return col, nil, nil
	}
	// in mysql, enums are specific to a column in a table, so all their data is
	// contained in the column they're used by.

	// column type should be enum('foo', 'bar')
	if len(c.ColumnType.String) < 5 {
		return nil, nil, errors.New("unexpected column type: " + c.ColumnType.String)
	}

	// we'll call the enum the same as the column name.
	// the function above will set the table name etc
	enum := &database.Enum{
		DBName: col.DBName,
	}
	// strip off the enum and parens
	s := c.ColumnType.String[5 : len(c.ColumnType.String)-1]
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
