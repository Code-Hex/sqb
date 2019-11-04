package sqb_test

import (
	"errors"
	"testing"

	"github.com/Code-Hex/sqb"
	"github.com/Code-Hex/sqb/stmt"
	"github.com/google/go-cmp/cmp"
)

func TestBuilder_Build(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		stmts    []stmt.Expr
		want     string
		wantArgs []interface{}
		wantErr  bool
	}{
		{
			name: "valid where between",
			sql:  "SELECT * FROM tables WHERE ?",
			stmts: []stmt.Expr{
				sqb.Between("name", 100, 200),
			},
			want:     "SELECT * FROM tables WHERE name BETWEEN ? AND ?",
			wantArgs: []interface{}{100, 200},
			wantErr:  false,
		},
		{
			name: "valid where not between",
			sql:  "SELECT * FROM tables WHERE ?",
			stmts: []stmt.Expr{
				sqb.NotBetween("name", 100, 200),
			},
			want:     "SELECT * FROM tables WHERE name NOT BETWEEN ? AND ?",
			wantArgs: []interface{}{100, 200},
			wantErr:  false,
		},
		{
			name: "valid condition",
			sql:  "SELECT * FROM tables WHERE ? AND ?",
			stmts: []stmt.Expr{
				sqb.Eq("name", "taro"),
				sqb.Ne("category", 10),
			},
			want:     "SELECT * FROM tables WHERE name = ? AND category != ?",
			wantArgs: []interface{}{"taro", 10},
			wantErr:  false,
		},
		{
			name: "valid conject twice",
			sql:  "SELECT * FROM tables WHERE ?",
			stmts: []stmt.Expr{
				sqb.And(
					sqb.Or(
						sqb.Eq("category", 1),
						sqb.Eq("category", 2),
					),
					sqb.Or(
						sqb.NotIn("brand", []string{
							"apple", "sony", "google",
						}),
						sqb.NotLike("name", "abc%"),
					),
				),
			},
			want:     "SELECT * FROM tables WHERE (category = ? OR category = ?) AND (brand NOT IN (?, ?, ?) OR name NOT LIKE ?)",
			wantArgs: []interface{}{1, 2, "apple", "sony", "google", "abc%"},
			wantErr:  false,
		},
		{
			name:     "invalid bindVars exceeds replaceable statements",
			sql:      "SELECT * FROM tables WHERE ?",
			stmts:    []stmt.Expr{},
			want:     "",
			wantArgs: nil,
			wantErr:  true,
		},
		{
			name: "invalid build error",
			sql:  "SELECT * FROM tables WHERE ?",
			stmts: []stmt.Expr{
				&ExprMock{
					WriteMock: func(stmt.Builder) error {
						return errors.New("error")
					},
				},
			},
			want:     "",
			wantArgs: nil,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := sqb.New(tt.sql)
			for _, expr := range tt.stmts {
				b = b.Bind(expr)
			}
			got, args, err := b.Build()
			if (err != nil) != tt.wantErr {
				t.Errorf("Builder.Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("sql\ngot = %q\nwant %q", got, tt.want)
			}
			if diff := cmp.Diff(tt.wantArgs, args); diff != "" {
				t.Errorf("args (-want, +got)\n%s", diff)
			}
		})
	}
}
