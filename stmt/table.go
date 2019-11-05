package stmt

import (
	"errors"
	"strconv"
)

// String is able to replace bindvars with string.
//
// i.e. "SELECT * FROM ?" => "SELECT * FROM string"
type String string

// Write writes the string.
func (s String) Write(b Builder) error {
	if s == "" {
		return errors.New("unspecified string")
	}
	b.WriteString(string(s))
	return nil
}

// Numeric is able to replace bindvars with numeric.
//
// i.e. "LIMIT ?" => "LIMIT numeric"
type Numeric int64

// Write writes the numeric.
func (n Numeric) Write(b Builder) error {
	b.WriteString(strconv.FormatInt(int64(n), 10))
	return nil
}

// Limit represents "LIMIT <limit_num>".
type Limit int64

// Write writes the number of limitations that the Limit has.
func (l Limit) Write(b Builder) error {
	b.WriteString("LIMIT ")
	return Numeric(l).Write(b)
}

// Offset represents "OFFSET <offset_num>".
type Offset int64

// Write writes the number of offsets that the Offset has.
func (o Offset) Write(b Builder) error {
	b.WriteString("OFFSET ")
	return Numeric(o).Write(b)
}

// OrderBy represents "<column_name>", "<column_name> DESC".
// If there is Next, it represents like "<column_name>, <column_name> DESC".
type OrderBy struct {
	Column string
	Desc   bool
	Next   *OrderBy
}

// Write writes expression for "ORDER BY".
func (o *OrderBy) Write(b Builder) error {
	b.WriteString(o.Column)
	if o.Desc {
		b.WriteString(" DESC")
	}
	if o.Next != nil {
		b.WriteString(", ")
		return o.Next.Write(b)
	}
	return nil
}
