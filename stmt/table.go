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
	b.WriteString(strconv.FormatInt(int64(l), 10))
	return nil
}
