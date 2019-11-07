package sqb

import (
	"errors"
	"strings"

	"github.com/Code-Hex/sqb/internal/pool"
	"github.com/Code-Hex/sqb/stmt"
)

// There is build logic using placeholder in internal/pool.go.
const (
	// Question represents a '?' placeholder parameter.
	Question = iota
	// Dollar represents a '$1', '$2'... placeholder parameters.
	Dollar
)

// Option represents options to build sql query.
type Option func(b *Builder)

// SetPlaceholder sets placeholder.
//
// Default value is zero uses Question '?' as a placeholder.
func SetPlaceholder(placeholder int) Option {
	return func(b *Builder) {
		b.placeholder = placeholder
	}
}

// Builder builds sql query string.
type Builder struct {
	placeholder int
	stmt        []stmt.Expr
}

// New returns sql query builder.
func New(opts ...Option) *Builder {
	b := &Builder{}
	for _, opt := range opts {
		opt(b)
	}
	return b
}

// Bind binds expression to bindVars. returns copied *Builder which
// bound expression.
func (b *Builder) Bind(expr stmt.Expr) *Builder {
	// copy
	ret := *b
	copy(ret.stmt, b.stmt)
	// append to copied builder
	ret.stmt = append(b.stmt, expr)
	return &ret
}

// Build builds sql query string, returning the built query string
// and a new arg list that can be executed by a database. The `query` should
// use the `?` bindVar. The return value uses the `?` bindVar.
func (b *Builder) Build(baseQuery string) (string, []interface{}, error) {
	q := baseQuery

	buf := pool.Get()
	defer pool.Put(buf)

	buf.Placeholder = b.placeholder

	// '?' <- bindVar
	var bindVars, offset int
	for i := strings.IndexByte(q, '?'); i != -1; i = strings.IndexByte(q, '?') {
		if bindVars >= len(b.stmt) {
			// If number of statements is less than bindVars, returns an error;
			return "", nil, errors.New("number of bindVars exceeds replaceable statements")
		}

		buf.WriteString(q[:i])
		if err := b.stmt[bindVars].Write(buf); err != nil {
			return "", nil, err
		}
		bindVars++
		offset += i + 1
		q = baseQuery[offset:]
	}
	buf.WriteString(q)

	return buf.String(), buf.Args(), nil
}
