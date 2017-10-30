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
	"gnorm.org/gnorm/database/drivers/mysql/gnorm/statistics"
	"gnorm.org/gnorm/database/drivers/mysql/gnorm/tables"
)

//go:generate gnorm gen

// MySQL implements drivers.Driver interface for MySQL database.
type MySQL struct{}

// Parse reads the mysql schemas for the given schemas and converts them into
// database.Info structs.
func (MySQL) Parse(log *log.Logger, conn string, schemaNames []string, filterTables func(schema, table string) bool) (*database.Info, error) {
	return parse(log, conn, schemaNames, filterTables)
}

func parse(log *log.Logger, conn string, schemaNames []string, filterTables func(schema, table string) bool) (*database.Info, error) {
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

	schemas := make(map[string]map[string][]*database.Column, len(schemaNames))
	for _, name := range schemaNames {
		schemas[name] = map[string][]*database.Column{}
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
			enum.Table = c.TableName
			enums[c.TableSchema] = append(enums[c.TableSchema], enum)
		}
	}

	indexes := make(map[string]map[string]map[string][]*database.Column)

	statistics, err := statistics.Query(db, statistics.TableSchemaCol.In(schemaNames))
	if err != nil {
		return nil, err
	}

	for _, s := range statistics {
		if !filterTables(s.TableSchema, s.TableName) {
			continue
		}

		schema, ok := schemas[s.TableSchema]
		if !ok {
			log.Printf("Should be impossible: index %q references unknown schema %q", s.IndexName, s.TableSchema)
			continue
		}

		table, ok := schema[s.TableName]
		if !ok {
			log.Printf("Should be impossible: index %q references unknown table %q", s.IndexName, s.TableName)
			continue
		}

		var column *database.Column
		for _, c := range table {
			if c.Name == s.ColumnName {
				column = c
				break
			}
		}
		if column == nil {
			log.Printf("Should be impossible: index %q references unknown column %q", s.IndexName, s.ColumnName)
			continue
		}

		schemaIndex, ok := indexes[s.TableSchema]
		if !ok {
			schemaIndex = make(map[string]map[string][]*database.Column)
			indexes[s.TableSchema] = schemaIndex
		}

		tableIndex, ok := schemaIndex[s.TableName]
		if !ok {
			tableIndex = make(map[string][]*database.Column)
			schemaIndex[s.TableName] = tableIndex
		}
		tableIndex[s.IndexName] = append(tableIndex[s.IndexName], column)
	}

	res := &database.Info{Schemas: make([]*database.Schema, 0, len(schemas))}
	for _, schema := range schemaNames {
		tables := schemas[schema]
		s := &database.Schema{
			Name:  schema,
			Enums: enums[schema],
		}

		dbtables := make(map[string]*database.Table, len(tables))
		for tname, columns := range tables {
			dbtables[tname] = &database.Table{Name: tname, Columns: columns}
		}
		for tname, index := range indexes[schema] {
			dbtables[tname].Indexes = make([]*database.Index, 0, len(indexes))
			for iname, columns := range index {
				dbtables[tname].Indexes = append(dbtables[tname].Indexes, &database.Index{Name: iname, Columns: columns})
			}
		}
		for _, table := range dbtables {
			s.Tables = append(s.Tables, table)
		}

		res.Schemas = append(res.Schemas, s)
	}

	return res, nil
}

func toDBColumn(c *columns.Row, log *log.Logger) (*database.Column, *database.Enum, error) {
	col := &database.Column{
		Name:         c.ColumnName,
		Nullable:     c.IsNullable == "YES",
		HasDefault:   c.ColumnDefault.String != "",
		Type:         c.DataType,
		Orig:         *c,
		IsPrimaryKey: strings.Contains(c.ColumnKey, "PRI"),
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
		Name: col.Name,
	}
	// strip off the enum and parens
	s := c.ColumnType[5 : len(c.ColumnType)-1]
	vals := strings.Split(s, ",")
	enum.Values = make([]*database.EnumValue, len(vals))
	for x := range vals {
		enum.Values[x] = &database.EnumValue{
			// strip off the quotes
			Name: vals[x][1 : len(vals[x])-1],
			// enum values start at 1 in mysql
			Value: x + 1,
		}
	}

	return col, enum, nil
}
