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
    fmt.Println("```plain\ngnorm preview\n")
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
```plain
gnorm preview

Reads your gnorm.toml file and connects to your database, translating the schema
just as it would be during a full run.  It is then printed out in an
easy-to-read format.  By default it prints out the data in a human-readable
plaintext tabular format.  You may specify a different format using the -format
flag, in which case you can print json, yaml, or types, where types is a list of
all types used by columns in your database.  The latter is useful when setting
up TypeMaps.

Usage:
  gnorm preview [flags]

Flags:
  -c, --config string   relative path to gnorm config file (default "gnorm.toml")
  -f, --format string   Specify output format: tabular, yaml, json, or types (default "tabular")
  -h, --help            help for preview
  -v, --verbose         show debugging output
```
<!-- {{{end}}} -->

Example output for the following schema:

```sql
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

```plain
$ gnorm preview
Schema: Public(public)

Enum: BookType(public.book_type)
+------------+------------+-------+
|    Name    |   DBName   | Value |
+------------+------------+-------+
| Fiction    | FICTION    |     1 |
| Nonfiction | NONFICTION |     2 |
+------------+------------+-------+


Table: Authors(public.authors)
+------+--------+-----------+--------+---------+--------+-------------+----------+------------+
| Name | DBName |   Type    | DBType | IsArray | Length | UserDefined | Nullable | HasDefault |
+------+--------+-----------+--------+---------+--------+-------------+----------+------------+
| ID   | id     | uuid.UUID | uuid   | false   |      0 | false       | false    | true       |
| Name | name   | string    | text   | false   |      0 | false       | false    | false      |
+------+--------+-----------+--------+---------+--------+-------------+----------+------------+


Table: Books(public.books)
+-----------+-----------+-----------+--------------------------+---------+--------+-------------+----------+------------+
|   Name    |  DBName   |   Type    |          DBType          | IsArray | Length | UserDefined | Nullable | HasDefault |
+-----------+-----------+-----------+--------------------------+---------+--------+-------------+----------+------------+
| ID        | id        | int       | integer                  | false   |      0 | false       | false    | true       |
| AuthorID  | author_id | uuid.UUID | uuid                     | false   |      0 | false       | false    | false      |
| Isbn      | isbn      | string    | text                     | false   |      0 | false       | false    | false      |
| Booktype  | booktype  |           | book_type                | false   |      0 | true        | false    | false      |
| Title     | title     | string    | text                     | false   |      0 | false       | false    | false      |
| Published | published | time.Time | timestamptz              | true    |      0 | false       | false    | false      |
| Years     | years     | int32     | int4                     | true    |      0 | false       | false    | false      |
| Pages     | pages     | int       | integer                  | false   |      0 | false       | false    | false      |
| Available | available | time.Time | timestamp with time zone | false   |      0 | false       | false    | true       |
| Tags      | tags      | string    | varchar                  | true    |      0 | false       | false    | true       |
+-----------+-----------+-----------+--------------------------+---------+--------+-------------+----------+------------+


```
