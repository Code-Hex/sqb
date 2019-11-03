package stmt

import "errors"

// Comparisoner implemented WriteComparison method.
//
// This interface represents a conditional expression.
type Comparisoner interface {
	WriteComparison(b Builder) error
}

var (
	_ Comparisoner = (*CompOp)(nil)
	_ Comparisoner = (*CompLike)(nil)
	_ Comparisoner = (*CompBetween)(nil)
	_ Comparisoner = (*CompIn)(nil)
)

// CompOp represents condition for using operators.
//
// Op field should contain "=", ">=", ">", "<=", "<", "!=", "IS", "IS NOT"
// Value field should set the value to use for comparison.
type CompOp struct {
	Op    string
	Value interface{}
}

// WriteComparison implemented Comparisoner interface.
func (c *CompOp) WriteComparison(b Builder) error {
	b.WriteString(c.Op)
	b.WriteString(" ?")
	b.AppendArgs(c.Value)
	return nil
}

// CompLike represents condition for using "LIKE".
//
// If enabled Negative field, it's meaning use "NOT LIKE".
// Value field should set the value to use for comparison.
type CompLike struct {
	Negative bool
	Value    interface{}
}

// WriteComparison implemented Comparisoner interface.
func (c *CompLike) WriteComparison(b Builder) error {
	if c.Negative {
		b.WriteString("NOT ")
	}
	b.WriteString("LIKE ?")
	b.AppendArgs(c.Value)
	return nil
}

// CompBetween represents condition for using "BETWEEN".
//
// If enabled Negative field, it's meaning use "NOT BETWEEN".
// This struct will convert to be like "BETWEEN left_expr AND right_expr".
type CompBetween struct {
	Negative bool
	Left     interface{}
	Right    interface{}
}

// WriteComparison implemented Comparisoner interface.
func (c *CompBetween) WriteComparison(b Builder) error {
	if c.Left == nil {
		return errors.New("unset Left Value in CompBetween")
	}
	if c.Right == nil {
		return errors.New("unset Right Value in CompBetween")
	}
	if c.Negative {
		b.WriteString("NOT ")
	}
	b.WriteString("BETWEEN ? AND ?")
	b.AppendArgs(c.Left, c.Right)
	return nil
}

// CompIn represents condition for using "IN".
//
// If enabled Negative field, it's meaning use "NOT IN".
// Values field should set list to use for comparison.
// This struct will convert to be like "IN (?, ?, ?)".
type CompIn struct {
	Negative bool
	Values   []interface{}
}

// WriteComparison implemented Comparisoner interface.
func (c *CompIn) WriteComparison(b Builder) error {
	if c.Negative {
		b.WriteString("NOT ")
	}
	b.WriteString("IN (")
	if err := makePlaceholders(b, c.Values); err != nil {
		return err
	}
	b.WriteString(")")
	b.AppendArgs(c.Values...)
	return nil
}
