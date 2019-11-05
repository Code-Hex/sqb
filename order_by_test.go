package sqb_test

import (
	"strings"
	"testing"

	"github.com/Code-Hex/sqb"
	"github.com/Code-Hex/sqb/stmt"
)

func TestOrderBy(t *testing.T) {
	type args struct {
		column string
		desc   bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ASC",
			args: args{
				column: "column",
				desc:   false,
			},
			want: "column",
		},
		{
			name: "DESC",
			args: args{
				column: "column",
				desc:   true,
			},
			want: "column DESC",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BuildCapture{
				buf:  strings.Builder{},
				Args: []interface{}{},
			}
			expr := sqb.OrderBy(tt.args.column, tt.args.desc)
			if err := expr.Write(b); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := b.buf.String(); tt.want != got {
				t.Errorf("\nwant: %q\ngot: %q", tt.want, got)
			}
		})
	}
}

func TestOrderByList(t *testing.T) {
	type args struct {
		expr  *stmt.OrderBy
		exprs []*stmt.OrderBy
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "unary",
			args: args{
				expr: sqb.OrderBy("column", false),
			},
			want: "column",
		},
		{
			name: "list",
			args: args{
				expr: sqb.OrderBy("column1", true),
				exprs: []*stmt.OrderBy{
					sqb.OrderBy("column2", true),
					sqb.OrderBy("column3", false),
				},
			},
			want: "column1 DESC, column2 DESC, column3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BuildCapture{
				buf:  strings.Builder{},
				Args: []interface{}{},
			}
			expr := sqb.OrderByList(tt.args.expr, tt.args.exprs...)
			if err := expr.Write(b); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := b.buf.String(); tt.want != got {
				t.Errorf("\nwant: %q\ngot: %q", tt.want, got)
			}
		})
	}
}
