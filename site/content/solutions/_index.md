+++
title= "Solutions"
weight=40
alwaysopen=true
+++

Gnorm Solutions are pre-packaged sets of configuration files, templates, and
static files that work together to produce code useful for a specific purpose.
They are much like themes for static site generators such as Hugo or Jekyll.

Solutions offer an easy way to distribute a ready-to-use gnorm configuration for
such things as generating Go or Python code to wrap a database, or generating
protobuf definitions, or even a full REST API for the DB. 


## How To Create a Solution

A solution is simply a gnorm.toml file, and any corresponding template files and
static files used by that gnorm.toml.  It is a good idea to put these files into
a source control repo, since they will likely evolve over time.

## Default Layout

A good default layout is to put the gnorm.toml file in the root of the repo, put
templates in a subdirectory called "templates", and put static files in a
subdirectory called "static".  Ensure that the template paths in your gnorm.toml
reference the template directory (e.g. "templates/foo.tpl").  Ensure the
StaticDir property in the gnorm.toml is set to "static".

It's good practice to set OutputDir to something other than the current
directory, so that users of the solution can delete the entire generated
directory without worrying about losing anything.  A good default for this
directory name is simply "gnorm".  Of course, users of your solution may
customize this name to whatever they want.

## Supporting Multiple Databases

Gnorm solutions are often database-specific.  How you interact with a mysql
database is generally different than how you interact with a postgres database,
for example. 

To support multiple databases with the same solution, you will almost always
need database-specific gnorm.toml files, and often need database-specific
template files.  The best way to provide a solution for multiple databases is to
create a separate gnorm.toml for each database (called gnorm-postgres.toml etc),
and then provide separate template directories for each configuration file, such
as templates/postgres/ and templates/mysql, and separate static directories (if
needed).  The each db-specific config can then reference the db-specific
templates and static files.

To use such a solution, the end user can then either rename the appropriate toml
file to gnorm.toml, or specify it on the command line with the -c option for
`gnorm gen` and `gnorm preview`.