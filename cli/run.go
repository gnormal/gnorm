package cli // import "gnorm.org/gnorm/cli"

import (
	"fmt"
	"io"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"

	"gnorm.org/gnorm/run"
)

// OSEnv encapsulates the environment.
type OSEnv struct {
	Args   []string
	Stderr io.Writer
	Stdout io.Writer
	Stdin  io.Reader
	Env    map[string]string
}

// ParseAndRun parses the environment to create a run.Command and runs it.  It
// returns the code that should be used for os.Exit.
func ParseAndRun(env OSEnv) int {
	cfg, err := Parse(env)
	if err != nil {
		fmt.Fprintln(env.Stderr, err)
		return 2
	}

	if err := run.Run(run.Config(cfg)); err != nil {
		fmt.Fprintln(env.Stderr, err)
		return 1
	}
	return 0
}

// Parse reads the configuration file and CLI args and returns a gnorm config
// value.
func Parse(env OSEnv) (Config, error) {
	c := Config{}
	m, err := toml.DecodeFile("gnorm.toml", &c)
	if err != nil {
		return Config{}, errors.WithMessage(err, "error parsing config file")
	}
	undec := m.Undecoded()
	if len(undec) > 0 {
		fmt.Fprintf(env.Stderr, "Warning: unknown values present in config file: %v\n", undec)
	}
	return c, nil
}
