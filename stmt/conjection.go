package stmt

import "errors"

var (
	_ Expr = (*Paren)(nil)
	_ Expr = (*Or)(nil)
	_ Expr = (*And)(nil)
)

// Paren represents a parenthesized expression.
type Paren struct {
	Expr Expr
}

// Write writes the expression with parentheses.
func (p *Paren) Write(b Builder) error {
	if p.Expr == nil {
		return errors.New("unset Expr in Paren")
	}
	b.WriteString("(")
	if err := p.Expr.Write(b); err != nil {
		return err
	}
	b.WriteString(")")
	return nil
}

// Or represents an OR boolean expression.
type Or struct {
	Left  Expr
	Right Expr
}

// Write writes the OR boolean expression with parentheses.
// Currently, the OR operator is the only one that's lower precedence
// than AND on most of databases.
func (o *Or) Write(b Builder) error {
	if o.Left == nil {
		return errors.New("unset Left Expr in OR")
	}
	if o.Right == nil {
		return errors.New("unset Right Expr in OR")
	}
	b.WriteString("(")
	if err := o.Left.Write(b); err != nil {
		return err
	}
	b.WriteString(" OR ")
	if err := o.Right.Write(b); err != nil {
		return err
	}
	b.WriteString(")")
	return nil
}

// And represents an And boolean expression.
type And struct {
	Left  Expr
	Right Expr
}

// Write writes the AND boolean expression.
func (a *And) Write(b Builder) error {
	if a.Left == nil {
		return errors.New("unset Left Expr in And")
	}
	if a.Right == nil {
		return errors.New("unset Right Expr in And")
	}
	if err := a.Left.Write(b); err != nil {
		return err
	}
	b.WriteString(" AND ")
	return a.Right.Write(b)
}
