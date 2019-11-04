package stmt

import (
	"errors"
)

var _ Expr = (*Condition)(nil)

// Condition represents condition for using Comparisoner interface.
//
// this struct creates "<column_name> <comparable_condition>"
// <comparable_condition> indicates Comparisoner interface.
type Condition struct {
	Column  string
	Compare Comparisoner
}

// Write implements Expr interface.
//
// For example:
// category = "music"
// category != "music"
// category LIKE "music"
// category NOT LIKE "music"
// category IN ("music", "video")
// category NOT IN ("music", "video")
func (c *Condition) Write(b Builder) error {
	b.WriteString(c.Column)
	if c.Compare == nil {
		return errors.New("unset Compare in condition")
	}
	b.WriteString(" ")
	return c.Compare.WriteComparison(b)
}
