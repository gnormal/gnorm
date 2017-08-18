package main // import "gnorm.org/gnorm"

import (
	"os"

	"gnorm.org/gnorm/cli"
)

//go:generate go get github.com/natefinch/gocog
//go:generate gocog ./site/content/code/gnorm.md

func main() {
	os.Exit(cli.Run())
}
