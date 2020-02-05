package mssql // import "gnorm.org/gnorm/database/drivers/mssql"

import (
	"database/sql"
	"gnorm.org/gnorm/database/drivers/mssql/gnorm"
	"log"
	// mssql driver
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/pkg/errors"

	"gnorm.org/gnorm/database"
	"gnorm.org/gnorm/database/drivers/mssql/gnorm/columns"
	"gnorm.org/gnorm/database/drivers/mssql/gnorm/statistics"
	"gnorm.org/gnorm/database/drivers/mssql/gnorm/tables"
)

//go:generate gnorm gen

// Mssql implements drivers.Driver interface for Mssql database.
type Mssql struct{}

// Parse reads the mssql schemas for the given schemas and converts them into
// database.Info structs.
func (Mssql) Parse(log *log.Logger, conn string, schemaNames []string, filterTables func(schema, table string) bool) (*database.Info, error) {

	log.Println("connecting to Mssql with DSN", conn)
	db, err := sql.Open("sqlserver", conn)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	log.Println("querying table schemas for", schemaNames)

	res := &database.Info{Schemas: make([]*database.Schema, 0, len(schemaNames))}
	for _, schema := range schemaNames {

		tbls, err := parseTables(log, db, schema, filterTables)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		enums, err := parseEnums(log, db, schema, filterTables)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		s := &database.Schema{
			Name:   schema,
			Tables: tbls,
			Enums:  enums,
		}
		res.Schemas = append(res.Schemas, s)
	}
	return res, nil
}

func parseTables(log *log.Logger, db gnorm.DB, schema string, filterTables func(schema, table string) bool) ([]*database.Table, error) {

	log.Println("Tables for schema", schema)

	result := make([]*database.Table, 0)
	tbls, err := tables.Query(db, schema)
	if err != nil {
		return nil, err
	}
	for _, t := range tbls {

		if !filterTables(schema, t.TableName) {
			continue
		}

		cols, err := parseColumns(log, db, t.TableName)
		if err != nil {
			return nil, err
		}

		indices, err := parseIndices(log, db, t.TableName)
		if err != nil {
			return nil, err
		}

		ta := &database.Table{
			Name:    t.TableName,
			Comment: t.TableComment,
			Columns: cols,
			Indexes: indices,
		}
		result = append(result, ta)
	}
	return result, nil
}

func parseColumns(log *log.Logger, db gnorm.DB, table string) ([]*database.Column, error) {

	log.Println("Columns for table", table)

	result := make([]*database.Column, 0)
	cols, err := columns.Query(db, table)
	if err != nil {
		return nil, err
	}
	log.Println("Columns for table ", len(cols))

	for _, c := range cols {

		result = append(result, &database.Column{
			Name:         c.ColumnName,
			Type:         c.DataType,
			IsArray:      false,
			Length:       0,
			UserDefined:  false,
			Nullable:     c.IsNullable,
			HasDefault:   c.ColumnDefault.Valid,
			IsPrimaryKey: c.ColumnKey == 1,
			Ordinal:      c.OrdinalPosition,
			IsForeignKey: false,
			ForeignKey:   nil,
			Orig:         nil,
		})
	}
	return result, nil
}

func parseIndices(log *log.Logger, db gnorm.DB, table string) ([]*database.Index, error) {

	log.Println("Indices for table", table)

	result := make([]*database.Index, 0)
	indices, err := statistics.Query(db, table)
	if err != nil {
		return nil, err
	}
	for _, idx := range indices {

		rcols := make([]*database.Column, 0)
		cols, err := statistics.QueryIndex(db, idx.IndexName)
		if err != nil {
			return nil, err
		}
		for _, col := range cols {
			rcols = append(rcols, &database.Column{
				Name: col.Name,
			})
		}

		result = append(result, &database.Index{
			Name:     idx.IndexName,
			IsUnique: idx.Unique,
			Columns:  rcols,
		})
	}
	return result, nil
}

func parseEnums(log *log.Logger, db gnorm.DB, schema string, filterTables func(schema, table string) bool) ([]*database.Enum, error) {

	log.Println("Eunms for schema NOT SUPPORTED", schema)

	result := make([]*database.Enum, 0)
	return result, nil
}
