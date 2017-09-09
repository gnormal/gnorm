package main // import "gnorm.org/gnorm"

import (
	"os"

	"gnorm.org/gnorm/cli"
)

//go:generate go get github.com/natefinch/gocog
//go:generate gocog -q ./site/content/gnorm.md
//go:generate gocog -q ./site/content/cli/_index.md --startmark={{{ --endmark=}}}
//go:generate gocog -q ./site/content/cli/commands/preview.md --startmark={{{ --endmark=}}}
//go:generate gocog -q ./site/content/cli/commands/version.md --startmark={{{ --endmark=}}}
//go:generate gocog -q ./site/content/cli/commands/init.md --startmark={{{ --endmark=}}}
//go:generate gocog -q ./site/content/cli/commands/gen.md --startmark={{{ --endmark=}}}
//go:generate gocog -q ./site/content/templates/data.md --startmark={{{ --endmark=}}}
//go:generate gocog -q ./site/content/templates/functions.md --startmark={{{ --endmark=}}}
//go:generate gocog ./site/content/cli/configuration.md --startmark={{{ --endmark=}}}

func main() {
	os.Exit(cli.Run())
}
