package database // import "gnorm.org/gnorm/database"
import (
	"log"
)

// Info is the collection of schema info from a database.
type Info struct {
	Schemas []*Schema // the list of schema info
}

// Schema is the information on a single named schema in the database.
type Schema struct {
	Name   string   // the original name of the schema in the DB
	Tables []*Table // the list of tables in this schema
	Enums  []*Enum  // the list of enums in this schema
}

// Enum represents a type that has a set of allowed values.
type Enum struct {
	Table  string       // (mysql) the original name of the table in the DB
	Name   string       // the original name of the enum in the DB
	Values []*EnumValue // the list of possible values for this enum
}

// EnumValue is one of the named values for an enum.
type EnumValue struct {
	Name  string // the original label of the enum in the DB
	Value int    // the value for this enum value (order)
}

// Table contains the definition of a database table.
type Table struct {
	Name    string    // the original name of the table in the DB
	Comment string    // the comment attached to the table
	Columns []*Column // ordered list of columns in this table
	Indexes []*Index  // list of indexes in this table
}

// Index contains the definition of a database index.
type Index struct {
	Name    string    // name of the index in the database
	Columns []*Column // list of columns in this index
}

// PrimaryKey contains the definition of a database primary key.
type PrimaryKey struct {
	SchemaName string // the original name of the schema in the db
	TableName  string // the original name of the table in the db
	ColumnName string // the original name of the column in the db
	Name       string // the original name of the key constraint in the db
}

// Foreign Key contains the definition of a database foreign key
type ForeignKey struct {
	SchemaName               string // the original name of the schema in the db
	TableName                string // the original name of the table in the db
	ColumnName               string // the original name of the column in the db
	Name                     string // the original name of the foreign key constraint in the db
	UniqueConstraintPosition int    // the position of the unique constraint in the db
	ForeignTableName         string // the original name of the table in the db for the referenced table
	ForeignColumnName        string // the original name of the column in the db for the referenced column
}

// Column contains data about a column in a table.
type Column struct {
	Name         string      // the original name of the column in the DB
	Type         string      // the original type of the column in the DB
	IsArray      bool        // true if the column type is an array
	Length       int         // non-zero if the type has a length (e.g. varchar[16])
	UserDefined  bool        // true if the type is user-defined
	Nullable     bool        // true if the column is not NON NULL
	HasDefault   bool        // true if the column has a default
	Comment      string      // the comment attached to the column
	IsPrimaryKey bool        // true if the column is a primary key
	IsForeignKey bool        // true if the column is a foreign key
	ForeignKey   *ForeignKey // foreign key database definition
	Orig         interface{} // the raw database column data
}

// Driver defines the base interface for databases that are supported by gnorm
type Driver interface {
	Parse(log *log.Logger, conn string, schemaNames []string, filterTables func(schema, table string) bool) (*Info, error)
}
