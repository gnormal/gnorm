package environ

import (
	"io"
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
