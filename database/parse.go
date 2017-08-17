package database

import (
	"database/sql"

	// import here for now
	_ "github.com/lib/pq"

	"github.com/gnormal/gnorm/database/pg"
)

//go:generate go get github.com/xoxo-go/xoxo
//go:generate xoxo pgsql://$DB_USER:$DB_PASSWORD@$DB_HOST/$DB_NAME?sslmode=$DB_SSL_MODE --schema information_schema -o pg --template-path ./templates

type Schema struct {
	Name   string
	Tables []Table
}

type Table struct {
	Name string
}

func Parse(conn string, schemas []string) ([]*pg.Table, error) {
	db, err := sql.Open("pq", conn)
	if err != nil {
		return nil, err
	}
	tables, err := pg.QueryTable(db, pg.TableTableSchemaWhere.In(schemas), pg.UnOrdered)
	if err != nil {
		return nil, err
	}
	return tables, nil
}
