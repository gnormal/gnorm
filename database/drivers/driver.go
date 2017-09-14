package drivers

import (
	"log"

	"gnorm.org/gnorm/database"
)

// Driver defines the base interface for databases that are supported by gnorm
type Driver interface {
	Name() string
	Parse(log *log.Logger, conn string, schemaNames []string, filterTables func(schema, table string) bool) (*database.Info, error)
}

// Templates is an interface for gnorm templates. The templates vary depending
// on the active driver.
type Templates interface {
	TplNames() []string
	Tpl(name string) ([]byte, error)
}
