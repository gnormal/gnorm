// Code generated by gnorm, DO NOT EDIT!

package columns

import (
	"database/sql"
	"log"
	"strings"

	"gnorm.org/gnorm/database/drivers/sqlite/gnorm"
)

// Row represents a row from 'COLUMNS'.
type Row struct {
	ColumnName      string         // COLUMN_NAME
	OrdinalPosition int64          // ORDINAL_POSITION
	ColumnDefault   sql.NullString // COLUMN_DEFAULT
	IsNullable      bool           // IS_NULLABLE
	DataType        string         // DATA_TYPE
	ColumnType      string         // COLUMN_TYPE
	ColumnKey       int            // COLUMN_KEY
}

// Query retrieves rows from 'COLUMNS' as a slice of Row.
func Query(db gnorm.DB, table string) ([]*Row, error) {
	const sqlstr = `

SELECT cid, 
       name, 
       type,
       "notnull",
       dflt_value, 
       pk 
  FROM pragma_table_info('XXXTABLEXXX'); 
`

	execsql := strings.Replace(sqlstr, "XXXTABLEXXX", table, -1)
	log.Println("querying columns ", execsql)
	log.Println("from table ", table)

	var vals []*Row
	q, err := db.Query(execsql)
	if err != nil {
		return nil, err
	}
	for q.Next() {
		r := Row{}

		err = q.Scan(&r.OrdinalPosition, &r.ColumnName, &r.DataType, &r.IsNullable, &r.ColumnDefault, &r.ColumnKey)
		if err != nil {
			return nil, err
		}

		vals = append(vals, &r)
	}
	return vals, nil
}
