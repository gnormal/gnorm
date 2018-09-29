+++
title= "Getting Started"
weight = 8
+++

## Set Up Your DB

Gnorm expects that you know how you want your data to be stored and organized in
your database.  Do that first, using whatever tools you're most comfortable
with.  Gnorm doesn't care how the DB gets set up, that's up to you.  If you're
porting a project with an existing database that you want to use with Gnorm,
great! There's nothing to do for this step.

## Get Gnorm

The easiest way is to download the latest
[release](https://github.com/gnormal/gnorm/releases).  There's no dependencies,
and all major platforms are supported.

If you have a Go programming environment set up, you can `go get
gnorm.org/gnorm`. Note that if you install this way, `gnorm version` will not
output any build info.  If you'd like the build info to be correct, run `mage build`,
which will compile the binary with the current git information.

## Set up your config

run `gnorm init` to get a default gnorm.toml file and template files in the
current directory.  You need to at least set the `ConnStr`, `DBType`, and
`Schemas` values.  `ConnStr` is the connection string for your database,
`DBType` tells Gnorm what kind of database it's working against, and `Schemas`
tells it what DB schemas to query.

If you want to use a template rendering engine other than Go's text/template,
fill out the TemplateEngine section of the configuration.

## Test your config 

Now, to test your configuration, run `gnorm preview`. This will query your
database and spit out all the information Gnorm knows about your database,
including schema names, table names, column info, custom types, etc. in a nice
tabular format.

## Tweak the output

There are a few more configuration values that are useful for most projects.  

`NameConversion` is a template that lets you set a standard string conversion
for all database names. For example, if your database is all snake_case but you
want those names to be PascalCase for your code, you can set NameConversion to
do that in one place, rather than having to do it manually in your templates
everywhere.  Check out the information on template
[functions](/templates/functions) to see possibilities.  And
don't worry, the original name in the database will still be available in the
[data](/templates/data)

`TypeMap` and `NullableTypeMap` are also conversion mechanisms for mapping
column types into programming language types.  Simply map the database typename
on the left to a programming type name on the right.  These may not be required
for all applications, if so, you may ignore them.

If you changed anything here, re-run `gnorm preview` to see the results. The
`NameConversion` field will change the Name of things, and the type maps will
change the Type of a column.

The final config value to look at is `OutputDir`.  This is the base directory
where all your generated file will be created.  If it is not set, it defaults to
the same directory where Gnorm is run, but it is often adviseable to use a
subdirectory to keep generated code separate from the code that generates it.
gnorm init creates a gnorm.toml file with `OutputDir = "gnorm"` which means your
generated code will be created in a subdirectory of the directory where gnorm is
run.  It's generally best to use a relative directory for OutputDir, as using an
absolute directory may not work on other people's machines.

## Let's generate! 

Gnorm init gives you very basic templates that do not really produce output that
would be useful in any real application. To produce something you can use, you
have to write templates or use a pre-made gnorm solution to format the data.
If you've ever used a static site generator like [Hugo](https://gohugo.io),
solutions are like themes.

## Learn More About Gnorm's Features

The best place to start to learn more about what Gnorm can do is to read about
the [configuration values](/cli/configuration/) in gnorm.toml that control how
gnorm behaves.