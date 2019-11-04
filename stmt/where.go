package stmt

import (
	"errors"
)

// Builder doesn't have the potential to return an error. But have the potential to panic.
//
// strings.Builder
// https://golang.org/src/strings/builder.go?s=3425:3477#L110
//
// bytes.Buffer
// https://golang.org/pkg/bytes/#Buffer.WriteString
type Builder interface {
	WriteString(string)
	AppendArgs(args ...interface{})
}

type Expr interface {
	Write(Builder) error
}

type Conjection struct {
	Left     Expr
	Combined string // AND OR
	Right    Expr
}

func (c *Conjection) Write(b Builder) error {
	if c.Left == nil {
		return errors.New("unset Left Expr in Conjection")
	}
	if _, ok := c.Left.(*Conjection); ok {
		b.WriteString("(")
		if err := c.Left.Write(b); err != nil {
			return err
		}
		b.WriteString(")")
	} else {
		if err := c.Left.Write(b); err != nil {
			return err
		}
	}

	if c.Combined == "" {
		return nil
	}
	b.WriteString(" ")
	b.WriteString(c.Combined)
	b.WriteString(" ")

	if c.Right == nil {
		return errors.New("unset Right Expr in Conjection")
	}

	if _, ok := c.Right.(*Conjection); ok {
		b.WriteString("(")
		if err := c.Right.Write(b); err != nil {
			return err
		}
		b.WriteString(")")
	} else {
		if err := c.Right.Write(b); err != nil {
			return err
		}
	}
	return nil
}

type Condition struct {
	Column  string
	Compare Comparisoner
}

func (c *Condition) Write(b Builder) error {
	// Column Negative Compare Value
	//
	// category = "music"
	// category != "music"
	// category LIKE "music"
	// category NOT LIKE "music"
	// category IN ("music", "video")
	// category NOT IN ("music", "video")
	b.WriteString(c.Column)
	if c.Compare == nil {
		return errors.New("unset Compare in condition")
	}
	b.WriteString(" ")
	return c.Compare.WriteComparison(b)
}

// Where "WHERE"
type Where struct {
	Expr Expr
}

func (w *Where) Write(b Builder) error {
	b.WriteString("WHERE ")
	if w.Expr == nil {
		return errors.New("unset Expr in where clause")
	}
	return w.Expr.Write(b)
}
