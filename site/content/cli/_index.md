+++
title= "Command Line"
date= 2017-08-17T13:16:04-04:00
description = ""
draft= false
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

usage: gnorm [command]

gnorm parses the gnorm.toml configuration file in the current directory and uses
that to connect to a database. Gnorm reads the schema from the database and
generates code according to your own templates.

commands

  version 	print version info and exit
  run 		generate code according to config in gnorm.toml
  

```
<!-- {{{end}}} -->