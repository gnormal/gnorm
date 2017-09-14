package run

import (
	"gnorm.org/gnorm/database"
	"gnorm.org/gnorm/environ"
)

func getDBInfo(env environ.Values, cfg *Config) (*database.Info, error) {
	info, err := cfg.Driver.Parse(env.Log, cfg.ConnStr, cfg.Schemas, makeFilter(cfg.IncludeTables, cfg.ExcludeTables))
	if err != nil {
		return nil, err
	}
	if err := convertNames(env.Log, info, cfg); err != nil {
		return nil, err
	}
	return info, nil
}

func makeFilter(include, exclude map[string][]string) func(schema, table string) bool {
	if sumLens(include) == 0 && sumLens(exclude) == 0 {
		return func(_, _ string) bool { return true }
	}
	if sumLens(include) == 0 {
		return func(schema, table string) bool {
			return !contains(exclude[schema], table)
		}
	}
	return func(schema, table string) bool {
		return contains(include[schema], table)
	}
}

func contains(vals []string, s string) bool {
	for x := range vals {
		if vals[x] == s {
			return true
		}
	}
	return false
}

// sumLens returns the sum of all the lengths of arrays in the map.
func sumLens(m map[string][]string) int {
	length := 0
	for k := range m {
		length += len(m[k])
	}
	return length
}
