package gnorm // import "gnorm.org/gnorm/database/drivers/sqlite/gnorm"

// Note that this file is *NOT* generated. :)

import (
	"database/sql"
)

// DB is the common interface for database operations.
//
// This should work with database/sql.DB and database/sql.Tx.
type DB interface {
	Exec(string, ...interface{}) (sql.Result, error)
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
}

// WhereClause has a String function should return a properly formatted where
// clause (not including the WHERE) for positional arguments starting at idx.
type WhereClause interface {
	String() string
	Values() []interface{}
}

type comparison string

const (
	compEqual   comparison = " = "
	compGreater comparison = " > "
	compLess    comparison = " < "
	compGTE     comparison = " >= "
	compLTE     comparison = " <= "
	compNE      comparison = " <> "
)

type inClause struct {
	field  string
	values []interface{}
}

func (in inClause) String() string {
	ret := in.field + " in ("
	for x := range in.values {
		if x != 0 {
			ret += ", "
		}
		ret += "?"
	}
	ret += ")"
	return ret
}

func (in inClause) Values() []interface{} {
	return in.values
}

type whereClause struct {
	field string
	comp  comparison
	value interface{}
}

func (w whereClause) String() string {
	ret := w.field + string(w.comp) + "?"
	return ret
}

func (w whereClause) Values() []interface{} {
	return []interface{}{w.value}
}
