+++
title= "Command Line"
weight=10
alwaysopen=true
+++

### Gnorm is...

<!-- {{{gocog
package main
import (
    "fmt"
    "os"
    "gnorm.org/gnorm/cli"
)
func main() {
    fmt.Println("```plain")
    os.Stderr = os.Stdout
    x := cli.Run()
    fmt.Println("```")
    os.Exit(x)
}
gocog}}} -->
```plain
A flexible code generator that turns your DB schema into
runnable code.  See full docs at https://gnorm.org

Usage:
  gnorm [command]

Available Commands:
  docs        Runs a local webserver serving gnorm documentation.
  gen         Generate code from DB schema
  help        Help about any command
  init        Generates the files needed to run GNORM.
  preview     Preview the data that will be sent to your templates
  version     Displays the version of GNORM.

Flags:
  -h, --help   help for gnorm

Use "gnorm [command] --help" for more information about a command.
```
<!-- {{{end}}} -->