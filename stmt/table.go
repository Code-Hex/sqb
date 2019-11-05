package stmt

import (
	"errors"
	"strconv"
)

// Table represents "FROM <table>", "INTO <table>", "UPDATE <table>".
type Table string

// Write writes the table name that the From has.
func (f Table) Write(b Builder) error {
	if f == "" {
		return errors.New("unspecified table")
	}
	b.WriteString(string(f))
	return nil
}

// Limit represents "LIMIT <limit_num>".
type Limit int64

// Write writes the number of limitations that the Limit has.
func (l Limit) Write(b Builder) error {
	b.WriteString("LIMIT ")
	b.WriteString(strconv.FormatInt(int64(l), 10))
	return nil
}

// Offset represents "OFFSET <offset_num>".
type Offset int64

// Write writes the number of offsets that the Offset has.
func (o Offset) Write(b Builder) error {
	b.WriteString("OFFSET ")
	b.WriteString(strconv.FormatInt(int64(o), 10))
	return nil
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
