package run

import (
	"bytes"
	"log"

	"github.com/pkg/errors"
	"gnorm.org/gnorm/database"
)

func convertNames(log *log.Logger, info *database.Info, cfg *Config) error {
	convert := func(s string) (string, error) {
		buf := &bytes.Buffer{}
		err := cfg.NameConversion.Execute(buf, s)
		if err != nil {
			return "", errors.WithMessage(err, "name conversion failed for "+s)
		}
		return buf.String(), nil
	}
	var err error
	for _, s := range info.Schemas {
		s.Name, err = convert(s.DBName)
		if err != nil {
			return errors.WithMessage(err, "schema")
		}
		for _, e := range s.Enums {
			e.Schema = s.Name
			e.DBSchema = s.DBName
			e.Name, err = convert(e.DBName)
			if err != nil {
				return errors.WithMessage(err, "enum")
			}
			for _, v := range e.Values {
				v.Name, err = convert(v.DBName)
				if err != nil {
					return errors.WithMessage(err, "enum value")
				}
			}
		}
		for _, t := range s.Tables {
			t.Schema = s.Name
			t.DBSchema = s.DBName
			t.Name, err = convert(t.DBName)
			if err != nil {
				return errors.WithMessage(err, "table")
			}
			for _, c := range t.Columns {
				c.Name, err = convert(c.DBName)
				if err != nil {
					return errors.WithMessage(err, "column")
				}
				var ok bool
				if c.Nullable {
					c.Type, ok = cfg.NullableTypeMap[c.DBType]
					if !ok {
						log.Println("Unmapped nullable type:", c.DBType)
					}
				} else {
					c.Type, ok = cfg.TypeMap[c.DBType]
					if !ok {
						log.Println("Unmapped type:", c.DBType)
					}
				}
			}
		}
	}
	return nil
}
