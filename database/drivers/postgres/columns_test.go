package postgres

import (
	"database/sql"
	"log"
	"testing"

	"gnorm.org/gnorm/database/drivers/postgres/pg"
)

// These values are actual values created by reading postgres 9.6.3.
var (
	// uuid
	AuthorIDCol = &pg.Column{
		TableCatalog:         sql.NullString{String: "gnorm-db", Valid: true},
		TableSchema:          sql.NullString{String: "public", Valid: true},
		TableName:            sql.NullString{String: "books", Valid: true},
		ColumnName:           sql.NullString{String: "author_id", Valid: true},
		OrdinalPosition:      sql.NullInt64{Int64: 2, Valid: true},
		IsNullable:           false,
		DataType:             sql.NullString{String: "uuid", Valid: true},
		UdtCatalog:           sql.NullString{String: "gnorm-db", Valid: true},
		UdtSchema:            sql.NullString{String: "pg_catalog", Valid: true},
		UdtName:              sql.NullString{String: "uuid", Valid: true},
		DtdIdentifier:        sql.NullString{String: "2", Valid: true},
		IsGenerated:          sql.NullString{String: "NEVER", Valid: true},
		GenerationExpression: sql.NullString{String: "", Valid: false},
		IsUpdatable:          true,
	}

	// character(32)
	ISBNCol = &pg.Column{
		TableCatalog:           sql.NullString{String: "gnorm-db", Valid: true},
		TableSchema:            sql.NullString{String: "public", Valid: true},
		TableName:              sql.NullString{String: "books", Valid: true},
		ColumnName:             sql.NullString{String: "isbn", Valid: true},
		OrdinalPosition:        sql.NullInt64{Int64: 3, Valid: true},
		IsNullable:             false,
		DataType:               sql.NullString{String: "character", Valid: true},
		CharacterMaximumLength: sql.NullInt64{Int64: 32, Valid: true},
		CharacterOctetLength:   sql.NullInt64{Int64: 128, Valid: true},
		UdtCatalog:             sql.NullString{String: "gnorm-db", Valid: true},
		UdtSchema:              sql.NullString{String: "pg_catalog", Valid: true},
		UdtName:                sql.NullString{String: "bpchar", Valid: true},
		DtdIdentifier:          sql.NullString{String: "3", Valid: true},
		IsGenerated:            sql.NullString{String: "NEVER", Valid: true},
		IsUpdatable:            true,
	}

	// autoincrement id
	BooksIDCol = &pg.Column{
		TableCatalog:          sql.NullString{String: "gnorm-db", Valid: true},
		TableSchema:           sql.NullString{String: "public", Valid: true},
		TableName:             sql.NullString{String: "books", Valid: true},
		ColumnName:            sql.NullString{String: "id", Valid: true},
		OrdinalPosition:       sql.NullInt64{Int64: 1, Valid: true},
		ColumnDefault:         sql.NullString{String: "nextval('books_id_seq'::regclass)", Valid: true},
		IsNullable:            false,
		DataType:              sql.NullString{String: "integer", Valid: true},
		NumericPrecision:      sql.NullInt64{Int64: 32, Valid: true},
		NumericPrecisionRadix: sql.NullInt64{Int64: 2, Valid: true},
		NumericScale:          sql.NullInt64{Int64: 0, Valid: true},
		UdtCatalog:            sql.NullString{String: "gnorm-db", Valid: true},
		UdtSchema:             sql.NullString{String: "pg_catalog", Valid: true},
		UdtName:               sql.NullString{String: "int4", Valid: true},
		DtdIdentifier:         sql.NullString{String: "1", Valid: true},
		IsGenerated:           sql.NullString{String: "NEVER", Valid: true},
		IsUpdatable:           true,
	}

	// Timestamp with timezone with a default of now()
	AvailableCol = &pg.Column{
		TableCatalog:      sql.NullString{String: "gnorm-db", Valid: true},
		TableSchema:       sql.NullString{String: "public", Valid: true},
		TableName:         sql.NullString{String: "books", Valid: true},
		ColumnName:        sql.NullString{String: "available", Valid: true},
		OrdinalPosition:   sql.NullInt64{Int64: 9, Valid: true},
		ColumnDefault:     sql.NullString{String: "'2017-09-04 20:04:33.854571-04'::timestamp with time zone", Valid: true},
		IsNullable:        false,
		DataType:          sql.NullString{String: "timestamp with time zone", Valid: true},
		DatetimePrecision: sql.NullInt64{Int64: 6, Valid: true},
		UdtCatalog:        sql.NullString{String: "gnorm-db", Valid: true},
		UdtSchema:         sql.NullString{String: "pg_catalog", Valid: true},
		UdtName:           sql.NullString{String: "timestamptz", Valid: true},
		DtdIdentifier:     sql.NullString{String: "9", Valid: true},
		IsGenerated:       sql.NullString{String: "NEVER", Valid: true},
		IsUpdatable:       true,
	}

	// nullable text
	SummaryCol = &pg.Column{
		TableCatalog:         sql.NullString{String: "gnorm-db", Valid: true},
		TableSchema:          sql.NullString{String: "public", Valid: true},
		TableName:            sql.NullString{String: "books", Valid: true},
		ColumnName:           sql.NullString{String: "summary", Valid: true},
		OrdinalPosition:      sql.NullInt64{Int64: 9, Valid: true},
		IsNullable:           true,
		DataType:             sql.NullString{String: "text", Valid: true},
		CharacterOctetLength: sql.NullInt64{Int64: 1073741824, Valid: true},
		UdtCatalog:           sql.NullString{String: "gnorm-db", Valid: true},
		UdtSchema:            sql.NullString{String: "pg_catalog", Valid: true},
		UdtName:              sql.NullString{String: "text", Valid: true},
		DtdIdentifier:        sql.NullString{String: "9", Valid: true},
		IsGenerated:          sql.NullString{String: "NEVER", Valid: true},
		IsUpdatable:          true,
	}

	// int4 array
	YearsCol = &pg.Column{
		TableCatalog:    sql.NullString{String: "gnorm-db", Valid: true},
		TableSchema:     sql.NullString{String: "public", Valid: true},
		TableName:       sql.NullString{String: "books", Valid: true},
		ColumnName:      sql.NullString{String: "years", Valid: true},
		OrdinalPosition: sql.NullInt64{Int64: 7, Valid: true},
		IsNullable:      false,
		DataType:        sql.NullString{String: "ARRAY", Valid: true},
		UdtCatalog:      sql.NullString{String: "gnorm-db", Valid: true},
		UdtSchema:       sql.NullString{String: "pg_catalog", Valid: true},
		UdtName:         sql.NullString{String: "_int4", Valid: true},
		DtdIdentifier:   sql.NullString{String: "7", Valid: true},
		IsGenerated:     sql.NullString{String: "NEVER", Valid: true},
		IsUpdatable:     true,
	}

	// user defined enum type
	BookTypeCol = &pg.Column{
		TableCatalog:    sql.NullString{String: "gnorm-db", Valid: true},
		TableSchema:     sql.NullString{String: "public", Valid: true},
		TableName:       sql.NullString{String: "books", Valid: true},
		ColumnName:      sql.NullString{String: "booktype", Valid: true},
		OrdinalPosition: sql.NullInt64{Int64: 4, Valid: true},
		IsNullable:      false,
		DataType:        sql.NullString{String: "USER-DEFINED", Valid: true},
		UdtCatalog:      sql.NullString{String: "gnorm-db", Valid: true},
		UdtSchema:       sql.NullString{String: "public", Valid: true},
		UdtName:         sql.NullString{String: "book_type", Valid: true},
		DtdIdentifier:   sql.NullString{String: "4", Valid: true},
		IdentityCycle:   false, IsGenerated: sql.NullString{String: "NEVER", Valid: true},
		IsUpdatable: true,
	}
)

