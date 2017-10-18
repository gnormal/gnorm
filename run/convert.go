package run

import (
	"bytes"
	"log"

	"github.com/pkg/errors"
	"gnorm.org/gnorm/database"
	"gnorm.org/gnorm/run/data"
)

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
				Schema:        sch,
				ColumnsByName: make(map[string]*data.Column, len(t.Columns)),
			}
			sch.Tables = append(sch.Tables, table)
			sch.TablesByName[table.DBName] = table
			table.Name, err = convert(t.Name)
			if err != nil {
				return nil, errors.WithMessage(err, "table")
			}
			for _, c := range t.Columns {
				col := &data.Column{
					DBName:       c.Name,
					DBType:       c.Type,
					IsArray:      c.IsArray,
					Length:       c.Length,
					UserDefined:  c.UserDefined,
					Nullable:     c.Nullable,
					HasDefault:   c.HasDefault,
					IsPrimaryKey: c.IsPrimaryKey,
					Orig:         c.Orig,
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
		}
	}
	return db, nil
}

func filterPrimaryKeyColumns(columns []*data.Column) []*data.Column {
	var pkColumns []*data.Column
	for _, column := range columns {
		if column.IsPrimaryKey {
			pkColumns = append(pkColumns, column)
		}
	}

	return pkColumns
}
