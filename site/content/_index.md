+++
title= "Home"
date= 2017-08-17T13:16:04-04:00
description = ""
draft= false
+++
# GNORM is Not an ORM

Gnorm is a database-first code generator that generates boilerplate that matches
your database's schema.

Gnorm uses templates you control, so that you can make the output look exactly
how you want it to look.  It can be used to generate type-safe database queries
that are faster than a traditional ORM.  It can also be used to create a REST or
RPC API that exposes the data in your database.

## Templates

Gnorm reads your database schema, then runs the resulting data through templates
you can customize in any way you like to produce code or documentation.

## Configuration

Configuring gnorm is as easy as creating a simple
[TOML](https://github.com/toml-lang/toml) file with a few configuration values.
Gnorm takes care of the rest.  See details about configuration
[here](/cli/configuration).

## Feature Support

Databases support a large variety of features, GNORM will add support for more
features as time allows, and as needed.  If you need a feature, feel free to
make an issue (and preferably also a PR).

## Database Support

Right now, GNORM only supports Postgres, however, adding more database support
is fairly easy and is something we'll be working on once we have the feature set
stabilized a little.  Contributions more than welcome.  Check out
database/drivers/postgres to get an idea of what is involved (it's not that
much). 

## Code

https://github.com/gnormal/gnorm