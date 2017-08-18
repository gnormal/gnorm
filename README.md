<h1 align="center">GNORM (WIP do not use)</h1>

<p align="center"><img src="https://user-images.githubusercontent.com/3185864/29083720-a7644ba2-7c37-11e7-8e3f-a9a73b7f83c5.png" /></p>
<p align="center">GNORM is Not an ORM.</p>

## About

Gnorm is a database-first code generator that generates boilerplate that matches
your database's schema.  

Gnorm uses templates you control, so that you can make the output look exactly
how you want it to look.  It can be used to generate type-safe database queries
that are faster than a traditional ORM.  It can also be used to create a REST or
RPC API that exposes the data in your database.

## Documentation

[https://gnorm.org/gnorm](https://gnorm.org/gnorm)

## Building

Gnorm uses a build script called make.go.  To install gnorm, just run 

`go run make.go`

You can also just run `go get` or `go install`, but the output of gnorm -version
will not reflect the build environment.

## Contributing

See CONTRIBUTING.md