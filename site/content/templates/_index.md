+++
title= "Templates"
weight=20
alwaysopen=true
+++

Gnorm uses go templates you write to convert database schemas into code (or
other text) that you can use.  To learn more about using go templates, [read the
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
`"{{.Schema}}/tables/{{.Table}}.html" = "table.tpl"`

This would tell gnorm to use the template at ./table.tpl and output data to
paths formed by the template on the left.  So, if you ran the "users" table from
the "Google" schema, the output file would be "Google/tables/users.html".  Gnorm
will run each table/enum/schema through their respective output targets, so it's
important that the filename template generates unique filenames.

If more than one entry is given, more than one file will be created for each
item.  Thus you could have an entry to generate a db wrapper for your
application, one entry to generate a protobuf definition, and one entry to
generate an HTML docs page.

You can define additional templates for import using TemplatePaths and
TemplateGlobs in your gnorm.toml file. Any files or patterns defined will be
made available via their filename in your table, schema, or enum templates.
Please note that the names of the templates are determined by the filename,
not the full path. This means that if you import two files with the same name
but different path, the last file imported will be the template available.

The following examples are equivalent:

TemplatePaths:
```toml
# gnorm.toml
TemplatePaths = ["templates/insert.tpl", "templates/upsert.tpl"]
```

TemplateGlobs:
```toml
# gnorm.toml
TemplateGlobs = ["templates/*.tpl"]
```

To use in your table, schema, or enum templates
```gotemplate
#table.tpl
{{ template "insert.tpl" .Table }}
{{ template "upsert.tpl" .Table }}
```

