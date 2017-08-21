package cli // import "gnorm.org/gnorm/cli"

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"

	"gnorm.org/gnorm/environ"
	"gnorm.org/gnorm/run"
)

var (
	timestamp  = "no timestamp, did you build with make.go?"
	commitHash = "no hash, did you build with make.go?"
)

// ParseAndRun parses the environment to create a run.Command and runs it.  It
// returns the code that should be used for os.Exit.
func ParseAndRun(env environ.Values) int {
	cfg, err := Parse(env)
	if err == errDone {
		return 0
	}
	if err != nil {
		fmt.Fprintln(env.Stderr, err)
		return 2
	}

	if err := run.Run(env, run.Config(cfg)); err != nil {
		fmt.Fprintln(env.Stderr, err)
		return 1
	}
	return 0
}

const usage = `
usage: gnorm [command]

gnorm parses the gnorm.toml configuration file in the current directory and uses
that to connect to a database. Gnorm reads the schema from the database and
generates code according to your own templates.

commands

  version 	print version info and exit
  run 		generate code according to config in gnorm.toml
  
`

// errDone indicates the program should exit normally.
var errDone = errors.New("done")

// Parse reads the configuration file and CLI args and returns a gnorm config
// value.
func Parse(env environ.Values) (Config, error) {
	if len(env.Args) == 0 {
		fmt.Fprintln(env.Stderr, usage)
		return Config{}, errDone
	}
	if len(env.Args) > 1 {
		return Config{}, errors.New("too many arguments")
	}

	switch env.Args[0] {
	case "version":
		fmt.Fprintf(env.Stderr, `
gnorm built at %s
commit hash: %s
`[1:], timestamp, commitHash)
		return Config{}, errDone
	case "run":
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
	return Config{}, errors.New("unknown command: " + env.Args[0])
}
