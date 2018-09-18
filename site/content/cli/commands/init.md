+++
title= "init"
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
    fmt.Println("```plain\ngnorm init\n")

    os.Stderr = os.Stdout
    x := cli.ParseAndRun(environ.Values{
        Stderr: os.Stdout,
        Stdout: os.Stdout,
        Args: []string{"help", "init"},
    })
    fmt.Println("```")
    os.Exit(x)
}
gocog}}} -->
```plain
gnorm init

Creates a default gnorm.toml and the various template files needed to run GNORM.

Usage:
  gnorm init [flags]

Flags:
  -h, --help   help for init
```
<!-- {{{end}}} -->
