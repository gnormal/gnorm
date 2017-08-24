+++
title= "preview"
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
        Args: []string{"help", "preview"},
    })
    fmt.Println("```")
    os.Exit(x)
}
gocog}}} -->
```
Reads your gnorm.toml file and connects to your database, translating the schema
just as it would be during a full run.  It is then printed out in an
easy-to-read format.

Usage:
  gnorm preview [flags]

Flags:
      --config string   relative path to gnorm config file (default "gnorm.toml")
  -h, --help            help for preview
      --verbose         show debugging output
      --yaml            show output in yaml instead of tabular
```
<!-- {{{end}}} -->