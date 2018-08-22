// Package data supplies the data for gnorm templates
package data

import (
	"fmt"
	"sort"
)

// This is all the data passed to templates.

// DBData is all the data about a database that we know.
type DBData struct {
	Schemas       []*Schema
	SchemasByName map[string]*Schema `yaml:"-" json:"-"` // dbname to schema
}

// SchemaData is the data passed to schema templates.
type SchemaData struct {
	Schema *Schema
	DB     *DBData
	Config ConfigData
	Params map[string]interface{}
}

// TableData is the data passed to table templates.
type TableData struct {
	Table  *Table
	DB     *DBData
	Config ConfigData
	Params map[string]interface{}
}

// EnumData is the data passed to enum templates.
type EnumData struct {
	Enum   *Enum
	DB     *DBData
	Config ConfigData
	Params map[string]interface{}
}

// Schema is the data about a DB schema.
type Schema struct {
	Name         string            // the converted name of the schema
	DBName       string            // the original name of the schema in the DB
	Tables       Tables            // the list of tables in this schema
	Enums        Enums             // the list of enums in this schema
	TablesByName map[string]*Table `yaml:"-" json:"-"` // dbnames to tables
}

// Table is the data about a DB Table.
type Table struct {
	Name           string                 // the converted name of the table
	DBName         string                 // the original name of the table in the DB
	Comment        string                 // the comment attached to the table
	Schema         *Schema                `yaml:"-" json:"-"` // the schema this table is in
	Columns        Columns                // Database columns
	ColumnsByName  map[string]*Column     `yaml:"-" json:"-"` // dbname to column
	PrimaryKeys    Columns                // Primary Key Columns
	Indexes        Indexes                // Table indexes
	IndexesByName  map[string]*Index      `yaml:"-" json:"-"` // indexname to index
	ForeignKeys    ForeignKeys            // Foreign Keys
	ForeignKeyRefs ForeignKeys            // Foreign Keys referencing this table
	FKByName       map[string]*ForeignKey `yaml:"-" json:"-"` // Foreign Keys by foreign key name
	FKRefsByName   map[string]*ForeignKey `yaml:"-" json:"-"` // Foreign Keys referencing this table by foreign key name
}

// HasPrimaryKey returns true if Table has one or more primary keys.
func (t *Table) HasPrimaryKey() bool {
	return len(t.PrimaryKeys) > 0
}

// Returns true if Table has one or more foreign keys
func (t *Table) HasForeignKeys() bool {
	return len(t.ForeignKeys) > 0
}

// Returns true if one or more foreign keys reference Table
func (t *Table) HasForeignKeyRefs() bool {
	return len(t.ForeignKeyRefs) > 0
}

// Column is the data about a DB column of a table.
type Column struct {
	Table              *Table                       `yaml:"-" json:"-"` // the table this column is in
	Name               string                       // the converted name of the column
	DBName             string                       // the original name of the column in the DB
	Type               string                       // the converted name of the type
	DBType             string                       // the original type of the column in the DB
	IsArray            bool                         // true if the column type is an array
	Length             int                          // non-zero if the type has a length (e.g. varchar[16])
	UserDefined        bool                         // true if the type is user-defined
	Nullable           bool                         // true if the column is not NON NULL
	HasDefault         bool                         // true if the column has a default
	Comment            string                       // the comment attached to the column
	IsPrimaryKey       bool                         // true if the column is a primary key
	IsFK               bool                         // true if the column is a foreign key
	HasFKRef           bool                         // true if the column is referenced by a foreign key
	FKColumn           *ForeignKeyColumn            // foreign key column definition
	FKColumnRefs       ForeignKeyColumns            // all foreign key columns referencing this column
	FKColumnRefsByName map[string]*ForeignKeyColumn `yaml:"-" json:"-"` // all foreign key columns referencing this column by foreign key name
	Orig               interface{}                  `yaml:"-" json:"-"` // the raw database column data
}

