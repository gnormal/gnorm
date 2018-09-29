+++
title= "Template Plugins"
weight=3
alwaysopen=true
+++

Note that template plugins are not available for external template engines.

You can call external applications as template functions. gnorm provides a
helper function for calling external functions called `plugin`.

Take for instance the following template extracted from the tests.

```plain
{{plugin "nix" "echoPlugin" . }}
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

The context sent to plugins are of the shape `{ "data": SOME_VALUE}` where
`SOME_VALUE` is any valid json object. It is up to the plugin author to decide
how to  handle the shape and form of data.

The output should be in the same format as input. The plugin can exit with non
`0` status and anything written to `stderr` will be included in the error
message when executing templates.


## Configuration

Set `PluginDirs` configuration value to the desired directories for plugin
lookup


```toml
## Add this to gnorm.toml
PluginDirs = ["/path/to/plugin/one" , "/path/to/plugin/two"]
```

So, say you call plugin `foo`, gnorm will look for foo binary in the specified
directories.The lookup is done in the same order as the directories are
specified.



