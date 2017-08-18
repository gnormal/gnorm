+++
title= "GNORM"
date= 2017-08-17T13:16:04-04:00
description = ""
draft= false
+++


## About

Gnorm is a datebase-first code generator that generates boilerplate that matches your database's schema.  

Gnorm uses templates you control, so that you can make the output look exactly
how you want it to look.  It can be used to generate type-safe database queries
that are faster than a traditional ORM.  It can also be used to create a REST or
RPC API that exposes the data in your database.

## Templates

Gnorm reads your database schema, then runs the resulting data through templates you can customize in any way you like.  Included in the repo are a set of templates produces Go structs and functions using the stdlib's database/sql package as a thorough example and usable database layer for most go projects.

However, your templates may generate whatever code or text files you wish, based on your templates.

## Configuration

Configuring gnorm is as easy as creating a simple [TOML](https://github.com/toml-lang/toml) file with a few configuration values.  Gnorm takes care of the rest.

## Code

https://github.com/gnormal/gnorm