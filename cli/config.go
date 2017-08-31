package cli

// Config holds the schema that is expected to exist in the gnorm.toml file.
type Config struct {
	// ConnStr is the connection string for the database.  environment variables
	// in $FOO form will be expanded.
	ConnStr string

	// Schemas holds the names of schemas to generate code for.
	Schemas []string

	// TemplateDir contains the relative path to the directory where gnorm
	// expects to find templates to render.  The default is the current
	// directory where gnorm is running.
	TemplateDir string

	// TablePath is a relative path for tables to be rendered.  The table
	// template will be rendered with each table in turn. If the path is empty,
	// tables will not be rendered.
	//
	// The table path may be a template, in which case the values .Schema and
	// .Table may be referenced, containing the name of the current schema and
	// table being rendered.  For example, "{{.Schema}}/{{.Table}}/{{.Table}}.go" would render
	// the "public.users" table to ./public/users/users.go.
	TablePath string

	// SchemaPath is a relative path for schemas to be rendered.  The schema
	// template will be rendered with each schema in turn. If the path is empty,
	// schema will not be rendered.
	//
	// The schema path may be a template, in which case the value .Schema may be
	// referenced, containing the name of the current schema being rendered. For
	// example, "schemas/{{.Schema}}/{{.Schema}}.go" would render the "public"
	// schema to ./schemas/public/public.go
	SchemaPath string

	// EnumPath is a relative path for enums to be rendered.  The enum.tpl template
	// will be rendered with each enum in turn. If the path is empty, enums will not
	// be rendered this way (thought you could render them via the schemas template).
	//
	// The enum path may be a template, in which case the values .Schema and .Enum
	// may be referenced, containing the name of the current schema and Enum being
	// rendered.  For example, "gnorm/{{.Schema}}/enums/{{.Enum}}.go" would render
	// the "public.book_type" enum to ./gnorm/public/enums/users.go.
	EnumPath string

	// TypeMap is a mapping of database type names to replacement type names
	// (generally types from your language for deserialization).  Types not in
	// this list will remain in their database form.  In the data sent to your
	// template, this is the Column.Type, and the original type is in
	// Column.OrigType.  Note that because of the way tables in TOML work,
	// TypeMap and NullableTypeMap must be at the end of your configuration
	// file.
	TypeMap map[string]string

	// NullableTypeMap is a mapping of database type names to replacement type
	// names (generally types from your language for deserialization)
	// specifically for database columns that are nullable.  Types not in this
	// list will remain in their database form.  In the data sent to your
	// template, this is the Column.Type, and the original type is in
	// Column.OrigType.   Note that because of the way tables in TOML work,
	// TypeMap and NullableTypeMap must be at the end of your configuration
	// file.
	NullableTypeMap map[string]string
}
