+++
title= "version"
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
    fmt.Println("```")
    os.Stderr = os.Stdout
    x := cli.ParseAndRun(environ.Values{
        Stderr: os.Stdout,
        Stdout: os.Stdout,
        Args: []string{"help", "version"},
    })
    fmt.Println("```")
    os.Exit(x)
}
gocog}}} -->
```
Shows the build date and commit hash used to build this binary.

Usage:
  gnorm version [flags]

Flags:
  -h, --help   help for version
```
<!-- {{{end}}} -->