package run

import (
	"bytes"
	"log"

	"github.com/pkg/errors"
	"gnorm.org/gnorm/database"
	"gnorm.org/gnorm/run/data"
)

type nameConverter func(s string) (string, error)

func makeData(log *log.Logger, info *database.Info, cfg *Config) (*data.DBData, error) {
	convert := func(s string) (string, error) {
		buf := &bytes.Buffer{}
		err := cfg.NameConversion.Execute(buf, s)
		if err != nil {
			return "", errors.WithMessage(err, "name conversion failed for "+s)
		}
		return buf.String(), nil
	}

	db := &data.DBData{
		SchemasByName: make(map[string]*data.Schema, len(info.Schemas)),
	}
	var err error
	for _, s := range info.Schemas {
		sch := &data.Schema{
			DBName:       s.Name,
			TablesByName: make(map[string]*data.Table, len(s.Tables)),
		}
		db.Schemas = append(db.Schemas, sch)
		db.SchemasByName[sch.DBName] = sch

		sch.Name, err = convert(s.Name)
		if err != nil {
			return nil, errors.WithMessage(err, "schema")
		}
		for _, e := range s.Enums {
			enum := &data.Enum{
				DBName: e.Name,
				Schema: sch,
			}
			sch.Enums = append(sch.Enums, enum)
			enum.Name, err = convert(e.Name)
			if err != nil {
				return nil, errors.WithMessage(err, "enum")
			}
			for _, v := range e.Values {
				val := &data.EnumValue{
					DBName: v.Name,
					Value:  v.Value,
				}
				enum.Values = append(enum.Values, val)
				val.Name, err = convert(v.Name)
				if err != nil {
					return nil, errors.WithMessage(err, "enum value")
				}
			}
		}
		for _, t := range s.Tables {
			table := &data.Table{
				DBName:        t.Name,
				Comment:       t.Comment,
				Schema:        sch,
				ColumnsByName: make(map[string]*data.Column, len(t.Columns)),
				IndexesByName: make(map[string]*data.Index, len(t.Indexes)),
				FKByName:      map[string]*data.ForeignKey{},
				FKRefsByName:  map[string]*data.ForeignKey{},
			}
			sch.Tables = append(sch.Tables, table)
			sch.TablesByName[table.DBName] = table
			table.Name, err = convert(t.Name)
			if err != nil {
				return nil, errors.WithMessage(err, "table")
			}
			for _, c := range t.Columns {
				col := &data.Column{
					Table:              table,
					DBName:             c.Name,
					DBType:             c.Type,
					ColumnType:         c.ColumnType,
					IsArray:            c.IsArray,
					Length:             c.Length,
					UserDefined:        c.UserDefined,
					Nullable:           c.Nullable,
					HasDefault:         c.HasDefault,
					Comment:            c.Comment,
					IsPrimaryKey:       c.IsPrimaryKey,
					IsFK:               c.IsForeignKey,
					FKColumnRefsByName: map[string]*data.ForeignKeyColumn{},
					Orig:               c.Orig,
				}
				table.Columns = append(table.Columns, col)
				table.ColumnsByName[col.DBName] = col
				col.Name, err = convert(c.Name)
				if err != nil {
					return nil, errors.WithMessage(err, "column")
				}
				var ok bool
				if c.Nullable {
					col.Type, ok = cfg.NullableTypeMap[c.Type]
					if !ok {
						log.Println("Unmapped nullable type:", c.Type)
					}
				} else {
					col.Type, ok = cfg.TypeMap[c.Type]
					if !ok {
						log.Println("Unmapped type:", c.Type)
					}
				}
			}
			table.PrimaryKeys = filterPrimaryKeyColumns(table.Columns)

			for _, i := range t.Indexes {
				index := &data.Index{
					DBName:   i.Name,
					IsUnique: i.IsUnique,
				}
				for _, c := range i.Columns {
					index.Columns = append(index.Columns, table.ColumnsByName[c.Name])
				}

				index.Name, err = convert(i.Name)
				if err != nil {
					return nil, errors.WithMessage(err, "index")
				}

				table.Indexes = append(table.Indexes, index)
				table.IndexesByName[index.DBName] = index
			}
		}
		if err = mapSchemaForeignKeyReferences(s, sch, convert); err != nil {
			return nil, err
		}
	}
	return db, nil
}

