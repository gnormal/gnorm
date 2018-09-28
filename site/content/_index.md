+++
title= "Home"
date= 2017-08-17T13:16:04-04:00
weight=20
+++
# GNORM

Gnorm converts your database's schema into in-memory data structures which you
can then feed into your own templates to produce code or documentation or
whatever. 

Gnorm is written in Go but can be used to generate any kind of textual output -
ruby, python, protobufs, html, javascript, etc.

Gnorm uses templates you control, so that you can make the output look exactly
how you want it to look.  

## Gnorm Works with Any Template Engine and Any Language

Gnorm can use any templating language, and can generate code or other textual
output for any language or system.  

It can also be used to create a REST or RPC API that exposes the data in your
database.  It can be used to generate type-safe database queries that are faster than a
traditional ORM.  It could even be used to create a yaml configuration for some tool
that depends on the database.

## Getting Started

Check out how to [get started](/getting-started)

## Demo

{{< youtube IPLpZj2SLWw >}}

## Templates

Gnorm reads your database schema, then runs the resulting data through templates
you can customize in any way you like to produce code or documentation.

## Configuration

Configuring gnorm is as easy as creating a simple
[TOML](https://github.com/toml-lang/toml) file with a few configuration values.
Gnorm takes care of the rest.  See details about configuration
[here](/cli/configuration).

## Database Support

GNORM supports Postgres and MySQL.  Further database support is planned for the
future.  Contributions more than welcome.  Check out database/drivers to get an
idea of what is involved (it's not that much). 

## Downloads

https://github.com/gnormal/gnorm/releases

## Feature Support

Databases support a large variety of features, GNORM will add support for more
features as time allows, and as needed.  If you need a feature, feel free to
make an issue (and preferably also a PR).

## What's Wrong With ORMs?

Gnorm believes that the code should not generate the database schema, but the
database schema should generate the code.  Schema is declarative, code is
imperative, so there's an inherent mismatch of doing code-first.  You end up
having to write your schema in special tags or structures purely to generate a
schema... when there's already a langauge to generate a database schema - SQL.
Why create a new one?  By creating your database schema exactly the way you want
it and then generating code to access it, you get the full power of the
database, and no hidden costs for what an ORM is doing behind your back.

## Discussion

If you have questions about Gnorm or want to hack on it, meet the devs on the #gnorm 
channel of [gopher slack](https://gophers.slack.com/).

## Code

https://github.com/gnormal/gnorm