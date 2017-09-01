package main // import "gnorm.org/gnorm"

import (
	"os"

	"gnorm.org/gnorm/cli"
)

//go:generate go get github.com/natefinch/gocog
//go:generate gocog ./site/content/gnorm.md
//go:generate gocog ./site/content/cli/_index.md --startmark={{{ --endmark=}}}
//go:generate gocog ./site/content/cli/commands/preview.md --startmark={{{ --endmark=}}}
//go:generate gocog ./site/content/cli/commands/version.md --startmark={{{ --endmark=}}}
//go:generate gocog ./site/content/cli/commands/gen.md --startmark={{{ --endmark=}}}
//go:generate gocog ./site/content/templates/data.md --startmark={{{ --endmark=}}}
//go:generate gocog ./site/content/templates/functions.md --startmark={{{ --endmark=}}}

func main() {
	os.Exit(cli.Run())
}
