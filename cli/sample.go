package cli

//go:generate gocog -q $GOFILE
// [[[gocog
// package main
// import (
// 	"fmt"
// 	"io"
// 	"os"
// )
// func main() {
// 	fmt.Print("const sample = `")
// 	f, err := os.Open("gnorm.toml")
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}
// 	defer f.Close()
// 	_, err = io.Copy(os.Stdout, f)
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}
// 	fmt.Println("`")
// }
// gocog]]]
const sample = `# ConnStr is the connection string for the database.  Any environment variables
# in this string will be expanded, so for example dbname=$MY_DDB will do the
# right thing.
# MySQL example:
# ConnStr = "root:admin@tcp/"
# Postgres example:
ConnStr = "dbname=mydb host=127.0.0.1 sslmode=disable user=admin"

# DBType holds the type of db you're connecting to.  Possible values are
# "postgres" or "mysql".
DBType = "postgres"

# Schemas holds the names of schemas to generate code for.
Schemas = ["public"]

# PluginDirs a list of paths that will be used for finding plugins.  The list
# will be traversed in order, looking for a specifically named plugin. The first
# plugin that is found will be the one used.
PluginDirs = ["plugins"]

# NameConversion defines how the DBName of tables, schemas, and enums are
# converted into their Name value.  This is a template that may use all the
# regular functions.  The "." value is the DB name of the item. Thus, to make an
# item's Name the same as its DBName, you'd use a template of "{{.}}". To make
# the Name the PascalCase version, you'd use "{{pascal .}}".
NameConversion = "{{.}}"

# IncludeTables is a whitelist of tables to generate data for. Tables not
# in this list will not be included in data geenrated by gnorm. You cannot
# set IncludeTables if ExcludeTables is set.  By default, tables will be
# included in all schemas.  To specify tables for a specific schema only,
# use the schema.tablenmae format.
IncludeTables = []

# ExcludeTables is a blacklist of tables to ignore while generating data.
# All tables in a schema that are not in this list will be used for
# generation. You cannot set ExcludeTables if IncludeTables is set.  By
# default, tables will be excluded from all schemas.  To specify tables for
# a specific schema only, use the schema.tablenmae format.
ExcludeTables = ["xyzzx"]

# PostRun is a command with arguments that is run after each file is generated
# by GNORM.  It is generally used to reformat the file, but it can be for any
# use. Environment variables will be expanded, and the special $GNORMFILE
# environment variable may be used, which will expand to the name of the file
# that was just generated.
# Example to run goimports on each output file:
PostRun = ["echo", "$GNORMFILE"]

# OutputDir is the directory relative to the project root (where the
# gnorm.toml file is located) in which all the generated files are written
# to.
#
# This defaults to the current working directory i.e the directory in which
# gnorm.toml is found.
OutputDir = "gnorm"

# StaticDir is the directory relative to the project root (where the
# gnorm.toml file is located) in which all static files , which are
# intended to be copied to the OutputDir are found.
#
# The directory structure is preserved when copying the files to the
# OutputDir
StaticDir = "static"

# NoOverwriteGlobs is a list of globs
# (https://golang.org/pkg/path/filepath/#Match). If a filename matches a glob
# *and* a file exists with that name, it will not be generated.
NoOverwriteGlobs = ["*.perm.go"]

# TablePaths is a map of output paths to template paths that tells Gnorm how to
# render and output its table info and where to save that output.  Each template
# will be rendered with each table in turn and written out to the given output
# path. If no pairs are specified, tables will not be rendered.  If multiple
# pairs are specified, each one will be generated in turn.
#
# The output path may be a template, in which case the values .Schema and .Table
# may be referenced, containing the name of the current schema and table being
# rendered.  For example, "{{.Schema}}/{{.Table}}/{{.Table}}.go" =
# "tables.gotmpl" would render tables.gotmpl template with data from the the
# "public.users" table to ./public/users/users.go.
[TablePaths]
"{{.Schema}}/tables/{{.Table}}.go" = "testdata/table.tpl"

# SchemaPaths iis a map of output paths to template paths that tells Gnorm how
# to render and output its schema info.  Each template will be rendered with
# each schema in turn and written out to the given output path. If no pairs are
# specified, schemas will not be rendered.  If multiple pairs are specified,
# each one will be generated in turn.
#
# The output path may be a template, in which case the value .Schema may be
# referenced, containing the name of the current schema being rendered. For
# example, "schemas/{{.Schema}}/{{.Schema}}.go" = "schemas.gotmpl" would render
# schemas.gotmpl template with the "public" schema and output to
# ./schemas/public/public.go
[SchemaPaths]
"{{.Schema}}.go" = "testdata/schema.tpl"

# EnumPaths is a is a map of output paths to template paths that tells Gnorm how
# to render and output its enum info.  Each template will be rendered with each
# enum in turn and written out to the given output path. If no pairs are
# specified, enums will not be rendered. If multiple pairs are specified, each
# one will be generated in turn.
#
# The enum path may be a template, in which case the values .Schema and .Enum
# may be referenced, containing the name of the current schema and Enum being
# rendered.  For mysql enums, which are specific to a table, .Table will be
# populated with the table name.  For example,
# "gnorm/{{.Schema}}/enums/{{.Enum}}.go" = "enums.gotmpl" would render the
# enums.gotmpl template with data from the "public.book_type" enum to
# ./gnorm/public/enums/users.go.
[EnumPaths]
"{{.Schema}}/enums/{{.Enum}}.go" = "testdata/enum.tpl"

# TypeMap is a mapping of database type names to replacement type names
# (generally types from your language for deserialization), specifically for
# database columns that are nullable.  In the data sent to your template, this
# is the mapping that translates Column.DBType into Column.Type.  If a DBType is
# not in this mapping, Column.Type will be an empty string.  Note that because
# of the way tables in TOML work, TypeMap and NullableTypeMap must be at the end
# of your configuration file.
# Example for mapping postgres types to Go types:
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
# Example for mapping postgres types to Go types:
[NullableTypeMap]
"timestamp with time zone" = "pq.NullTime"
"text" = "sql.NullString"
"boolean" = "sql.NullBool"
"uuid" = "uuid.NullUUID"
"character varying" = "sql.NullString"
"integer" = "sql.NullInt64"
"numeric" = "sql.NullFloat64"


# Params contains any data you may want to pass to your templates.  This is a
# good way to make templates reusable with different configuration values for
# different situations.  The values in this field will be available in the
# .Params value for all templates.
[Params]
mySpecialValue = "some value"

# TemplateEngine, if specified, describes a command line tool to run to
# render your templates, allowing you to use your preferred templating
# engine.  If not specified, go's text/template will be used to render.
# [TemplateEngine]
    # CommandLine is the command to run to render the template.  You may
    # pass the following variables to the command line -
    # {{.Data}} the name of a .json file containing the gnorm data serialized into json
    # {{.Template}} - the name of the template file being rendered
    # {{.Output}} the target file where output should be written
    # CommandLine = ["yasha", "-v", "{{.Data}}", "{{.Template}}", "-o", "{{.Output}}"]

    # If true, the json data will be sent via stdin to the rendering tool.
    # UseStdin = false

    # If true, the standard output of the tool will be written to the target file.
    # UseStdout = false`
// [[[end]]]
