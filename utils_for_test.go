package sqb_test

import (
	"strings"

	"github.com/Code-Hex/sqb/stmt"
)

var _ stmt.Builder = (*BuildCapture)(nil)

type BuildCapture struct {
	buf  strings.Builder
	Args []interface{}
}

func (b *BuildCapture) WriteString(s string) {
	b.buf.WriteString(s)
}

func (b *BuildCapture) AppendArgs(args ...interface{}) {
	b.Args = append(b.Args, args...)
}

var _ stmt.Expr = (*ExprMock)(nil)

type ExprMock struct {
	WriteMock func(stmt.Builder) error
}

func (e *ExprMock) Write(b stmt.Builder) error {
	return e.WriteMock(b)
}
