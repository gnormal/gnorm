+++
title= "Template Data"
weight=1
alwaysopen=true
+++

The data passed to the templates is defined below.

<!-- {{{gocog
package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	fmt.Println("```")
	for _, s := range []string{"Schema", "Table", "Tables", "Column", "Columns", "Enum", "EnumValue"} {
		c := exec.Command("go", "doc", "gnorm.org/gnorm/database."+s)
		b, err := c.CombinedOutput()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(string(b))
	}
	fmt.Println("```")
}
gocog}}} -->
```
type Schema struct {
	Name   string  // the converted name of this schema
	DBName string  // the original name of the schema in the DB
	Tables Tables  // the list of tables in this schema
	Enums  []*Enum // the list of enums in this schema
}
    Schema is the information on a single named schema in the database.


type Table struct {
	Schema   string  // the converted name of the schema this table is in
	DBSchema string  // the original name of the schema in the DB
	Name     string  // the name of the table
	DBName   string  // the original name of the table in the DB
	Columns  Columns // ordered list of columns in this table
}
    Table contains the definiiton of a database table.


type Tables []*Table
    Tables is a list of tables in this schema.


func (t Tables) DBNames() []string
func (t Tables) Names() []string

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
    Column contains data about a column in a table.


type Columns []*Column
    Columns represents the ordered list of columns in a table.


func (c Columns) DBNames() []string
func (c Columns) Names() []string

type Enum struct {
	Schema   string       // the converted name of the schema this enum is in
	DBSchema string       // the original name of the schema in the DB
	Table    string       // (mysql) the converted name of the table this enum is in
	DBTable  string       // (mysql) the original name of the table in the DB
	Name     string       // the converted name of the enum
	DBName   string       // the original name of the enum in the DB
	Values   []*EnumValue // the list of possible values for this enum
}
    Enum represents a type that has a set of allowed values.


type EnumValue struct {
	Name   string // the converted name for this enum value
	DBName string // the original label of the enum in the DB
	Value  int    // the value for this enum value (order)
}
    EnumValue is one of the named values for an enum.


```
<!-- {{{end}}} -->
