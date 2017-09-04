package run

import (
	"github.com/pkg/errors"

	"gnorm.org/gnorm/database"
	"gnorm.org/gnorm/database/drivers/mysql"
	"gnorm.org/gnorm/database/drivers/postgres"
	"gnorm.org/gnorm/environ"
)

func getDBInfo(env environ.Values, cfg *Config) (*database.Info, error) {
	var info *database.Info
	var err error
	switch cfg.DBType {
	case Postgres:
		info, err = postgres.Parse(env.Log, cfg.ConnStr, cfg.Schemas)
	case Mysql:
		info, err = mysql.Parse(env.Log, cfg.ConnStr, cfg.Schemas)
	default:
		return nil, errors.Errorf("unknown database type: %v", cfg.DBType)
	}
	if err != nil {
		return nil, err
	}
	if err := convertNames(env.Log, info, cfg); err != nil {
		return nil, err
	}
	return info, nil
}
