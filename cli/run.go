package cli // import "gnorm.org/gnorm/cli"

import (
	"flag"
	"fmt"
	"io"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"

	"gnorm.org/gnorm/run"
)

var (
	timestamp  = "no timestamp, did you build with make.go?"
	commitHash = "no hash, did you build with make.go?"
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
	if err == errDone {
		return 0
	}
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

const usage = `
usage: gnorm [options]

gnorm parses the gnorm.toml configuration file in the current directory and uses
that to connect to a database. Gnorm reads the schema from the database and
generates code according to your own templates.

options:

  -version 	print version info and exit

`

// errDone indicates the program should exit normally.
var errDone = errors.New("done")

// Parse reads the configuration file and CLI args and returns a gnorm config
// value.
func Parse(env OSEnv) (Config, error) {
	if len(env.Args) == 0 {
		fmt.Fprintln(env.Stderr, usage)
		return Config{}, errDone
	}
	fmt.Printf("%q\n", env.Args)
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	version := fs.Bool("version", false, "")

	if err := fs.Parse(env.Args); err != nil {
		return Config{}, errors.WithStack(err)
	}

	if *version {
		fmt.Fprintf(env.Stderr, `
gnorm built at %s
commit hash: %s
`[1:], timestamp, commitHash)
		return Config{}, errDone
	}

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
