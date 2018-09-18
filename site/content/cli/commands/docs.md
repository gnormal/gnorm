+++
title= "docs"
date= 2017-08-17T13:16:04-04:00
description = ""
+++

<!-- {{{gocog
package main
import (
    "fmt"
    "os"
    "gnorm.org/gnorm/cli"
    "gnorm.org/gnorm/environ"
)
func main() {
    fmt.Println("```plain\ngnorm docs\n")
    os.Stderr = os.Stdout
    x := cli.ParseAndRun(environ.Values{
        Stderr: os.Stdout,
        Stdout: os.Stdout,
        Args: []string{"help", "docs"},
    })
    fmt.Println("```")
    os.Exit(x)
}
gocog}}} -->
```plain
gnorm docs

Starts a web server running at localhost:8080 that serves docs for this version
of Gnorm.

Usage:
  gnorm docs [flags]

Flags:
  -h, --help   help for docs
```
<!-- {{{end}}} -->