type testLog struct {
	t *testing.T
}

func (l testLog) Write(b []byte) (n int, err error) {
	l.t.Log(string(b))
	return len(b), nil
}

func tLog(t *testing.T) *log.Logger {
	return log.New(testLog{t}, "", 0)
}

func TestLength(t *testing.T) {
	col := toDBColumn(SummaryCol, tLog(t))
	if col.Length != 0 {
		t.Errorf("Text column length should be zero but was %v", col.Length)
	}

	col = toDBColumn(ISBNCol, tLog(t))
	if col.Length != 32 {
		t.Errorf("character(32) column length should be 32 but was %v", col.Length)
	}

	col = toDBColumn(AuthorIDCol, tLog(t))
	if col.Length != 0 {
		t.Errorf("uuid column length should be 0 but was %v", col.Length)
	}
}

func TestArray(t *testing.T) {
	col := toDBColumn(SummaryCol, tLog(t))
	if col.IsArray {
		t.Error("Text column length should not be labelled IsArray, but is.")
	}

	col = toDBColumn(YearsCol, tLog(t))
	if !col.IsArray {
		t.Error("int4 array not marked as IsArray")
	}
}

func TestUserDefined(t *testing.T) {
	col := toDBColumn(SummaryCol, tLog(t))
	if col.UserDefined {
		t.Error("Text column should not be labelled UserDefined, but is.")
	}

	col = toDBColumn(BookTypeCol, tLog(t))
	if !col.UserDefined {
		t.Error("user defined enum not marked as UserDefined")
	}
	if col.DBType != BookTypeCol.UdtName.String {
		t.Errorf("Expected column to have UdtName %q as DBType, but instead got %s", BookTypeCol.UdtName.String, col.DBType)
	}
}

