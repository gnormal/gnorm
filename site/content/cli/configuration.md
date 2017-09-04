+++
title = "Configuration"
date = 2017-08-23T13:16:04-04:00
+++

Gnorm is configured using a configuration file written in
[TOML](https://github.com/toml-lang/toml).  The file must be called gnorm.toml
and must live in the directory where you call gnorm.

### example configuration file
```
# ConnStr is the connection string for the database.
ConnStr = "dbname=mydb host=127.0.0.1 sslmode=disable user=admin"

# DBType holds the type of db you're connecting to.  Possible values are
# "postgres" or "mysql". 
DBType = "postgres"

# Schemas holds the names of schemas to generate code for.
Schemas = ["public", "foobar"]

# TemplateDir contains the relative path to the directory where gnorm
# expects to find templates to render.  The default is the current
# directory where gnorm is running.
TemplateDir = "templates"

# PostRun is a command with arguments that is run after each file is generated
# by GNORM.  It is generally used to reformat the file, but it can be for any
# use. Environment variables will be expanded, and the special $GNORMFILE
# environment variable may be used, which will expand to the name of the file
# that was just generated.
PostRun = ["goimports", "-w", "$GNORMFILE"]

# TablePath is a relative path for tables to be rendered.  The table.tpl
# template will be rendered with each table in turn. If the path is empty,
# tables will not be rendered this way (though you could render them via the
# schemas template).
#
# The table path may be a template, in which case the values .Schema and .Table
# may be referenced, containing the name of the current schema and table being
# rendered.  For example, "gnorm/{{.Schema}}/tables/{{.Table}}/{{.Table}}.go"
# would render the "public.users" table to ./gnorm/public/tables/users/users.go.
TablePath = "schemas/{{.Schema}}/{{.Table}}/{{.Table}}.go"

# SchemaPath is a relative path for schemas to be rendered.  The schema.tpl
# template will be rendered with each schema in turn. If the path is empty,
# schemas will not be rendered.
#
# The schema path may be a template, in which case the value .Schema may be
# referenced, containing the name of the current schema being rendered. For
# example, "schemas/{{.Schema}}/{{.Schema}}.go" would render the "public"
# schema to ./schemas/public/public.go
SchemaPath = "schemas/{{.Schema}}/{{.Schema}}.go"

# EnumPath is a relative path for enums to be rendered.  The enum.tpl template
# will be rendered with each enum in turn. If the path is empty, enums will not
# be rendered this way (thought you could render them via the schemas template).
#
# The enum path may be a template, in which case the values .Schema and .Enum
# may be referenced, containing the name of the current schema and Enum being
# rendered.  For example, "gnorm/{{.Schema}}/enums/{{.Enum}}.go" would render
# the "public.book_type" enum to ./gnorm/public/enums/users.go.
EnumPath =  "gnorm/{{.Schema}}/enums/{{.Enum}}.go"

# NameConversion defines how the DBName of tables, schemas, and enums are
# converted into their Name value.  This is a template that may use all the
# regular functions.  The "." value is the DB name of the item. Thus, to make an
# item's Name the same as its DBName, you'd use a template of "{{.}}". To make
# the Name the PascalCase version, you'd use "{{pascal .}}".
NameConversion = "{{pascal .}}"

# TypeMap is a mapping of database type names to replacement type names
# (generally types from your language for deserialization), specifically for
# database columns that are nullable.  In the data sent to your template, this
# is the mapping that translates Column.DBType into Column.Type.  If a DBType is
# not in this mapping, Column.Type will be an empty string.  Note that because
# of the way tables in TOML work, TypeMap and NullableTypeMap must be at the end
# of your configuration file.
[TypeMap]
"timestamp with time zone" = "time.Time"
"text" = "string"
"boolean" = "bool"
"uuid" = "uuid.UUID"
"character varying" = "string"
"integer" = "int"
"numeric" = "float64"

# NullableTypeMap is a mapping of database type names to replacement type names
# (generally types from your language for deserialization), specifically for
# database columns that are nullable.  In the data sent to your template, this
# is the mapping that translates Column.DBType into Column.Type.  If a DBType is
# not in this mapping, Column.Type will be an empty string.  Note that because
# of the way tables in TOML work, TypeMap and NullableTypeMap must be at the end
# of your configuration file.
[NullableTypeMap]
"timestamp with time zone" = "pq.NullTime"
"text" = "sql.NullString"
"boolean" = "sql.NullBool"
"uuid" = "uuid.NullUUID"
"character varying" = "sql.NullString"
"integer" = "sql.NullInt64"
"numeric" = "sql.NullFloat64"

```