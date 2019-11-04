package sqb_test

import (
	"strings"
	"testing"
	"time"

	"github.com/Code-Hex/sqb"
	"github.com/Code-Hex/sqb/stmt"
	"github.com/google/go-cmp/cmp"
)

func TestAnd(t *testing.T) {
	type args struct {
		left  stmt.Expr
		right stmt.Expr
		exprs []stmt.Expr
	}
	tests := []struct {
		name     string
		args     args
		want     string
		wantArgs []interface{}
	}{
		{
			name: "simple",
			args: args{
				left:  sqb.Eq("category", "music"),
				right: sqb.Eq("category", "home appliances"),
			},
			want:     "category = ? AND category = ?",
			wantArgs: []interface{}{"music", "home appliances"},
		},
		{
			name: "complex",
			args: args{
				left:  sqb.Eq("category", "music"),
				right: sqb.Eq("category", "home appliances"),
				exprs: []stmt.Expr{
					sqb.Ne("sub_category", 1),
					sqb.Le("sub_category", 2),
				},
			},
			want:     "category = ? AND category = ? AND sub_category != ? AND sub_category <= ?",
			wantArgs: []interface{}{"music", "home appliances", 1, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BuildCapture{
				buf:  strings.Builder{},
				Args: []interface{}{},
			}
			expr := sqb.And(tt.args.left, tt.args.right, tt.args.exprs...)
			if err := expr.Write(b); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := b.buf.String(); tt.want != got {
				t.Errorf("\nwant: %q\ngot: %q", tt.want, got)
			}
			if diff := cmp.Diff(tt.wantArgs, b.Args); diff != "" {
				t.Errorf("args (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestOr(t *testing.T) {
	type args struct {
		left  stmt.Expr
		right stmt.Expr
		exprs []stmt.Expr
	}
	tests := []struct {
		name     string
		args     args
		want     string
		wantArgs []interface{}
	}{
		{
			name: "simple",
			args: args{
				left:  sqb.Eq("category", "music"),
				right: sqb.Eq("category", "home appliances"),
			},
			want:     "(category = ? OR category = ?)",
			wantArgs: []interface{}{"music", "home appliances"},
		},
		{
			name: "complex",
			args: args{
				left:  sqb.Eq("category", "music"),
				right: sqb.Eq("category", "home appliances"),
				exprs: []stmt.Expr{
					sqb.Ge("sub_category", 1),
					sqb.Lt("sub_category", 2),
				},
			},
			want:     "(((category = ? OR category = ?) OR sub_category >= ?) OR sub_category < ?)",
			wantArgs: []interface{}{"music", "home appliances", 1, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BuildCapture{
				buf:  strings.Builder{},
				Args: []interface{}{},
			}
			expr := sqb.Or(tt.args.left, tt.args.right, tt.args.exprs...)
			if err := expr.Write(b); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := b.buf.String(); tt.want != got {
				t.Errorf("\nwant: %q\ngot: %q", tt.want, got)
			}
			if diff := cmp.Diff(tt.wantArgs, b.Args); diff != "" {
				t.Errorf("args (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestAndFromMap(t *testing.T) {
	type args struct {
		f sqb.ConditionalFunc
		m map[string]interface{}
	}
	tests := []struct {
		name     string
		args     args
		want     string
		wantArgs []interface{}
	}{
		{
			name: "an argument",
			args: args{
				f: sqb.Gt,
				m: map[string]interface{}{
					"col": 2,
				},
			},
			want:     "col > ?",
			wantArgs: []interface{}{2},
		},
		{
			name: "two arguments",
			args: args{
				f: sqb.Ne,
				m: map[string]interface{}{
					"col":    2,
					"colcol": time.Time{},
				},
			},
			want:     "col != ? AND colcol != ?",
			wantArgs: []interface{}{2, time.Time{}},
		},
		{
			name: "three arguments",
			args: args{
				f: sqb.Eq,
				m: map[string]interface{}{
					"col":       2,
					"colcol":    time.Time{},
					"colcolcol": []byte("hello"),
				},
			},
			want:     "col = ? AND colcol = ? AND colcolcol = ?",
			wantArgs: []interface{}{2, time.Time{}, []byte("hello")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BuildCapture{
				buf:  strings.Builder{},
				Args: []interface{}{},
			}
			expr := sqb.AndFromMap(tt.args.f, tt.args.m)
			if err := expr.Write(b); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := b.buf.String(); tt.want != got {
				t.Errorf("\nwant: %q\ngot: %q", tt.want, got)
			}
			if diff := cmp.Diff(tt.wantArgs, b.Args); diff != "" {
				t.Errorf("args (-want, +got)\n%s", diff)
			}
		})
	}
}
func TestAndByMap_nil_panic(t *testing.T) {
	got := sqb.AndFromMap(sqb.Lt, map[string]interface{}{})
	if got != nil {
		t.Fatalf("expected nil")
	}
	defer func() {
		if v := recover(); v == nil {
			panic("expected panic")
		}
	}()
	sqb.AndFromMap(nil, map[string]interface{}{})
}

func TestOrFromMap(t *testing.T) {
	type args struct {
		f sqb.ConditionalFunc
		m map[string]interface{}
	}
	tests := []struct {
		name     string
		args     args
		want     string
		wantArgs []interface{}
	}{
		{
			name: "an argument",
			args: args{
				f: sqb.Gt,
				m: map[string]interface{}{
					"col": 2,
				},
			},
			want:     "col > ?",
			wantArgs: []interface{}{2},
		},
		{
			name: "two arguments",
			args: args{
				f: sqb.Ne,
				m: map[string]interface{}{
					"col":    2,
					"colcol": time.Time{},
				},
			},
			want:     "(col != ? OR colcol != ?)",
			wantArgs: []interface{}{2, time.Time{}},
		},
		{
			name: "three arguments",
			args: args{
				f: sqb.Eq,
				m: map[string]interface{}{
					"col":       2,
					"colcol":    time.Time{},
					"colcolcol": []byte("hello"),
				},
			},
			want:     "((col = ? OR colcol = ?) OR colcolcol = ?)",
			wantArgs: []interface{}{2, time.Time{}, []byte("hello")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BuildCapture{
				buf:  strings.Builder{},
				Args: []interface{}{},
			}
			expr := sqb.OrFromMap(tt.args.f, tt.args.m)
			if err := expr.Write(b); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := b.buf.String(); tt.want != got {
				t.Errorf("\nwant: %q\ngot: %q", tt.want, got)
			}
			if diff := cmp.Diff(tt.wantArgs, b.Args); diff != "" {
				t.Errorf("args (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestOrByMap_nil_panic(t *testing.T) {
	got := sqb.OrFromMap(sqb.Lt, map[string]interface{}{})
	if got != nil {
		t.Fatalf("expected nil")
	}
	defer func() {
		if v := recover(); v == nil {
			panic("expected panic")
		}
	}()
	sqb.OrFromMap(nil, map[string]interface{}{})
}

func TestParen(t *testing.T) {
	tests := []struct {
		name     string
		args     stmt.Expr
		want     string
		wantArgs []interface{}
	}{
		{
			name:     "valid",
			args:     sqb.Eq("col", true),
			want:     "(col = ?)",
			wantArgs: []interface{}{true},
		},
		{
			name: "valid AND",
			args: sqb.And(
				sqb.Eq("col", true),
				sqb.Ne("col2", 10),
			),
			want:     "(col = ? AND col2 != ?)",
			wantArgs: []interface{}{true, 10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BuildCapture{
				buf:  strings.Builder{},
				Args: []interface{}{},
			}
			expr := sqb.Paren(tt.args)
			if err := expr.Write(b); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := b.buf.String(); tt.want != got {
				t.Errorf("\nwant: %q\ngot: %q", tt.want, got)
			}
			if diff := cmp.Diff(tt.wantArgs, b.Args); diff != "" {
				t.Errorf("args (-want, +got)\n%s", diff)
			}
		})
	}
}
