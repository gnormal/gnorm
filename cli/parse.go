package cli // import "gnorm.org/gnorm/cli"

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"gnorm.org/gnorm/environ"
	"gnorm.org/gnorm/run"
)

var (
	timestamp  = "no timestamp, did you build with make.go?"
	commitHash = "no hash, did you build with make.go?"
)

// ParseAndRun parses the environment and runs the command.
func ParseAndRun(env environ.Values) int {
	// the error from the executed command
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
	rootCmd.AddCommand(versionCmd(env))
	if err := rootCmd.Execute(); err != nil {
		// cobra outputs the error itself.
		return 2
	}
	return code
}

// parse reads the configuration file and returns a gnorm config value.
func parse(log *log.Logger, file string) (Config, error) {
	c := Config{}
	m, err := toml.DecodeFile(file, &c)
	if err != nil {
		return Config{}, errors.WithMessage(err, "error parsing config file")
	}
	undec := m.Undecoded()
	if len(undec) > 0 {
		log.Println("Warning: unknown values present in config file:", undec)
	}
	return c, nil
}

func previewCmd(env environ.Values, code *int) *cobra.Command {
	var cfgFile string
	var useYaml bool
	var verbose bool
	preview := &cobra.Command{
		Use:   "preview",
		Short: "Preview the data that will be sent to your templates",
		Long: `
Reads your gnorm.toml file and connects to your database, translating the schema
just as it would be during a full run.  It is then printed out in an
easy-to-read format.`[1:],
		Run: func(cmd *cobra.Command, args []string) {
			if verbose {
				env.Log = log.New(env.Stderr, "", 0)
			} else {
				env.Log = log.New(ioutil.Discard, "", 0)
			}

			cfg, err := parse(env.Log, cfgFile)
			if err != nil {
				fmt.Fprintln(env.Stderr, err)
				*code = 2
				return
			}
			if err := run.Preview(env, run.Config(cfg), useYaml, verbose); err != nil {
				env.Log.Println(err)
				*code = 1
			}
		},
	}
	preview.Flags().StringVar(&cfgFile, "config", "gnorm.toml", "relative path to gnorm config file")
	preview.Flags().BoolVar(&useYaml, "yaml", false, "show output in yaml instead of tabular")
	preview.Flags().BoolVar(&verbose, "verbose", false, "show debugging output")
	return preview
}

func versionCmd(env environ.Values) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Displays the version of GNORM.",
		Long: `
Shows the build date and commit hash used to build this binary.`[1:],
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(env.Stdout, "built at: %s\ncommit hash: %s", timestamp, commitHash)
		},
	}
}
