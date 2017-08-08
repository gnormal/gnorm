<h1 align="center">GNORM (WIP do not use)</h1>

<p align="center"><img src="https://user-images.githubusercontent.com/3185864/29083720-a7644ba2-7c37-11e7-8e3f-a9a73b7f83c5.png" /></p>
<p align="center">GNORM is Not an ORM.</p>

## About

Gnorm is a datebase-first code generator that generates boilerplate that matches your database's schema.  This allows for faster, type-safe queries than you get from a traditional ORM, without losing any of the development speed.

## Templates

Gnorm reads your database schema, then runs the resulting data through templates you can customize in any way you like.  A default set of templates produces Go structs and functions using the stdlib's database/sql package as a thorough example and usable database layer for most go projects.

However, your templates may generate whatever code or text files you wish, based on your templates.

## Configuration

Configuring gnorm is as easy as creating a simple [TOML](https://github.com/toml-lang/toml) file with a few configuration values.  Gnorm takes care of the rest.
