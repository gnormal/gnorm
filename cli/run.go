package cli

import (
	"log"
	"os"
	"strings"

	"gnorm.org/gnorm/environ"
)

// Run captures the OS environment and passes it to ParseAndRun.  It returns
// the code that the executable should exit with.
func Run() int {
	env := environ.Values{
		Args:   make([]string, len(os.Args)-1),
		Stderr: os.Stderr,
		Stdout: os.Stdout,
		Stdin:  os.Stdin,
		Env:    getenv(os.Environ()),
		Log:    log.New(os.Stderr, "", 0),
	}
	copy(env.Args, os.Args[1:])
	return ParseAndRun(env)
}

func getenv(env []string) map[string]string {
	ret := make(map[string]string, len(env))
	for _, s := range env {
		parts := strings.SplitN(s, "=", 2)
		if len(parts) != 2 {
			panic("invalid environment variable set: " + s)
		}
		ret[parts[0]] = parts[1]
	}
	return ret
}
