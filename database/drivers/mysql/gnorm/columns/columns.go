// Code generated by gnorm, DO NOT EDIT!

package columns

import (
	"database/sql"

	"gnorm.org/gnorm/database/drivers/mysql/gnorm"
)

// Row represents a row from 'COLUMNS'.
type Row struct {
	TableCatalog           string         // TABLE_CATALOG
	TableSchema            string         // TABLE_SCHEMA
	TableName              string         // TABLE_NAME
	ColumnName             string         // COLUMN_NAME
	OrdinalPosition        int64          // ORDINAL_POSITION
	ColumnDefault          sql.NullString // COLUMN_DEFAULT
	IsNullable             string         // IS_NULLABLE
	DataType               string         // DATA_TYPE
	CharacterMaximumLength sql.NullInt64  // CHARACTER_MAXIMUM_LENGTH
	CharacterOctetLength   sql.NullInt64  // CHARACTER_OCTET_LENGTH
	NumericPrecision       sql.NullInt64  // NUMERIC_PRECISION
	NumericScale           sql.NullInt64  // NUMERIC_SCALE
	DatetimePrecision      sql.NullInt64  // DATETIME_PRECISION
	CharacterSetName       sql.NullString // CHARACTER_SET_NAME
	CollationName          sql.NullString // COLLATION_NAME
	ColumnType             string         // COLUMN_TYPE
	ColumnKey              string         // COLUMN_KEY
	Extra                  string         // EXTRA
	Privileges             string         // PRIVILEGES
	ColumnComment          string         // COLUMN_COMMENT
	GenerationExpression   string         // GENERATION_EXPRESSION
}

// Field values for every column in Columns.
var (
	TableCatalogCol           gnorm.StringField        = "TABLE_CATALOG"
	TableSchemaCol            gnorm.StringField        = "TABLE_SCHEMA"
	TableNameCol              gnorm.StringField        = "TABLE_NAME"
	ColumnNameCol             gnorm.StringField        = "COLUMN_NAME"
	OrdinalPositionCol        gnorm.Int64Field         = "ORDINAL_POSITION"
	ColumnDefaultCol          gnorm.SqlNullStringField = "COLUMN_DEFAULT"
	IsNullableCol             gnorm.StringField        = "IS_NULLABLE"
	DataTypeCol               gnorm.StringField        = "DATA_TYPE"
	CharacterMaximumLengthCol gnorm.SqlNullInt64Field  = "CHARACTER_MAXIMUM_LENGTH"
	CharacterOctetLengthCol   gnorm.SqlNullInt64Field  = "CHARACTER_OCTET_LENGTH"
	NumericPrecisionCol       gnorm.SqlNullInt64Field  = "NUMERIC_PRECISION"
	NumericScaleCol           gnorm.SqlNullInt64Field  = "NUMERIC_SCALE"
	DatetimePrecisionCol      gnorm.SqlNullInt64Field  = "DATETIME_PRECISION"
	CharacterSetNameCol       gnorm.SqlNullStringField = "CHARACTER_SET_NAME"
	CollationNameCol          gnorm.SqlNullStringField = "COLLATION_NAME"
	ColumnTypeCol             gnorm.StringField        = "COLUMN_TYPE"
	ColumnKeyCol              gnorm.StringField        = "COLUMN_KEY"
	ExtraCol                  gnorm.StringField        = "EXTRA"
	PrivilegesCol             gnorm.StringField        = "PRIVILEGES"
	ColumnCommentCol          gnorm.StringField        = "COLUMN_COMMENT"
	GenerationExpressionCol   gnorm.StringField        = "GENERATION_EXPRESSION"
)

// Query retrieves rows from 'COLUMNS' as a slice of Row.
func Query(db gnorm.DB, where gnorm.WhereClause) ([]*Row, error) {
	const origsqlstr = `SELECT 
		TABLE_CATALOG, TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ORDINAL_POSITION, COLUMN_DEFAULT, IS_NULLABLE, DATA_TYPE, CHARACTER_MAXIMUM_LENGTH, CHARACTER_OCTET_LENGTH, NUMERIC_PRECISION, NUMERIC_SCALE, DATETIME_PRECISION, CHARACTER_SET_NAME, COLLATION_NAME, COLUMN_TYPE, COLUMN_KEY, EXTRA, PRIVILEGES, COLUMN_COMMENT, GENERATION_EXPRESSION
		FROM information_schema.COLUMNS WHERE (`

	sqlstr := origsqlstr + where.String() + ") "

	var vals []*Row
	q, err := db.Query(sqlstr, where.Values()...)
	if err != nil {
		return nil, err
	}
	for q.Next() {
		r := Row{}

		err = q.Scan(&r.TableCatalog, &r.TableSchema, &r.TableName, &r.ColumnName, &r.OrdinalPosition, &r.ColumnDefault, &r.IsNullable, &r.DataType, &r.CharacterMaximumLength, &r.CharacterOctetLength, &r.NumericPrecision, &r.NumericScale, &r.DatetimePrecision, &r.CharacterSetName, &r.CollationName, &r.ColumnType, &r.ColumnKey, &r.Extra, &r.Privileges, &r.ColumnComment, &r.GenerationExpression)
		if err != nil {
			return nil, err
		}

		vals = append(vals, &r)
	}
	return vals, nil
}
