package database // import "gnorm.org/gnorm/database"
import (
	"fmt"
	"strings"
)

// Info is the collection of schema info from a database.
type Info struct {
	Schemas []*Schema // the list of schema info
}

// Schema is the information on a single named schema in the database.
type Schema struct {
	Name   string  // the converted name of this schema
	DBName string  // the original name of the schema in the DB
	Tables Tables  // the list of tables in this schema
	Enums  []*Enum // the list of enums in this schema
}

// Enum represents a type that has a set of allowed values.
type Enum struct {
	Schema   string       // the converted name of the schema this enum is in
	DBSchema string       // the original name of the schema in the DB
	Table    string       // (mysql) the converted name of the table this enum is in
	DBTable  string       // (mysql) the original name of the table in the DB
	Name     string       // the converted name of the enum
	DBName   string       // the original name of the enum in the DB
	Values   []*EnumValue // the list of possible values for this enum
}

// EnumValue is one of the named values for an enum.
type EnumValue struct {
	Name   string // the converted name for this enum value
	DBName string // the original label of the enum in the DB
	Value  int    // the value for this enum value (order)
}

// Table contains the definiiton of a database table.
type Table struct {
	Schema   string  // the converted name of the schema this table is in
	DBSchema string  // the original name of the schema in the DB
	Name     string  // the name of the table
	DBName   string  // the original name of the table in the DB
	Columns  Columns // ordered list of columns in this table
}

// Tables is a list of tables in this schema.
type Tables []*Table

// Names returns a list of table Names in this schema.
func (t Tables) Names() Strings {
	names := make([]string, len(t))
	for x := range t {
		names[x] = t[x].Name
	}
	return names
}

// DBNames returns a list of table DBNames in this schema.
func (t Tables) DBNames() Strings {
	names := make([]string, len(t))
	for x := range t {
		names[x] = t[x].DBName
	}
	return names
}

// Column contains data about a column in a table.
type Column struct {
	Name        string      // the converted name of the column
	DBName      string      // the original name of the column in the DB
	Type        string      // the mapped type of the column
	DBType      string      // the original type of the column in the DB
	IsArray     bool        // true if the column type is an array
	Length      int         // non-zero if the type has a length (e.g. varchar[16])
	UserDefined bool        // true if the type is user-defined
	Nullable    bool        // true if the column is not NON NULL
	HasDefault  bool        // true if the column has a default
	Orig        interface{} `yaml:"-"` // the raw database column data
}

// Columns represents the ordered list of columns in a table.
type Columns []*Column

// Names returns the ordered list of column Names in this table.
func (c Columns) Names() Strings {
	names := make(Strings, len(c))
	for x := range c {
		names[x] = c[x].Name
	}
	return names
}

// DBNames returns the ordered list of column DBNames in this table.
func (c Columns) DBNames() Strings {
	names := make(Strings, len(c))
	for x := range c {
		names[x] = c[x].DBName
	}
	return names
}

// Strings is a named type of []string to allow us to put methods on it.
type Strings []string

// Sprintf calls fmt.Sprintf(format, str) for every string in this value and
// returns the results as a new Strings.
func (s Strings) Sprintf(format string) Strings {
	ret := make(Strings, len(s))
	for x := range s {
		ret[x] = fmt.Sprintf(format, s[x])
	}
	return ret
}

// Join returbs strings.Join(s, sep).
func (s Strings) Join(sep string) string {
	return strings.Join([]string(s), sep)
}
