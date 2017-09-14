package drivers

import (
	"log"
	"sync"

	"github.com/pkg/errors"
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

var regDriver = &sync.Map{}

// Get returns a registered driverby name.
func Get(name string) (Driver, error) {
	if d, ok := regDriver.Load(name); ok {
		return d.(Driver), nil
	}
	return nil, errors.Errorf("unknown database type: %v", name)
}

// Register registers driver d.
func Register(d Driver) {
	regDriver.Store(d.Name(), d)
}
