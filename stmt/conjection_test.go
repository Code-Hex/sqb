package stmt

import (
	"errors"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParen_Write(t *testing.T) {
	tests := []struct {
		name     string
		p        *Paren
		want     string
		wantArgs []interface{}
		wantErr  bool
	}{
		{
			name: "valid",
			p: &Paren{
				Expr: &Condition{
					Column: "hello",
					Compare: &CompOp{
						Op:    "=",
						Value: 10,
					},
				},
			},
			want:     "(hello = ?)",
			wantArgs: []interface{}{10},
			wantErr:  false,
		},
		{
			name: "invalid nil",
			p: &Paren{
				Expr: nil,
			},
			want:     "",
			wantArgs: []interface{}{},
			wantErr:  true,
		},
		{
			name: "invalid",
			p: &Paren{
				Expr: &ExprMock{
					WriteMock: func(b Builder) error {
						return errors.New("error")
					},
				},
			},
			want:     "",
			wantArgs: []interface{}{},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BuildCapture{
				buf:  strings.Builder{},
				Args: []interface{}{},
			}
			if err := tt.p.Write(b); (err != nil) != tt.wantErr {
				t.Errorf("Paren.Write() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if got := b.buf.String(); tt.want != got {
					t.Errorf("\nwant: %q\ngot: %q", tt.want, got)
				}
				if diff := cmp.Diff(tt.wantArgs, b.Args); diff != "" {
					t.Errorf("args (-want, +got)\n%s", diff)
				}
			}
		})
	}
}

func TestOr_Write(t *testing.T) {
	tests := []struct {
		name     string
		o        *Or
		want     string
		wantArgs []interface{}
		wantErr  bool
	}{
		{
			name: "valid",
			o: &Or{
				Left: &Condition{
					Column: "hello",
					Compare: &CompLike{
						Negative: false,
						Value:    "world",
					},
				},
				Right: &Condition{
					Column: "world",
					Compare: &CompBetween{
						Negative: false,
						Left:     10,
						Right:    300,
					},
				},
			},
			want: "(hello LIKE ? OR world BETWEEN ? AND ?)",
			wantArgs: []interface{}{
				"world",
				10, 300,
			},
			wantErr: false,
		},
		{
			name: "invalid Left is nil",
			o: &Or{
				Left: nil,
			},
			want:     "",
			wantArgs: []interface{}{},
			wantErr:  true,
		},
		{
			name: "invalid Left",
			o: &Or{
				Left: &ExprMock{
					WriteMock: func(Builder) error {
						return errors.New("error")
					},
				},
				Right: &ExprMock{
					WriteMock: func(Builder) error {
						return nil
					},
				},
			},
			want:     "",
			wantArgs: []interface{}{},
			wantErr:  true,
		},
		{
			name: "invalid Right is nil",
			o: &Or{
				Left:  &ExprMock{},
				Right: nil,
			},
			want:     "",
			wantArgs: []interface{}{},
			wantErr:  true,
		},
		{
			name: "invalid Right",
			o: &Or{
				Left: &ExprMock{
					WriteMock: func(Builder) error {
						return nil
					},
				},
				Right: &ExprMock{
					WriteMock: func(Builder) error {
						return errors.New("error")
					},
				},
			},
			want:     "",
			wantArgs: []interface{}{},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BuildCapture{
				buf:  strings.Builder{},
				Args: []interface{}{},
			}
			if err := tt.o.Write(b); (err != nil) != tt.wantErr {
				t.Errorf("Or.Write() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if got := b.buf.String(); tt.want != got {
					t.Errorf("\nwant: %q\ngot: %q", tt.want, got)
				}
				if diff := cmp.Diff(tt.wantArgs, b.Args); diff != "" {
					t.Errorf("args (-want, +got)\n%s", diff)
				}
			}
		})
	}
}

func TestAnd_Write(t *testing.T) {
	tests := []struct {
		name     string
		a        *And
		want     string
		wantArgs []interface{}
		wantErr  bool
	}{
		{
			name: "valid",
			a: &And{
				Left: &Condition{
					Column: "hello",
					Compare: &CompLike{
						Negative: false,
						Value:    "world",
					},
				},
				Right: &Condition{
					Column: "world",
					Compare: &CompBetween{
						Negative: false,
						Left:     10,
						Right:    300,
					},
				},
			},
			want: "hello LIKE ? AND world BETWEEN ? AND ?",
			wantArgs: []interface{}{
				"world",
				10, 300,
			},
			wantErr: false,
		},
		{
			name: "invalid Left is nil",
			a: &And{
				Left: nil,
			},
			want:     "",
			wantArgs: []interface{}{},
			wantErr:  true,
		},
		{
			name: "invalid Left",
			a: &And{
				Left: &ExprMock{
					WriteMock: func(Builder) error {
						return errors.New("error")
					},
				},
				Right: &ExprMock{
					WriteMock: func(Builder) error {
						return nil
					},
				},
			},
			want:     "",
			wantArgs: []interface{}{},
			wantErr:  true,
		},
		{
			name: "invalid Right is nil",
			a: &And{
				Left:  &ExprMock{},
				Right: nil,
			},
			want:     "",
			wantArgs: []interface{}{},
			wantErr:  true,
		},
		{
			name: "invalid Right",
			a: &And{
				Left: &ExprMock{
					WriteMock: func(Builder) error {
						return nil
					},
				},
				Right: &ExprMock{
					WriteMock: func(Builder) error {
						return errors.New("error")
					},
				},
			},
			want:     "",
			wantArgs: []interface{}{},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BuildCapture{
				buf:  strings.Builder{},
				Args: []interface{}{},
			}
			if err := tt.a.Write(b); (err != nil) != tt.wantErr {
				t.Errorf("And.Write() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if got := b.buf.String(); tt.want != got {
					t.Errorf("\nwant: %q\ngot: %q", tt.want, got)
				}
				if diff := cmp.Diff(tt.wantArgs, b.Args); diff != "" {
					t.Errorf("args (-want, +got)\n%s", diff)
				}
			}
		})
	}
}
