package cli

import (
	"os"
	"strings"

	"github.com/spf13/cobra"

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

// ParseAndRun parses the environment and runs the command.
func ParseAndRun(env environ.Values) int {
	// the return code from the executed command
	var code int
	rootCmd := &cobra.Command{
		Use:   "gnorm",
		Short: "GNORM is Not an ORM, it's a db schema->code generator",
		Long: `
A flexible code generator that turns your DB schema into
runnable code.  See full docs at https://gnorm.org`[1:],
	}
	rootCmd.SetArgs(env.Args)
	rootCmd.SetOutput(env.Stderr)

	rootCmd.AddCommand(previewCmd(env, &code))
	rootCmd.AddCommand(genCmd(env, &code))
	rootCmd.AddCommand(versionCmd(env))
	rootCmd.AddCommand(initCmd(env, &code))
	if err := rootCmd.Execute(); err != nil {
		// cobra outputs the error itself.
		return 2
	}
	return code
}
