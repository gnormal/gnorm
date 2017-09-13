+++
title= "Gnorm vs ..."
date= 2017-08-17T13:16:04-04:00
description = ""
weight=9
+++


## [GORM](https://github.com/jinzhu/gorm)

GORM is a Go ORM.  

### Differences:

1. GORM is code-first.  It generates the database from your code, using
   GORM-specific techniques that obscure what's really happening in the DB.
   - Gnorm is database-first.  You generate your database using the first class tools for your DB (sql, pgAdmin, whatever). 
1. GORM is only for go.  
    - Gnorm can generate code for any language, or even docs or protobufs.
1. GORM is a framework (it calls you) that uses a ton of reflection, special tags, and naming conventions to create behaviors.  This makes it run slow, and makes reading and debugging the code difficult.
    - Gnorm generates simple, easy-to-understand boilerplate that would be tedious and error-prone to write, but offers the best speed and least complexity.

## [SqlBoiler](https://github.com/volatiletech/sqlboiler) 

SqlBoiler is a database-first fluent-api generator for Go. 

### Differences:

1. SqlBoiler only outputs Go code.
    - Gnorm outputs any langauge, even non-code things, like protobufs or docs.
1. SqlBoiler produces a fluent API style of code that is non-idiomatic for Go.
    - Gnorm can generate idiomatic go code.
1. SqlBoiler requires magic strings for queries (e.g. Where("age > ?", 30), which offer no type safety or compile-time checks.  If you change the name or type of a column, that code won't fail until it's triggered by a test or in production.
    - Gnorm can generate type-safe queries that ensure that if a column's type or name change, your code will fail to compile.

## [XO](https://github.com/knq/xo) 

XO is a code generator focused on producing Go code for your database.

### Differences

1. XO is mainly aimed at producing Go code and thus has a lot of assumptions built in (like what types you want to use to represent specific DB types).
    - Gnorm doesn't make any assumptions about what you're outputting, so you're never "swimming upstream".
1. XO only produces code in a single directory with specifically named files. If you have a big DB schema (or even if you don't), this limits you to long ugly names in order to namespace things, and prevents many more advanced layouts.
    - Gnorm lets you control the directory and filename generated for each schema, table, and custom type.
1. XO has a mode that will output code that mimics a specific query.
    - Gnorm doesn't have this, but you can generally generate code to support almost any query, and then just write the query yourself.


## Caveats

Many of these tools support more databases that Gnorm at this time (gnorm only
supports Postgres and MySQL).  All of these tools are fine works of
craftsmanship and if they fit your need, you should absolutely use them.

