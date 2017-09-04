package mysql // import "gnorm.org/gnorm/database/drivers/mysql"

import (
	"database/sql"
	"log"
	"strconv"
	"strings"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"

	"gnorm.org/gnorm/database"
	"gnorm.org/gnorm/database/drivers/mysql/pg"
)

// Parse reads the postgres schemas for the given schemas and converts them into
// database.Info structs.
func Parse(log *log.Logger, conn string, schemaNames []string) (*database.Info, error) {
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
		Orig:       *c,
	}

	typ := c.DataType.String

	var enum *database.Enum

	// in mysql, enums are specific to a column in a table
	if typ == "enum" {
		// column type should be enum('foo', 'bar')
		if len(c.ColumnType.String) < 5 {
			return nil, nil, errors.New("unexpected column type: " + c.ColumnType.String)
		}
		// we'll call the enum the same as the column name.
		// the function above will set the table name etc
		enum = &database.Enum{
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
	}

	length, newtyp, err := calculateLength(typ)
	switch {
	case err != nil:
		return nil, nil, err
	case length > 0:
		col.Length = length
		typ = newtyp
	}
	col.DBType = typ

	return col, enum, nil
}

// calculateLength tries to convert a type that contains a length specification
// to a length number and a type name without the brackets.  Thus varchar(32)
// would return 32, "varchar".  It's up to the consumer to understand that
// sometimes length is a maximum and sometimes it's a requirement (i.e.
// varchar(32) vs char(32), since this information is intrinsic to the type
// name.
func calculateLength(typ string) (length int, newtyp string, err error) {
	idx := strings.Index(typ, "(")
	if idx == -1 {
		// no length indicated
		return 0, "", nil
	}
	end := strings.LastIndex(typ, ")")
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
