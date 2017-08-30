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

Example output for the following schema:

```
CREATE TABLE authors (
  id uuid DEFAULT uuid_generate_v4() NOT NULL primary key,
  name text NOT NULL
);

CREATE INDEX authors_name_idx ON authors(name);

CREATE TYPE book_type AS ENUM (
  'FICTION',
  'NONFICTION'
);

CREATE TABLE books (
  id SERIAL PRIMARY KEY,
  author_id uuid NOT NULL REFERENCES authors(id),
  isbn text NOT NULL UNIQUE,
  booktype book_type NOT NULL,
  title text NOT NULL,
  published timestamptz[] NOT NULL,
  years integer[] NOT NULL,
  pages integer NOT NULL,
  available timestamptz NOT NULL DEFAULT 'NOW()',
  tags varchar[] NOT NULL DEFAULT '{}'
);

CREATE INDEX books_title_idx ON books(author_id, title);
```

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

Enum: public.book_type
+------------+-------+
|    NAME    | VALUE |
+------------+-------+
| FICTION    |     1 |
| NONFICTION |     2 |
+------------+-------+


Table: public.books
+-----------+-----------+--------------------------+---------+--------+-------------+----------+------------+
|  COLUMN   |   TYPE    |          DBTYPE          | ISARRAY | LENGTH | USERDEFINED | NULLABLE | HASDEFAULT |
+-----------+-----------+--------------------------+---------+--------+-------------+----------+------------+
| id        | int       | integer                  | false   |      0 | false       | false    | true       |
| author_id | uuid.UUID | uuid                     | false   |      0 | false       | false    | false      |
| isbn      | string    | text                     | false   |      0 | false       | false    | false      |
| booktype  | BookType  | book_type                | false   |      0 | true        | false    | false      |
| title     | string    | text                     | false   |      0 | false       | false    | false      |
| published | time.Time | timestamptz              | true    |      0 | false       | false    | false      |
| years     | int32     | int4                     | true    |      0 | false       | false    | false      |
| pages     | int       | integer                  | false   |      0 | false       | false    | false      |
| available | time.Time | timestamp with time zone | false   |      0 | false       | false    | true       |
| tags      | string    | varchar                  | true    |      0 | false       | false    | true       |
+-----------+-----------+--------------------------+---------+--------+-------------+----------+------------+


Table: public.schema_version
+---------+------+---------+---------+--------+-------------+----------+------------+
| COLUMN  | TYPE | DBTYPE  | ISARRAY | LENGTH | USERDEFINED | NULLABLE | HASDEFAULT |
+---------+------+---------+---------+--------+-------------+----------+------------+
| version | int  | integer | false   |      0 | false       | false    | false      |
+---------+------+---------+---------+--------+-------------+----------+------------+


Table: public.authors
+--------+-----------+--------+---------+--------+-------------+----------+------------+
| COLUMN |   TYPE    | DBTYPE | ISARRAY | LENGTH | USERDEFINED | NULLABLE | HASDEFAULT |
+--------+-----------+--------+---------+--------+-------------+----------+------------+
| id     | uuid.UUID | uuid   | false   |      0 | false       | false    | true       |
| name   | string    | text   | false   |      0 | false       | false    | false      |
+--------+-----------+--------+---------+--------+-------------+----------+------------+

```
<!-- {{{end}}} -->
