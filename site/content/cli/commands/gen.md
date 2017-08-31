+++
title= "gen"
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
        Args: []string{"help", "gen"},
    })
    fmt.Println("```")
    os.Exit(x)
}
gocog}}} -->
```
Reads your gnorm.toml file and connects to your database, translating the schema
into in-memory objects.  Then reads your templates and writes files to disk
based on those templates.

Usage:
  gnorm gen [flags]

Flags:
      --config string   relative path to gnorm config file (default "gnorm.toml")
  -h, --help            help for gen
      --verbose         show debugging output
```
<!-- {{{end}}} -->