// ForeignKey contains the
type ForeignKey struct {
	DBName         string            // the original name of the foreign key constraint in the db
	Name           string            // the converted name of the foreign key constraint
	TableDBName    string            // the original name of the table in the db
	RefTableDBName string            // the original name of the foreign table in the db
	Table          *Table            `yaml:"-" json:"-"` // the foreign key table
	RefTable       *Table            `yaml:"-" json:"-"` // the foreign key foreign table
	FKColumns      ForeignKeyColumns // all foreign key columns belonging to the foreign key
}

// ForeignKeyColumn contains the definition of a database foreign key at the kcolumn level
type ForeignKeyColumn struct {
	DBName          string  // the original name of the foreign key constraint in the db
	ColumnDBName    string  // the original name of the column in the db
	RefColumnDBName string  // the original name of the foreign column in the db
	Column          *Column `yaml:"-" json:"-"` // the foreign key column
	RefColumn       *Column `yaml:"-" json:"-"` // the referenced column
}

// Index is the data about a table index.
type Index struct {
	Name     string  // the converted name of the index
	DBName   string  // dbname of the index
	IsUnique bool    // true if index is unique
	Columns  Columns // columns used in the index
}

// Enum represents a type that has a set of allowed values.
type Enum struct {
	Name   string       // the converted name of the enum
	DBName string       // the original name of the enum in the DB
	Schema *Schema      `yaml:"-" json:"-"` // the schema the enum is in
	Table  *Table       `yaml:"-" json:"-"` // (mysql) the table this enum is part of
	Values []*EnumValue // the list of possible values for this enum
}

// EnumValue is one of the named values for an enum.
type EnumValue struct {
	Name   string // the converted label of the enum
	DBName string // the original label of the enum in the DB
	Value  int    // the value for this enum value (order)
}

