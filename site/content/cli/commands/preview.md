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

Example output:

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
    fmt.Println("$ gnorm preview")
    os.Stderr = os.Stdout
    x := cli.ParseAndRun(environ.Values{
        Stderr: os.Stdout,
        Stdout: os.Stdout,
        Args: []string{"preview"},
    })
    fmt.Println("```")
    os.Exit(x)
}
gocog}}} -->
```
$ gnorm preview
Schema: public

Table: public.authors
+--------+-----------+--------+---------+--------+-------------+----------+
| COLUMN |   TYPE    | DBTYPE | ISARRAY | LENGTH | USERDEFINED | NULLABLE |
+--------+-----------+--------+---------+--------+-------------+----------+
| id     | uuid.UUID | uuid   | false   |      0 | false       | false    |
| name   | string    | text   | false   |      0 | false       | false    |
+--------+-----------+--------+---------+--------+-------------+----------+


Table: public.books
+-----------+-----------+--------------------------+---------+--------+-------------+----------+
|  COLUMN   |   TYPE    |          DBTYPE          | ISARRAY | LENGTH | USERDEFINED | NULLABLE |
+-----------+-----------+--------------------------+---------+--------+-------------+----------+
| id        | int       | integer                  | false   |      0 | false       | false    |
| author_id | uuid.UUID | uuid                     | false   |      0 | false       | false    |
| isbn      | string    | text                     | false   |      0 | false       | false    |
| booktype  | BookType  | book_type                | false   |      0 | true        | false    |
| title     | string    | text                     | false   |      0 | false       | false    |
| published | time.Time | timestamptz              | true    |      0 | false       | false    |
| years     | int32     | int4                     | true    |      0 | false       | false    |
| pages     | int       | integer                  | false   |      0 | false       | false    |
| available | time.Time | timestamp with time zone | false   |      0 | false       | false    |
| tags      | string    | varchar                  | true    |      0 | false       | false    |
+-----------+-----------+--------------------------+---------+--------+-------------+----------+


Table: public.schema_version
+---------+------+---------+---------+--------+-------------+----------+
| COLUMN  | TYPE | DBTYPE  | ISARRAY | LENGTH | USERDEFINED | NULLABLE |
+---------+------+---------+---------+--------+-------------+----------+
| version | int  | integer | false   |      0 | false       | false    |
+---------+------+---------+---------+--------+-------------+----------+

```
<!-- {{{end}}} -->
