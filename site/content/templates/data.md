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
| Name  | string | the converted name of the column
| DBName | string | the original name of the column in the DB
| Type |string | the converted name of the type
| DBType | string | the original type name of the column in the DB 
| IsArray | boolean | true if the column type is an array 
| Length | integer | non-zero if the type has a length (e.g. varchar[16])
| UserDefined | boolean | true if the type is user-defined
| Nullable | boolean | true if the column is not NON NULL
| HasDefault | boolean | true if the column has a default
| Orig | db-specific | the raw database column data (different per db type)



### Columns

Columns is an ordered list of [Column](#column) values from a table.  Columns
have the following properties:

| Property | Type | Description |
| --- | ---- | --- |
| DBNames | [Strings](#strings) | the ordered list of DBNames of all the columns 
| Names | [Strings](#stings) | the ordered list of Names of all the columns

### ConfigData

| Property | Type | Description |
| --- | ---- | --- |
| ConnStr | string | the connection string for the database
| Schemas | list of string | the schema names to generate files for
| IncludeTables | map[string] list of string | whitelist map of schema names to table names in that schema to generate files for.
| ExcludeTables | map[string] list of string | blacklist map of schema names to table names in that schema to not generate files for.
| PostRun | list of string | the command to run on files after generation
| TypeMap | map[string]string | map of DBNames to converted names for column types
| NullableTypeMap | map[string]string | map of DBNames to converted names for column types (used when Nullable=true)

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
| Names | [Strings](#stings) | the ordered list of Names of all the enums

### EnumValue

| Property | Type | Description |
| --- | ---- | --- |
|Name |   string | the converted label of the enum
|DBName | string | the original label of the enum in the DB
|Value |  int | the value for this enum value (order)

### Schema

A schema represents a namespace of tables and enums in a database.

| Property | Type | Description |
| --- | ---- | --- |
| Name | string the converted name of the schema
| DBName | string | the original name of the schema in the DB 
| Tables | [Tables](#tables) | the list of [Table](#table) values in this schema 
| Enums | [Enums](#enums) | the list of [Enum](#enum) values in this schema
| TablesByName | map\[string\][Table](#table) | map of DBName to Table.

### Strings

Strings is a list of string values with the following methods

| Method | Arguments | Description |
| --- | ---- | --- |
| Sprintf | format (string) | Sprintf calls [fmt.Sprintf](https://golang.org/pkg/fmt/#Sprintf)(format, str) for every string in this value and returns the results as a new Strings value.

### Table

| Property | Type | Description |
| --- | ---- | --- |
| Name | string   | the converted name of the table
| DBName | string | the original name of the table in the DB
| Schema | [Schema](#schema)  | the schema this table is in
| Columns | [Columns](#columns) | ordered list of Database columns
| ColumnsByName | map[string][Column](#column) | map of column dbname to column

### Tables

Tables is list of [Table](#table) values from a schema.  Tables
have the following properties:

| Property | Type | Description |
| --- | ---- | --- |
| DBNames | [Strings](#strings) | the list of DBNames of all the tables
| Names | [Strings](#stings) | the list of Names of all the tables
