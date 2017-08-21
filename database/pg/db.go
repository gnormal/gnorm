package pg // import "gnorm.org/gnorm/database/pg"


// Note that this file is *NOT* generated. :)

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type SQLIdentifier = sql.NullString
type SQLIdentifierField = NullStringField

type CardinalNumber = sql.NullInt64
type CardinalNumberField = NullInt64Field

type CharacterData = sql.NullString
type CharacterDataField = NullStringField

type Oid = sql.NullInt64
type OidField = NullInt64Field

type Name = string
type NameField = StringField

type TimeStamp = time.Time
type TimeStampField = TimeField

type YesOrNo bool 

// Value marshals the value into the database
func (y YesOrNo) Value() (driver.Value, error) {
	if y {
		return "YES", nil
	}
	return "NO", nil
}

// Scan Unmarshalls the bytes[] back into a YesOrNo object
func (y *YesOrNo) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	s, ok := src.(string)
	if !ok {
		return errors.Errorf("expected YesOrNo to be a string, but was %T", src)
	}
	switch s {
		case "YES":
		*y = true
		return nil
		case "NO":
		*y = false 
		return nil
		default:
		return errors.New("unexpected value for YesOrNo: "+ s)
	}
}

// XODB is the common interface for database operations.
//
// This should work with database/sql.DB and database/sql.Tx.
type XODB interface {
	Exec(string, ...interface{}) (sql.Result, error)
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
}

// Jsonb is a wrapper for map[string]interface{} for storing json into postgres
type Jsonb map[string]interface{}

// Value marshals the json into the database
func (j Jsonb) Value() (driver.Value, error) {
	return json.Marshal(j)
}

// Scan Unmarshalls the bytes[] back into a Jsonb object
func (j *Jsonb) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Type assertion .([]byte) failed")
	}

	var i interface{}
	err := json.Unmarshal(source, &i)
	if err != nil {
		return err
	}

	if i == nil {
		return nil
	}

	*j, ok = i.(map[string]interface{})
	if !ok {
		return errors.New("reading from DB into Jsonb, failed to convert to map[string]interface{}")
	}

	return nil
}

// UnOrdered is a convenience value to make it clear you're not sorting a query.
var UnOrdered = OrderBy{}

// OrderByDesc returns a sort order descending by the given field.
func OrderByDesc(field string) OrderBy {
	return OrderBy{
		Field: field,
		Order: OrderDesc,
	}
}

// OrderByAsc returns a sort order ascending by the given field.
func OrderByAsc(field string) OrderBy {
	return OrderBy{
		Field: field,
		Order: OrderAsc,
	}
}

// OrderBy indicates how rows should be sorted.
type OrderBy struct {
	Field string
	Order SortOrder
}

func (o OrderBy) String() string {
	if o.Order == OrderNone {
		return ""
	}
	return " ORDER BY " + o.Field + " " + o.Order.String() + " "
}

// SortOrder defines how to order rows returned.
type SortOrder int

// Defined sort orders for not sorted, descending and ascending.
const (
	OrderNone SortOrder = iota
	OrderDesc
	OrderAsc
)

// String returns the sql string representation of this sort order.
func (s SortOrder) String() string {
	switch s {
	case OrderDesc:
		return "DESC"
	case OrderAsc:
		return "ASC"
	}
	return ""
}

// WhereClause has a String function should return a properly formatted where
// clause (not including the WHERE) for positional arguments starting at idx.
type WhereClause interface {
	String(idx *int) string
	Values() []interface{}
}

type comparison string

const (
	compEqual   comparison = " = "
	compGreater comparison = " > "
	compLess    comparison = " < "
	compGTE     comparison = " >= "
	compLTE     comparison = " <= "
	compNE      comparison = " <> "
)

type inClause struct {
	field  string
	values []interface{}
}

func (in inClause) String(idx *int) string {
	ret := in.field + " in ("
	for x := range in.values {
		if x != 0 {
			ret += ", "
		}
		ret += "$" + strconv.Itoa(*idx)
		(*idx)++
	}
	ret += ")"
	return ret
}

func (in inClause) Values() []interface{} {
	return in.values
}

type whereClause struct {
	field string
	comp  comparison
	value interface{}
}

func (w whereClause) String(idx *int) string {
	ret := w.field + string(w.comp) + "$" + strconv.Itoa(*idx)
	(*idx)++
	return ret
}

func (w whereClause) Values() []interface{} {
	return []interface{}{w.value}
}

// AndClause returns a WhereClause that serializes to the AND
// of all the given where clauses.
func AndClause(wheres ...WhereClause) WhereClause {
	return andClause(wheres)
}

type andClause []WhereClause

func (a andClause) String(idx *int) string {
	wheres := make([]string, len(a))
	for x := 0; x < len(a); x++ {
		wheres[x] = a[x].String(idx)
	}
	return strings.Join(wheres, " AND ")
}

func (a andClause) Values() []interface{} {
	vals := make([]interface{}, 0, len(a))
	for x := 0; x < len(a); x++ {
		vals = append(vals, a[x].Values()...)
	}
	return vals
}

// OrClause returns a WhereClause that serializes to the OR
// of all the given where clauses.
func OrClause(wheres ...WhereClause) WhereClause {
	return orClause(wheres)
}

type orClause []WhereClause

func (o orClause) String(idx *int) string {
	wheres := make([]string, len(o))
	for x := 0; x < len(wheres); x++ {
		wheres[x] = o[x].String(idx)
	}
	return strings.Join(wheres, " OR ")
}

func (o orClause) Values() []interface{} {
	vals := make([]interface{}, len(o))
	for x := 0; x < len(o); x++ {
		vals = append(vals, o[x].Values()...)
	}
	return vals
}
