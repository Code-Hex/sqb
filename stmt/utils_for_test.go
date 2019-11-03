package stmt

import "strings"

var _ Builder = (*BuildCapture)(nil)

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

var _ Expr = (*ExprMock)(nil)

type ExprMock struct {
	WriteMock func(Builder) error
}

func (e *ExprMock) Write(b Builder) error {
	return e.WriteMock(b)
}
