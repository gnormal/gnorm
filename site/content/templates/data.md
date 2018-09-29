+++
title= "Template Data"
weight=1
alwaysopen=true
+++

The data passed to the templates is defined below.

## __Schema Data__

Data passed to each schema template:

| Property | Type | Description |
| --- | ---- | --- |
| Schema | [Schema](#schema) | The schema being rendered
| DB | [DB](#db) | The data for the whole DB
| Config | [Config](#config) | Gnorm config values from the gnorm.toml file
| Params | map[string]anything | the values from the Params entry in the config file


## __Table Data__

Data passed to each table template:

| Property | Type | Description |
| --- | ---- | --- |
| Table | [Table](#table) | the table being rendered
| DB | [DB](#db) | The data for the whole DB
| Config | [Config](#config) | Gnorm config values from the gnorm.toml file
| Params | map[string]anything | the values from the Params entry in the config file



## __Enum Data__

Data passed to each enum template:

| Property | Type | Description |
| --- | ---- | --- |
| Enum | [Enum](#enum) | the enum being rendered
| DB | [DB](#db) | The data for the whole DB
| Config | [Config](#config) | Gnorm config values from the gnorm.toml file
| Params | map[string]anything | the values from the Params entry in the config file


## __Type Definitions__
-----
These are the definitions of all the complex types referenced by the above.

### DB

Data representing the entire DB schema(s).

| Property | Type | Description |
| --- | ---- | --- |
| Schemas | list of [Schemas](#schema) | all the schemas parsed by gnorm
| SchemasByName | map[string][Schema](#schema) | map of schema DBName to Schema

### Column

Column is the data about a DB column of a table.

| Property | Type | Description |
| --- | ---- | --- |
| Table | [Table](#table) | the table this column is in
| Name  | string | the converted name of the column
| DBName | string | the original name of the column in the DB
| Type |string | the converted name of the type
| DBType | string | the original type name of the column in the DB
| IsArray | boolean | true if the column type is an array
| Length | integer | non-zero if the type has a length (e.g. varchar[16])
| UserDefined | boolean | true if the type is user-defined
| Nullable | boolean | true if the column is not NON NULL
| HasDefault | boolean | true if the column has a default
| Comment | string | the comment attached to the column
| IsPrimaryKey | boolean | true if the column is a primary key
| IsFK | boolean | true if the column is a foreign key
| HasFKRef | boolean | true if the column is referenced by a foreign key
| FKColumn | [ForeignKeyColumn](#foreignkeycolumn) | foreign key column definition
| FKColumnRefs | [ForeignKeyColumns](#foreignkeycolumns) | all foreign key columns referencing this column
| FKColumnRefsByName | map[string][ForeignKeyColumn](#foreignkeycolumn) | all foreign key columns referencing this column by foreign key name
| Orig | db-specific | the raw database column data (different per db type)

### Columns

Columns is an ordered list of [Column](#column) values from a table.  Columns
have the following properties:

| Property | Type | Description |
| --- | ---- | --- |
| DBNames | [Strings](#strings) | the ordered list of DBNames of all the columns
| Names | [Strings](#strings) | the ordered list of Names of all the columns

### ConfigData

| Property | Type | Description |
| --- | ---- | --- |
| ConnStr | string | the connection string for the database
| DBType | string | the type of db
| Schemas | list of string | the schema names to generate files for
| IncludeTables | map[string] list of string | whitelist map of schema names to table names in that schema to generate files for.
| ExcludeTables | map[string] list of string | blacklist map of schema names to table names in that schema to not generate files for.
| PostRun | list of string | the command to run on files after generation
| TypeMap | map[string]string | map of DBNames to converted names for column types
| NullableTypeMap | map[string]string | map of DBNames to converted names for column types (used when Nullable=true)
| PluginDirs | list of string | ordered list of directories to look in for plugins
| OutputDir | string | the directory where gnorm should output all its data
| StaticDir | string | the directory from which to statically copy files to outputdir

### Enum

An enum is a user-defined column type that has a set of allowable values.

| Property | Type | Description |
| --- | ---- | --- |
| Name  | string | the converted name of the enum
| DBName | string | the original name of the enum in the DB
| Schema | [Schema](#schema) | the schema the enum is in
| Table |  [Table](#table)  | (mysql only) the table this enum is part of
| Values | list of [EnumValue](#enumvalue)| the list of possible values for this enum

### Enums

Enums is a list of [Enum](#enum) values from a schem. Enums
have the following properties:

| Property | Type | Description |
| --- | ---- | --- |
| DBNames | [Strings](#strings) | the ordered list of DBNames of all the enums
| Names | [Strings](#strings) | the ordered list of Names of all the enums

### EnumValue

| Property | Type | Description |
| --- | ---- | --- |
|Name |   string | the converted label of the enum
|DBName | string | the original label of the enum in the DB
|Value |  int | the value for this enum value (order)

### ForeignKey

| Property | Type | Description |
| --- | ---- | ---|
| DBName | string | the original name of the foreign key constraint in the db
| Name | string | the converted name of the foreign key constraint
| TableDBName | string | the original name of the table in the db
| RefTableDBName | string | the original name of the foreign table in the db
| Table | [Table](#table) | the foreign key table
| RefTable | [Table](#table) | the foreign key foreign table
| FKColumns | [ForeignKeyColumns](#foreignkeycolumns) | all foreign key columns belonging to the foreign key

### ForeignKeys
ForeignKeys is a list of ForeignKey objects. The list has the following methods on it:

| Property | Type | Description |
| --- | ---- | ---|
| DBNames | [Strings](#strings) | the list of DBNames of all the foreign keys
| Names | [Strings](#strings) | the list of converted names of all the foreign keys

### ForeignKeyColumn
ForeignKeyColumn represents a column in the current table that references the value of a column in another table.

| Property | Type | Description |
| --- | ---- | ---|
| DBName | string | the original name of the foreign key constraint in the db
| ColumnDBName | string | the original name of the column in the db
| RefColumnDBName | string | the original name of the foreign column in the db
| Column | [Column](#column) | the foreign key column
| RefColumn | [Column](#column) | the referenced column

### ForeignKeyColumns
ForeignKeyColumns is a list of ForeignKeyColumn objects.  The list has the following methods:

| Property | Type | Description |
| --- | ---- | ---|
| DBNames | [Strings](#strings) | the list of DBNames of all the foreign keys
| ColumnDBNames | [Strings](#strings) | the list of column database names
| RefColumnDBNames | [Strings](#strings) | the list of foreign column database names

### Schema

A schema represents a namespace of tables and enums in a database.

| Property | Type | Description |
| --- | ---- | --- |
| Name | string | the converted name of the schema
| DBName | string | the original name of the schema in the DB
| Tables | [Tables](#tables) | the list of [Table](#table) values in this schema
| Enums | [Enums](#enums) | the list of [Enum](#enum) values in this schema
| TablesByName | map\[string\][Table](#table) | map of DBName to Table.

### Strings

Strings is a list of string values with the following methods (not avaialable with external template engines)

| Method | Arguments | Description |
| --- | ---- | --- |
| Except | vals ([]string) | Except returns a Strings value with the given values removed from the list (if they existed).  The check is case sensitive.
| Sorted | vals () | Sorted returns a sorted Strings value. 
| Sprintf | format (string) | Sprintf calls [fmt.Sprintf](https://golang.org/pkg/fmt/#Sprintf)(format, str) for every string in this value and returns the results as a new Strings value.

### Table

| Property | Type | Description |
| --- | ---- | --- |
| Name | string   | the converted name of the table
| DBName | string | the original name of the table in the DB
| Comment | string | the comment attached to the table
| Schema | [Schema](#schema)  | the schema this table is in
| Columns | [Columns](#columns) | ordered list of Database columns
| ColumnsByName | map[string][Column](#column) | map of column dbname to column
| PrimaryKeys | [Columns](#columns) | primary key columns
| HasPrimaryKey | bool | does the column have at least one primary key
| Indexes | [Indexes](#indexes) | the list of indexes on the table
| IndexesByName | map[string][Index](#index) | map index dbname to index
| ForeignKeys | [ForeignKey](#foreignKeys | list of foreign keys
| ForeignKeyRefs | [ForeignKeys](#foreignKeys) | foreign keys referencing this table
| FKByName | map[string][ForeignKey](#foreignKey) | foreign keys by foreign key name
| FKRefsByName | map[string][ForeignKey](#foreignKey) | foreign keys referencing this table by name

### Tables

Tables is list of [Table](#table) values from a schema.  Tables
have the following properties:

| Property | Type | Description |
| --- | ---- | --- |
| DBNames | [Strings](#strings) | the list of DBNames of all the tables
| Names | [Strings](#strings) | the list of Names of all the tables

### Index

| Property | Type | Description |
| --- | --- | --- |
| Name | string | the converted name of the index
| DBName | string | the name of the index from the database
| IsUnique | bool | true if the index is unique
| Columns | [Columns](#columns) | the list of the columns used in the index

### Indexes

Indexes is a list of [Index](#index) values for a table.

| Property | Type | Description
| --- | --- | --- |
| Names | [Strings](#strings) | the list of Names of the indexes
| DBNames | [Strings](#strings) | the list of DBNames of the Indexes