func TestNullable(t *testing.T) {
	col := toDBColumn(ISBNCol, tLog(t))
	if col.Nullable {
		t.Errorf("Column should not be nullable but was set as Nullable = true.")
	}

	col = toDBColumn(AuthorIDCol, tLog(t))
	if col.Nullable {
		t.Errorf("Column should not be nullable but was set as Nullable = true.")
	}

	col = toDBColumn(BookTypeCol, tLog(t))
	if col.Nullable {
		t.Errorf("Column should not be nullable but was set as Nullable = true.")
	}

	col = toDBColumn(SummaryCol, tLog(t))
	if !col.Nullable {
		t.Errorf("Column should be nullable but was set as not Nullable = true.")
	}
}

func TestHasDefault(t *testing.T) {
	col := toDBColumn(ISBNCol, tLog(t))
	if col.HasDefault {
		t.Errorf("Column has not default but was set as HasDefault = true.")
	}

	col = toDBColumn(AuthorIDCol, tLog(t))
	if col.HasDefault {
		t.Errorf("Column has not default but was set as HasDefault = true.")
	}

	col = toDBColumn(BookTypeCol, tLog(t))
	if col.HasDefault {
		t.Errorf("Column has not default but was set as HasDefault = true.")
	}

	col = toDBColumn(AvailableCol, tLog(t))
	if !col.HasDefault {
		t.Errorf("Column has default but was set as HasDefault = false.")
	}

	col = toDBColumn(BooksIDCol, tLog(t))
	if !col.HasDefault {
		t.Errorf("Column has default but was set as HasDefault = false.")
	}
}

func TestDBType(t *testing.T) {
	col := toDBColumn(SummaryCol, tLog(t))
	if col.DBType != SummaryCol.DataType.String {
		t.Errorf("Column should have name %q but is %q.", SummaryCol.DataType.String, col.DBType)
	}

	// User-defined columns are different, their data type is USER-DEFINED which
	// is less than useful, so we copy the UdtName into the column Type.
	col = toDBColumn(BookTypeCol, tLog(t))
	if col.DBType != BookTypeCol.UdtName.String {
		t.Errorf("Expected column to have UdtName %q as DBType, but instead got %s", BookTypeCol.UdtName.String, col.DBType)
	}
}