func filterPrimaryKeyColumns(columns data.Columns) data.Columns {
	var pkColumns data.Columns
	for _, column := range columns {
		if column.IsPrimaryKey {
			pkColumns = append(pkColumns, column)
		}
	}

	return pkColumns
}

func mapSchemaForeignKeyReferences(isch *database.Schema, sch *data.Schema, convert nameConverter) error {
	for _, t := range isch.Tables {
		table, ok := sch.TablesByName[t.Name]
		if !ok {
			log.Printf("Unmapped table %v in %v", t.Name, isch.Name)
			continue
		}

		fkColumnsByFKNames := map[string]data.ForeignKeyColumns{}

		for _, c := range t.Columns {
			column, ok := table.ColumnsByName[c.Name]
			if !ok {
				log.Printf("Unmapped column %v in %v.%v", c.Name, isch.Name, t.Name)
				continue
			}

			if column.IsFK {
				refTable, ok := sch.TablesByName[c.ForeignKey.ForeignTableName]
				if !ok {
					log.Printf("Unmapped foreign table %v in %v", c.ForeignKey.ForeignTableName, isch.Name)
					continue
				}
				refColumn, ok := refTable.ColumnsByName[c.ForeignKey.ForeignColumnName]
				if !ok {
					log.Printf("Unmapped foreign column %v in %v.%v", c.ForeignKey.ForeignColumnName, isch.Name, c.ForeignKey.ForeignTableName)
					continue
				}

				fkColumn := &data.ForeignKeyColumn{
					DBName:          c.ForeignKey.Name,
					ColumnDBName:    column.DBName,
					RefColumnDBName: refColumn.DBName,
					Column:          column,
					RefColumn:       refColumn,
				}
				column.FKColumn = fkColumn

				refColumn.HasFKRef = true
				refColumn.FKColumnRefs = append(refColumn.FKColumnRefs, fkColumn)
				refColumn.FKColumnRefsByName[fkColumn.DBName] = fkColumn

				if _, ok := fkColumnsByFKNames[fkColumn.DBName]; !ok {
					fkColumnsByFKNames[fkColumn.DBName] = data.ForeignKeyColumns{fkColumn}
				} else {
					fkColumnsByFKNames[fkColumn.DBName] = append(fkColumnsByFKNames[fkColumn.DBName], fkColumn)
				}
			}
		}

		for _, fkc := range fkColumnsByFKNames {
			err := mapForeignTable(fkc, convert)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func mapForeignTable(fkc data.ForeignKeyColumns, convert nameConverter) error {
	if len(fkc) == 0 {
		return nil
	}

	// All ForeignKeyColumns will point to same table/refTable and have the same name, use first one
	table := fkc[0].Column.Table
	refTable := fkc[0].RefColumn.Table
	cName, err := convert(fkc[0].DBName)
	if err != nil {
		return errors.Wrap(err, "foreign key")
	}

	fk := &data.ForeignKey{
		DBName:         fkc[0].DBName,
		Name:           cName,
		TableDBName:    table.DBName,
		RefTableDBName: refTable.DBName,
		Table:          table,
		RefTable:       refTable,
		FKColumns:      fkc,
	}

	table.ForeignKeys = append(table.ForeignKeys, fk)
	refTable.ForeignKeyRefs = append(table.ForeignKeyRefs, fk)
	table.FKByName[fk.DBName] = fk
	refTable.FKRefsByName[fk.DBName] = fk

	return nil
}