// ConfigData holds the portion of the config that will be available to
// templates.  Note that Params are added to the data at a higher level.
type ConfigData struct {
	// ConnStr is the connection string for the database.  Environment variables
	// in $FOO form will be expanded.
	ConnStr string

	// The type of DB you're connecting to.  Currently the possible values are
	// "postgres" or "mysql".
	DBType string

	// Schemas holds the names of schemas to generate code for.
	Schemas []string

	// IncludeTables is a map of schema names to table names. It is whitelist of
	// tables to generate data for. Tables not in this list will not be included
	// in data generated by gnorm. You cannot set IncludeTables if ExcludeTables
	// is set.
	IncludeTables map[string][]string

	// ExcludeTables is a map of schema names to table names.  It is a blacklist
	// of tables to ignore while generating data. All tables in a schema that
	// are not in this list will be used for generation. You cannot set
	// ExcludeTables if IncludeTables is set.
	ExcludeTables map[string][]string

	// PostRun is a command with arguments that is run after each file is
	// generated by GNORM.  It is generally used to reformat the file, but it
	// can be for any use. Environment variables will be expanded, and the
	// special $GNORMFILE environment variable may be used, which will expand to
	// the name of the file that was just generated.
	PostRun []string

	// TypeMap is a mapping of database type names to replacement type names
	// (generally types from your language for deserialization).  Types not in
	// this list will remain in their database form.  In the data sent to your
	// template, this is the Column.Type, and the original type is in
	// Column.OrigType.  Note that because of the way tables in TOML work,
	// TypeMap and NullableTypeMap must be at the end of your configuration
	// file.
	TypeMap map[string]string

	// NullableTypeMap is a mapping of database type names to replacement type
	// names (generally types from your language for deserialization)
	// specifically for database columns that are nullable.  Types not in this
	// list will remain in their database form.  In the data sent to your
	// template, this is the Column.Type, and the original type is in
	// Column.OrigType.   Note that because of the way tables in TOML work,
	// TypeMap and NullableTypeMap must be at the end of your configuration
	// file.
	NullableTypeMap map[string]string

	// PluginDirs a set of absolute/relative  paths that will be used for
	// plugin lookup.
	PluginDirs []string

	// OutputDir is the directory relative to the project root (where the
	// gnorm.toml file is located) in which all the generated files are written
	// to.
	//
	// This defaults to the current working directory i.e the directory in which
	// gnorm.toml is found.
	OutputDir string

	// StaticDir is the directory relative to the project root (where the
	// gnorm.toml file is located) in which all static files , which are
	// intended to be copied to the OutputDir are found.
	//
	// The directory structure is preserved when copying the files to the
	// OutputDir
	StaticDir string

	// NoOverwriteGlobs is a list of globs
	// (https://golang.org/pkg/path/filepath/#Match). If a filename matches a glob
	// *and* a file exists with that name, it will not be generated.
	NoOverwriteGlobs []string
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

// Except returns a copy of the Strings with the given values removed.
func (s Strings) Except(excludes []string) Strings {
	ret := make(Strings, 0, len(s))
	for x := range s {
		if !contains(excludes, s[x]) {
			ret = append(ret, s[x])
		}
	}
	return ret
}

// Sorted returns a sorted copy of the Strings.
func (s Strings) Sorted() Strings {
	ret := make(Strings, len(s))
	copy(ret, s)
	sort.Strings([]string(ret))
	return ret
}

func contains(list []string, s string) bool {
	for x := range list {
		if s == list[x] {
			return true
		}
	}
	return false
}

// Foreign keys represents a list of ForeignKey
type ForeignKeys []*ForeignKey

// DBNames returns the list of db foreign key names
func (fk ForeignKeys) DBNames() Strings {
	names := make(Strings, len(fk))
	for x := range fk {
		names[x] = fk[x].DBName
	}
	return names
}

// Names returns the list of converted foreign key names
func (fk ForeignKeys) Names() Strings {
	names := make(Strings, len(fk))
	for x := range fk {
		names[x] = fk[x].Name
	}
	return names
}

// ForeignKeyColumns represents a list of ForeignKeyColumn
type ForeignKeyColumns []*ForeignKeyColumn

// DBNames returns the list of db foreign key names
func (fkc ForeignKeyColumns) DBNames() Strings {
	names := make(Strings, len(fkc))
	for x := range fkc {
		names[x] = fkc[x].DBName
	}
	return names
}

// ColumnDBNames returns the list of column database names.
func (fkc ForeignKeyColumns) ColumnDBNames() Strings {
	names := make(Strings, len(fkc))
	for x := range fkc {
		names[x] = fkc[x].ColumnDBName
	}
	return names
}

// RefColumnDBNames returns the list of foreign column database names.
func (fkc ForeignKeyColumns) RefColumnDBNames() Strings {
	names := make(Strings, len(fkc))
	for x := range fkc {
		names[x] = fkc[x].RefColumnDBName
	}
	return names
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

// Enums represents all the enums in a schema.
type Enums []*Enum

// Names returns the list of enum Names in this schema.
func (c Enums) Names() Strings {
	names := make(Strings, len(c))
	for x := range c {
		names[x] = c[x].Name
	}
	return names
}

// DBNames returns the list of enum DBNames in this schema.
func (c Enums) DBNames() Strings {
	names := make(Strings, len(c))
	for x := range c {
		names[x] = c[x].DBName
	}
	return names
}

// Indexes represents all the indexes on a table.
type Indexes []*Index

// DBNames returns the list of index DBNames in this table.
func (i Indexes) DBNames() Strings {
	names := make(Strings, len(i))
	for x := range i {
		names[x] = i[x].DBName
	}
	return names
}

// Names returns the list of index Name in this table.
func (i Indexes) Names() Strings {
	names := make(Strings, len(i))
	for x := range i {
		names[x] = i[x].Name
	}
	return names
}
