package environ // import "gnorm.org/gnorm/environ"

import (
	"io"
	"io/ioutil"
	"log"
)

// Values encapsulates the environment of the OS.
type Values struct {
	Args   []string
	Stderr io.Writer
	Stdout io.Writer
	Stdin  io.Reader
	Env    map[string]string
	Log    *log.Logger
}

// InitLog sets up env.Log to print to stderr if verbose is true, otherwise log
// messages will be discarded.
func (env *Values) InitLog(verbose bool) {
	if verbose {
		env.Log = log.New(env.Stderr, "", 0)
	} else {
		env.Log = log.New(ioutil.Discard, "", 0)
	}
}
