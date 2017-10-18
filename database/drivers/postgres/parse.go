package postgres // import "gnorm.org/gnorm/database/drivers/postgres"

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	// register postgres driver
	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"gnorm.org/gnorm/database"
	"gnorm.org/gnorm/database/drivers/postgres/gnorm/columns"
	"gnorm.org/gnorm/database/drivers/postgres/gnorm/tables"
)

// PG implements drivers.Driver interface for interacting with postgresql
// database.
type PG struct{}

// Parse reads the postgres schemas for the given schemas and converts them into
// database.Info structs.
func (PG) Parse(log *log.Logger, conn string, schemaNames []string, filterTables func(schema, table string) bool) (*database.Info, error) {
	return parse(log, conn, schemaNames, filterTables)
}

func parse(log *log.Logger, conn string, schemaNames []string, filterTables func(schema, table string) bool) (*database.Info, error) {
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
	tables, err := tables.Query(db, tables.TableSchemaCol.In(sch))
	if err != nil {
		return nil, err
	}

	log.Printf("found %v tables", len(tables))
	schemas := make(map[string]map[string][]*database.Column, len(schemaNames))
	for _, name := range schemaNames {
		schemas[name] = map[string][]*database.Column{}
	}

	for _, t := range tables {
		if !filterTables(t.TableSchema.String, t.TableName.String) {
			log.Printf("skipping filtered-out table %v.%v", t.TableSchema.String, t.TableName.String)
			continue
		}

		s, ok := schemas[t.TableSchema.String]
		if !ok {
			log.Printf("Should be impossible: table %q references unknown schema %q", t.TableName.String, t.TableSchema.String)
			continue
		}
		s[t.TableName.String] = nil
	}

	columns, err := columns.Query(db, columns.TableSchemaCol.In(sch))
	if err != nil {
		return nil, err
	}
	log.Printf("found %v columns for all tables in all specified schemas", len(columns))
	for _, c := range columns {
		if !filterTables(c.TableSchema.String, c.TableName.String) {
			log.Printf("skipping column %q because it is for filtered-out table %v.%v", c.ColumnName.String, c.TableSchema.String, c.TableName.String)
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

		col := toDBColumn(c, log)
		schema[c.TableName.String] = append(schema[c.TableName.String], col)
	}

	primaryKeys, err := queryPrimaryKeys(log, db, schemaNames)
	if err != nil {
		return nil, err
	}
	log.Printf("found %v primary keys", len(primaryKeys))
	for _, pk := range primaryKeys {
		if !filterTables(pk.SchemaName, pk.TableName) {
			log.Printf("skipping constraint %q because it is for filtered-out table %v.%v", pk.Name, pk.SchemaName, pk.TableName)
			continue
		}

		schema, ok := schemas[pk.SchemaName]
		if !ok {
			log.Printf("Should be impossible: constraint %q references unknown schema %q", pk.Name, pk.SchemaName)
			continue
		}
		table, ok := schema[pk.TableName]
		if !ok {
			log.Printf("Should be impossible: constraint %q references unknown table %q in schema %q", pk.Name, pk.TableName, pk.SchemaName)
			continue
		}

		for _, col := range table {
			if pk.ColumnName != col.Name {
				continue
			}
			col.IsPrimaryKey = true
		}
	}

	enums, err := queryEnums(log, db, schemaNames)
	if err != nil {
		return nil, err
	}
	log.Printf("found %v enums for all schemas", len(enums))
	res := &database.Info{Schemas: make([]*database.Schema, 0, len(schemas))}
	for _, schema := range schemaNames {
		tables := schemas[schema]
		s := &database.Schema{
			Name:  schema,
			Enums: enums[schema],
		}
		for tname, columns := range tables {
			s.Tables = append(s.Tables, &database.Table{Name: tname, Columns: columns})
		}
		res.Schemas = append(res.Schemas, s)
	}

	return res, nil
}

func toDBColumn(c *columns.Row, log *log.Logger) *database.Column {
	col := &database.Column{
		Name:       c.ColumnName.String,
		Nullable:   c.IsNullable.String == "YES",
		HasDefault: c.ColumnDefault.String != "",
		Length:     int(c.CharacterMaximumLength.Int64),
		Orig:       *c,
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

	col.Type = typ

	return col
}

func queryPrimaryKeys(log *log.Logger, db *sql.DB, schemas []string) ([]*database.PrimaryKey, error) {
	// TODO: make this work with Gnorm generated types
	const q = `
	SELECT k.table_schema, k.table_name, k.column_name, k.constraint_name
	FROM information_schema.key_column_usage k
	LEFT JOIN information_schema.table_constraints c
    	ON k.table_schema = c.table_schema
    	AND k.table_name = c.table_name
    	AND k.constraint_name = c.constraint_name
	WHERE c.constraint_type='PRIMARY KEY' AND k.table_schema IN (%s)`
	spots := make([]string, len(schemas))
	vals := make([]interface{}, len(schemas))
	for x := range schemas {
		spots[x] = fmt.Sprintf("$%v", x+1)
		vals[x] = schemas[x]
	}
	query := fmt.Sprintf(q, strings.Join(spots, ", "))
	rows, err := db.Query(query, vals...)
	defer rows.Close()
	if err != nil {
		return nil, errors.WithMessage(err, "error querying keys")
	}
	var ret []*database.PrimaryKey

	for rows.Next() {
		kc := &database.PrimaryKey{}
		if err := rows.Scan(&kc.SchemaName, &kc.TableName, &kc.ColumnName, &kc.Name); err != nil {
			return nil, errors.WithMessage(err, "error scanning key constraint")
		}
		ret = append(ret, kc)
	}
	return ret, nil
}

func queryEnums(log *log.Logger, db *sql.DB, schemas []string) (map[string][]*database.Enum, error) {
	// TODO: make this work with Gnorm generated types
	const q = `
	SELECT      n.nspname, t.typname as type 
	FROM        pg_type t 
	LEFT JOIN   pg_catalog.pg_namespace n ON n.oid = t.typnamespace 
	WHERE       (t.typrelid = 0 OR (SELECT c.relkind = 'c' FROM pg_catalog.pg_class c WHERE c.oid = t.typrelid)) 
	AND     NOT EXISTS(SELECT 1 FROM pg_catalog.pg_type el WHERE el.oid = t.typelem AND el.typarray = t.oid)
	AND     n.nspname IN (%s)`
	spots := make([]string, len(schemas))
	vals := make([]interface{}, len(schemas))
	for x := range schemas {
		spots[x] = fmt.Sprintf("$%v", x+1)
		vals[x] = schemas[x]
	}
	query := fmt.Sprintf(q, strings.Join(spots, ", "))
	rows, err := db.Query(query, vals...)
	defer rows.Close()
	if err != nil {
		return nil, errors.WithMessage(err, "error querying enum names")
	}
	ret := map[string][]*database.Enum{}
	for rows.Next() {
		var name, schema string
		if err := rows.Scan(&schema, &name); err != nil {
			return nil, errors.WithMessage(err, "error scanning enum name into string")
		}
		vals, err := queryValues(log, db, schema, name)
		if err != nil {
			return nil, err
		}
		enum := &database.Enum{
			Name:   name,
			Values: vals,
		}
		ret[schema] = append(ret[schema], enum)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.WithMessage(err, "error reading enum names")
	}
	return ret, nil
}

func queryValues(log *log.Logger, db *sql.DB, schema, enum string) ([]*database.EnumValue, error) {
	// TODO: make this work with Gnorm generated types
	rows, err := db.Query(`
	SELECT
	e.enumlabel,
	e.enumsortorder
	FROM pg_type t
	JOIN ONLY pg_namespace n ON n.oid = t.typnamespace
	LEFT JOIN pg_enum e ON t.oid = e.enumtypid
	WHERE n.nspname = $1 AND t.typname = $2`, schema, enum)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query enum values for %s.%s", schema, enum)
	}
	defer rows.Close()
	var vals []*database.EnumValue
	for rows.Next() {
		var name sql.NullString
		var val sql.NullInt64
		if err := rows.Scan(&name, &val); err != nil {
			return nil, errors.Wrapf(err, "failed reading enum values for %s.%s", schema, enum)
		}
		vals = append(vals, &database.EnumValue{Name: name.String, Value: int(val.Int64)})
	}
	log.Printf("found %d values for enum %v.%v", len(vals), schema, enum)
	return vals, nil
}
