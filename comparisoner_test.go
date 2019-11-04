package sqb_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Code-Hex/sqb"
	"github.com/google/go-cmp/cmp"
)

func TestConditional(t *testing.T) {
	tests := []struct {
		f    sqb.ConditionalFunc
		want string
	}{
		{
			f:    sqb.Eq,
			want: "=",
		},
		{
			f:    sqb.Ge,
			want: ">=",
		},
		{
			f:    sqb.Gt,
			want: ">",
		},
		{
			f:    sqb.Le,
			want: "<=",
		},
		{
			f:    sqb.Lt,
			want: "<",
		},
		{
			f:    sqb.Ne,
			want: "!=",
		},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			b := &BuildCapture{
				buf:  strings.Builder{},
				Args: []interface{}{},
			}
			wantArg := 1
			expr := tt.f("col", wantArg)
			if err := expr.Write(b); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			wantQuery := fmt.Sprintf("col %s ?", tt.want)
			if got := b.buf.String(); wantQuery != got {
				t.Errorf("\nwant: %q\ngot: %q", wantQuery, got)
			}
			if diff := cmp.Diff([]interface{}{wantArg}, b.Args); diff != "" {
				t.Errorf("args (-want, +got)\n%s", diff)
			}
		})
	}
}
func TestLike(t *testing.T) {
	type args struct {
		column string
		value  interface{}
	}
	tests := []struct {
		name     string
		args     args
		want     string
		wantArgs []interface{}
	}{
		{
			name: "valid",
			args: args{
				column: "col",
				value:  "abc%",
			},
			want:     "col LIKE ?",
			wantArgs: []interface{}{"abc%"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BuildCapture{
				buf:  strings.Builder{},
				Args: []interface{}{},
			}
			expr := sqb.Like(tt.args.column, tt.args.value)
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

func TestNotLike(t *testing.T) {
	type args struct {
		column string
		value  interface{}
	}
	tests := []struct {
		name     string
		args     args
		want     string
		wantArgs []interface{}
	}{
		{
			name: "valid",
			args: args{
				column: "col",
				value:  "abc%",
			},
			want:     "col NOT LIKE ?",
			wantArgs: []interface{}{"abc%"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BuildCapture{
				buf:  strings.Builder{},
				Args: []interface{}{},
			}
			expr := sqb.NotLike(tt.args.column, tt.args.value)
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

func TestIn(t *testing.T) {
	type args struct {
		column string
		args   []interface{}
	}
	tests := []struct {
		name     string
		args     args
		want     string
		wantArgs []interface{}
	}{
		{
			name: "valid",
			args: args{
				column: "col",
				args: []interface{}{
					[]uint8("hello"),
				},
			},
			want: "col IN (?)",
			wantArgs: []interface{}{
				[]uint8("hello"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BuildCapture{
				buf:  strings.Builder{},
				Args: []interface{}{},
			}
			expr := sqb.In(tt.args.column, tt.args.args...)
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

func TestNotIn(t *testing.T) {
	type args struct {
		column string
		args   []interface{}
	}
	tests := []struct {
		name     string
		args     args
		want     string
		wantArgs []interface{}
	}{
		{
			name: "valid",
			args: args{
				column: "col",
				args: []interface{}{
					[]uint8("hello"),
				},
			},
			want: "col NOT IN (?)",
			wantArgs: []interface{}{
				[]uint8("hello"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BuildCapture{
				buf:  strings.Builder{},
				Args: []interface{}{},
			}
			expr := sqb.NotIn(tt.args.column, tt.args.args...)
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

func TestBetween(t *testing.T) {
	type args struct {
		column string
		left   interface{}
		right  interface{}
	}
	tests := []struct {
		name     string
		args     args
		want     string
		wantArgs []interface{}
	}{
		{
			name: "valid",
			args: args{
				column: "col",
				left:   1,
				right:  2,
			},
			want: "col BETWEEN ? AND ?",
			wantArgs: []interface{}{
				1, 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BuildCapture{
				buf:  strings.Builder{},
				Args: []interface{}{},
			}
			expr := sqb.Between(tt.args.column, tt.args.left, tt.args.right)
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

func TestNotBetween(t *testing.T) {
	type args struct {
		column string
		left   interface{}
		right  interface{}
	}
	tests := []struct {
		name     string
		args     args
		want     string
		wantArgs []interface{}
	}{
		{
			name: "valid",
			args: args{
				column: "col",
				left:   1,
				right:  2,
			},
			want: "col NOT BETWEEN ? AND ?",
			wantArgs: []interface{}{
				1, 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BuildCapture{
				buf:  strings.Builder{},
				Args: []interface{}{},
			}
			expr := sqb.NotBetween(tt.args.column, tt.args.left, tt.args.right)
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
