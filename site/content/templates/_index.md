+++
title= "Templates"
weight=20
alwaysopen=true
+++

Gnorm uses go templates you write to convert database schemas into code (or
other text) that you can use.

By default, gnorm uses Go's text/template.  However, you can configure gnorm
to use whatever templating engine you want, so long as there is a command
line tool that lets you render files.  To use a different templating engine,
fill out the `TemplateEngine` section in the [configuration](/cli/configuration)


To learn more about using go templates, [read the
documentation](https://golang.org/pkg/text/template/).

There are three templates that gnorm uses to generate code: 

- Table templates
- Schema templates
- Enum templates

The location of these templates is defined in your gnorm.toml file in
`SchemaPaths`, `TablePaths`, and `EnumPaths` values.  Each value has 0 or more sub
values in the following format:

"output filename template" = "contents template filename"

The right hand side is the filename of a template you've written, which will
format the schema information into code, html, whatever. 

The left hand side determines the filename where those contents will be written.

For example, an entry in TablePaths might looks like this:
`"{{.Schema}}/tables/{{.Table}}.html" = "table.tpl"

This would tell gnorm to use the template at ./table.tpl and output data to
paths formed by the template on the left.  So, if you ran the "users" table from
the "Google" schema, the output file would be "Google/tables/users.html".  Gnorm
will run each table/enum/schema through their respective output targets, so it's
important that the filename template generates unique filenames.

If more than one entry is given, more than one file will be created for each
item.  Thus you could have an entry to generate a db wrapper for your
application, one entry to generate a protobuf definition, and one entry to
generate an HTML docs page.