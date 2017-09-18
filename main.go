package main // import "gnorm.org/gnorm"

import (
	"os"

	"gnorm.org/gnorm/cli"
)

// This space used for gocog commands run across the website.
// If we use gocog on .go files, put the go generate command in that file.

//go:generate go get github.com/natefinch/gocog
//go:generate gocog ./site/content/gnorm.md
//go:generate gocog ./site/content/cli/_index.md --startmark={{{ --endmark=}}}
//go:generate gocog ./site/content/cli/commands/preview.md --startmark={{{ --endmark=}}}
//go:generate gocog ./site/content/cli/commands/version.md --startmark={{{ --endmark=}}}
//go:generate gocog ./site/content/cli/commands/init.md --startmark={{{ --endmark=}}}
//go:generate gocog ./site/content/cli/commands/gen.md --startmark={{{ --endmark=}}}
//go:generate gocog ./site/content/templates/functions.md --startmark={{{ --endmark=}}}
//go:generate gocog ./site/content/cli/configuration.md --startmark={{{ --endmark=}}}

func main() {
	os.Exit(cli.Run())
}
