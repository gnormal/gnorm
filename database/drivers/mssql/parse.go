package mssql // import "gnorm.org/gnorm/database/drivers/mssql"

import (
	"database/sql"
	"gnorm.org/gnorm/database/drivers/mssql/gnorm"
	"gnorm.org/gnorm/database/drivers/mssql/gnorm/parameters"
	"gnorm.org/gnorm/database/drivers/mssql/gnorm/procs"
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

		procs, err := parseProcs(log, db, schema, filterTables)
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
			Procs:  procs,
			Enums:  enums,
		}
		res.Schemas = append(res.Schemas, s)
	}
	return res, nil
}

func parseProcs(log *log.Logger, db gnorm.DB, schema string, filterTables func(schema, table string) bool) ([]*database.Proc, error) {

	log.Println("Procs for schema", schema)

	result := make([]*database.Proc, 0)
	tbls, err := procs.Query(db, schema)
	if err != nil {
		return nil, err
	}
	for _, t := range tbls {

		if !filterTables(schema, t.ProcName) {
			continue
		}

		parameters, err := parseParameters(log, db, schema, t.ProcName)
		if err != nil {
			return nil, err
		}

		ta := &database.Proc{
			Name:       t.ProcName,
			Comment:    t.ProcComment,
			Parameters: parameters,
		}
		result = append(result, ta)
	}
	return result, nil
}

func parseParameters(log *log.Logger, db gnorm.DB, schema string, proc string) ([]*database.Parameter, error) {

	log.Println("Columns for proc", proc)

	result := make([]*database.Parameter, 0)
	params, err := parameters.Query(db, schema, proc)
	if err != nil {
		return nil, err
	}
	log.Println("Parameters for proc ", len(params))

	for _, p := range params {

		param := &database.Parameter{
			Name:        p.ParameterName,
			Type:        p.DataType,
			IsArray:     false,
			Length:      0,
			UserDefined: false,
			Nullable:    p.IsNullable == "YES",
			HasDefault:  p.ParameterDefault.Valid,
			Comment:     "",
			Ordinal:     int64(p.OrdinalPosition),
			Orig:        nil,
		}

		result = append(result, param)

	}
	return result, nil
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

		cols, err := parseColumns(log, db, schema, t.TableName)
		if err != nil {
			return nil, err
		}

		indices, err := parseIndices(log, db, schema, t.TableName)
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

func parseColumns(log *log.Logger, db gnorm.DB, schema string, table string) ([]*database.Column, error) {

	log.Println("Columns for table", table)

	result := make([]*database.Column, 0)
	cols, err := columns.Query(db, schema, table)
	if err != nil {
		return nil, err
	}
	log.Println("Columns for table ", len(cols))

	foreigns, err := columns.QueryForeign(db, schema, table)
	if err != nil {
		return nil, err
	}
	log.Println("Foreigns for table ", len(foreigns))

	fks := make(map[string]*columns.ForeignRow)
	for _, fk := range foreigns {
		fks[fk.FkColumnName] = fk
	}

	for _, c := range cols {

		fk, isForeignKey := fks[c.ColumnName]

		var colfk *database.ForeignKey = nil

		if isForeignKey {
			colfk = &database.ForeignKey{
				SchemaName:               schema,
				TableName:                fk.PkTableName,
				ColumnName:               fk.PkColumnName,
				Name:                     fk.FkName,
				UniqueConstraintPosition: 0,
				ForeignTableName:         fk.PkTableName,
				ForeignColumnName:        fk.PkColumnName,
			}
		}

		col := &database.Column{
			Name:         c.ColumnName,
			Type:         c.DataType,
			IsArray:      false,
			Length:       0,
			UserDefined:  false,
			Nullable:     c.IsNullable == "YES",
			HasDefault:   c.ColumnDefault.Valid,
			Comment:      "",
			IsPrimaryKey: c.Pk,
			Ordinal:      int64(c.OrdinalPosition),
			IsForeignKey: isForeignKey,
			ForeignKey:   colfk,
			Orig:         nil,
		}

		result = append(result, col)

	}
	return result, nil
}

func parseIndices(log *log.Logger, db gnorm.DB, schema string, table string) ([]*database.Index, error) {

	log.Println("Indices for table", table)

	result := make([]*database.Index, 0)
	indices, err := statistics.Query(db, table)
	if err != nil {
		return nil, err
	}
	for _, idx := range indices {

		rcols := make([]*database.Column, 0)
		cols, err := statistics.QueryIndex(db, idx.Name)
		if err != nil {
			return nil, err
		}
		for _, col := range cols {
			rcols = append(rcols, &database.Column{
				Name: col.ColumnName,
			})
		}

		result = append(result, &database.Index{
			Name:     idx.Name,
			IsUnique: idx.IsUnique,
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
