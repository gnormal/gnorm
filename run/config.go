package run

import (
	"text/template"

	"gnorm.org/gnorm/database"
	"gnorm.org/gnorm/run/data"
)

// Config holds the schema that is expected to exist in the gnorm.toml file.
type Config struct {
	data.ConfigData

	// TablePaths is a list of output targets used to render table data.
	//
	// The filename template may reference the values .Schema and .Table,
	// containing the name of the current schema and table being rendered.  For
	// example, "{{.Schema}}/{{.Table}}/{{.Table}}.go" would render the
	// "public.users" table to ./public/users/users.go.
	TablePaths []OutputTarget

	// SchemaPaths is a list of output targets used to render schema data.
	//
	// The filename template may reference the value .Schema, containing the
	// name of the current schema being rendered. For example,
	// "schemas/{{.Schema}}/{{.Schema}}.go" would render the "public" schema to
	// ./schemas/public/public.go
	SchemaPaths []OutputTarget

	// EnumPaths is a list of output targets used to render enum data.
	//
	// The filename template may reference the values .Schema and .Enum,
	// containing the name of the current schema and Enum being rendered.  For
	// example, "gnorm/{{.Schema}}/enums/{{.Enum}}.go" would render the
	// "public.book_type" enum to ./gnorm/public/enums/users.go.
	EnumPaths []OutputTarget

	// NameConversion defines how the DBName of tables, schemas, and enums are
	// converted into their Name value.  This is a template that may use all the
	// regular functions.  The "." value is the DB name of the item. Thus, to
	// make an item's Name the same as its DBName, you'd use a template of
	// "{{.}}". To make the Name the PascalCase version of DBName, you'd use
	// "{{pascal .}}".
	NameConversion *template.Template

	// Driver holds a reference to the current database driver that was
	// registered for the DBType and can connect using ConnStr.
	Driver database.Driver

	// Params contains any data you may want to pass to your templates.  This is
	// a good way to make templates reusable with different configuration values
	// for different situations.  The values in this field will be available in
	// the .Params value for all templates.
	Params map[string]interface{}

	// TemplateEngine, if specified, describes a command line tool to run to
	// render your templates, allowing you to use your preferred templating
	// engine.  If not specified, go's text/template will be used to render.
	TemplateEngine struct {
		// CommandLine is the command to run to render the template.  You may
		// pass the following variables to the command line - {{.Data}} the name
		// of a .json file containing the gnorm data serialized into json,
		// {{.Template}} - the name of the template file being rendered,
		// {{.Output}} the target file where output should be written.
		// text/template.
		CommandLine []*template.Template

		// If true, the json data will be sent via stdin to the rendering tool.
		UseStdin bool

		// If true, the output of the tool will be written to the target file.
		UseStdout bool
	}
}

// OutputTarget contains a template that generates a filename to write to, and a
// template that generates the contents for that file.  If an external template
// engine is used, Contents will be nil, and the template at ContentsPath should
// be used.
type OutputTarget struct {
	Filename     *template.Template
	Contents     *template.Template
	ContentsPath string
}
