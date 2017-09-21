+++
title= "Template Plugins"
weight=3
alwaysopen=true
+++

You can call external applications as template functions. ngorm provides a
helper function for calling external functions called `plugin`.

Take for instance the following template extracted from the tests.

```
{{range plugin "nix" "echoPlugin" . }}{{.}}{{end}}
```

The snippet says call plugin `nix` and pass  `echoPlugin` as the first argument,
we refer this argument as a function call on the plugin, and then pass the
current context to the plugin.

## What is a plugin?

A plugin is a commandline executable with the following properties

- Accepts a function name as the first argument.
- Receives context in json format from `stdin` i.e it reads the context from
  `stdin`
- Writes the processed output to `stdout` in json format

The conext send to plugins are of the shape `{ "data": SOME_VALUE}` where
`SOME_VALUE` is any valid json object. It is up to the plugin author to decide
how to  handle the shape and form of data.

The output should be in the same format as input. The plugin can exit with non
`0` status and anything written to `stderr` will be included in the error
message when executing templates.


## Configuration
By default ngorm looks for the plugin in the following order

- `$PWD/plugins`
- System `$PATH`

However there is an option to specify custom directories to look for plugins  via the configuration file.

Note that, at the moment only absolute paths are supported. For example

```
## Add this to ngorm.toml
PluginDirs = ["/path/to/plugin/one" , "/path/to/plugin/two"]
```

So, say you call plugin `foo`, gnorm will look for foo binary in the following
order

- `/path/to/plugin/one`
- `/path/to/plugin/two`
- `$PWD/plugins`
- System `$PATH`