package database // import "gnorm.org/gnorm/database"

// Info is the collection of schema info from a database.
type Info struct {
	Schemas []Schema
}

// Schema is the information on a single named schema in the database.
type Schema struct {
	Name   string
	Tables []Table
	Enums  []Enum
}

// Enum represents a type that has a set of allowed values.
type Enum struct {
	Name   string
	Values []EnumValue
}

// EnumValue is one of the named values for an enum.
type EnumValue struct {
	Name  string
	Value int
}

// Table contains the definiiton of a database table.
type Table struct {
	Name    string
	Columns []Column
}

// LengthType defines what the Column Length means.
type LengthType string

const (
	// MaxLength means the Column Length is a maximum size.
	MaxLength LengthType = "max"

	// FixedLength means the Column Length is a fixed length.
	FixedLength LengthType = "fixed"
)

// Column contains data about a column in a table.
type Column struct {
	Name        string
	Type        string
	DBType      string
	IsArray     bool
	Length      int
	UserDefined bool
	Nullable    bool
	Orig        interface{}
}
