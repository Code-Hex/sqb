package stmt

// Builder an interface used to build SQL queries.
//
// WriteString method uses the passed string is writing to the query builder.
// AppendArgs method uses for pass arguments corresponding to variables.
type Builder interface {
	WritePlaceholder()
	WriteString(string)
	AppendArgs(args ...interface{})
}

// Expr implemented Write method.
//
// This interface represents an expression.
type Expr interface {
	Write(Builder) error
}

// Comparisoner implemented WriteComparison method.
//
// This interface represents a conditional expression.
type Comparisoner interface {
	WriteComparison(b Builder) error
}
