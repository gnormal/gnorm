+++
title= "Command Line"
date= 2017-08-17T13:16:04-04:00
alwaysopen=true
+++

<!-- {{{gocog
package main
import (
    "fmt"
    "os"
    "gnorm.org/gnorm/cli"
)
func main() {
    fmt.Println("```")
    os.Stderr = os.Stdout
    x := cli.Run()
    fmt.Println("```")
    os.Exit(x)
}
gocog}}} -->
```
A flexible code generator that turns your DB schema into
runnable code.  See full docs at https://gnorm.org

Usage:
  gnorm [command]

Available Commands:
  help        Help about any command
  preview     Preview the data that will be sent to your templates
  version     Displays the version of GNORM.

Flags:
  -h, --help   help for gnorm

Use "gnorm [command] --help" for more information about a command.
```
<!-- {{{end}}} -->