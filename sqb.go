package sqb

import (
	"errors"
	"strings"

	"github.com/Code-Hex/sqb/internal/pool"
	"github.com/Code-Hex/sqb/stmt"
)

// Builder builds sql query string.
type Builder struct {
	baseQuery string
	stmt      []stmt.Expr
}

// New returns sql query builder.
func New(sqlstr string) *Builder {
	return &Builder{
		baseQuery: sqlstr,
	}
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
func (b *Builder) Build() (string, []interface{}, error) {
	q := b.baseQuery

	buf := pool.Get()
	defer pool.Put(buf)

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
		q = b.baseQuery[offset:]
	}
	buf.WriteString(q)

	return buf.String(), buf.Args(), nil
}
