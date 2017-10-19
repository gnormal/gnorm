+++
title= "Static Files"
weight=30
alwaysopen=true
+++

Gnorm supports copying static files from an input directory to the output
directory.  Set the `StaticDir` field in your gnorm.toml to the directory you
would like to copy files from, and those files will get copied to the
`OutputDir` (which defaults to the directory where Gnorm is running).  The
directory layout under StaticDir will be maintained in the output dir, so for
example, if you have `<StaticDir>/foo/bar.txt`, it will get copied to
`<OutputDir>/foo/bar.txt`. 

This can be useful for a couple purposes - first, it lets you keep non-generated
files outside the output directory, so that you can delete and recreate the
entire directory. Second, it allows users to distribute [gnorm
solutions](/solutions/) which include templates and non-generated files, similar
to how themes are distributed for static website generators.
