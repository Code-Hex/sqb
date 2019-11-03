package stmt

import (
	"errors"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWhere_Write(t *testing.T) {
	tests := []struct {
		name     string
		where    *Where
		want     string
		wantArgs []interface{}
		wantErr  bool
	}{
		{
			name: "invalid",
			where: &Where{
				Expr: nil,
			},
			want:     "",
			wantArgs: []interface{}{},
			wantErr:  true,
		},
		{
			name: "valid simple",
			where: &Where{
				Expr: &Condition{
					Column: "category",
					Compare: &CompOp{
						Op:    "=",
						Value: 10,
					},
				},
			},
			want:     "WHERE category = ?",
			wantArgs: []interface{}{10},
			wantErr:  false,
		},
		{
			name: "valid conject",
			where: &Where{
				Expr: &Conjection{
					Left: &Condition{
						Column: "category",
						Compare: &CompOp{
							Op:    "=",
							Value: 10,
						},
					},
					Combined: "AND",
					Right: &Condition{
						Column: "brand",
						Compare: &CompOp{
							Op:    "!=",
							Value: "apple",
						},
					},
				},
			},
			want:     "WHERE category = ? AND brand != ?",
			wantArgs: []interface{}{10, "apple"},
			wantErr:  false,
		},
		{
			name: "valid conject twice",
			where: &Where{
				Expr: &Conjection{
					Left: &Conjection{
						Left: &Condition{
							Column: "category",
							Compare: &CompOp{
								Op:    "=",
								Value: 1,
							},
						},
						Combined: "OR",
						Right: &Condition{
							Column: "category",
							Compare: &CompOp{
								Op:    "=",
								Value: 2,
							},
						},
					},
					Combined: "AND",
					Right: &Conjection{
						Left: &Condition{
							Column: "brand",
							Compare: &CompIn{
								Negative: true,
								Values: []interface{}{
									"apple", "sony", "google",
								},
							},
						},
						Combined: "OR",
						Right: &Condition{
							Column: "name",
							Compare: &CompLike{
								Negative: true,
								Value:    "abc%",
							},
						},
					},
				},
			},
			want:     "WHERE (category = ? OR category = ?) AND (brand NOT IN (?, ?, ?) OR name NOT LIKE ?)",
			wantArgs: []interface{}{1, 2, "apple", "sony", "google", "abc%"},
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BuildCapture{
				buf:  strings.Builder{},
				Args: []interface{}{},
			}
			if err := tt.where.Write(b); (err != nil) != tt.wantErr {
				t.Errorf("Where.Write() error = %v, wantErr %v", err, tt.wantErr)
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

func TestConjection_Write(t *testing.T) {
	tests := []struct {
		name     string
		c        *Conjection
		want     string
		wantArgs []interface{}
		wantErr  bool
	}{
		{
			name: "valid",
			c: &Conjection{
				Left: &Condition{
					Column: "category",
					Compare: &CompOp{
						Op:    "=",
						Value: "music",
					},
				},
				Combined: "AND",
				Right: &Condition{
					Column: "id",
					Compare: &CompOp{
						Op:    "<",
						Value: 3,
					},
				},
			},
			want:     "category = ? AND id < ?",
			wantArgs: []interface{}{"music", 3},
			wantErr:  false,
		},
		{
			name: "valid nested",
			c: &Conjection{
				Left: &Conjection{
					Left: &Condition{
						Column: "category",
						Compare: &CompOp{
							Op:    "=",
							Value: "music",
						},
					},
					Combined: "OR",
					Right: &Condition{
						Column: "category",
						Compare: &CompOp{
							Op:    "=",
							Value: "video",
						},
					},
				},
				Combined: "AND",
				Right: &Condition{
					Column: "id",
					Compare: &CompOp{
						Op:    "<",
						Value: 3,
					},
				},
			},
			want:     "(category = ? OR category = ?) AND id < ?",
			wantArgs: []interface{}{"music", "video", 3},
			wantErr:  false,
		},
		{
			name: "valid nested in nested",
			c: &Conjection{
				Left: &Conjection{
					Left: &Condition{
						Column: "category",
						Compare: &CompOp{
							Op:    "=",
							Value: "music",
						},
					},
					Combined: "OR",
					Right: &Conjection{
						Left: &Condition{
							Column: "sub_category",
							Compare: &CompOp{
								Op:    "=",
								Value: "jpop",
							},
						},
						Combined: "AND",
						Right: &Condition{
							Column: "sub_category",
							Compare: &CompOp{
								Op:    "=",
								Value: "hiphop",
							},
						},
					},
				},
				Combined: "AND",
				Right: &Condition{
					Column: "id",
					Compare: &CompOp{
						Op:    "<",
						Value: 3,
					},
				},
			},
			want:     "(category = ? OR (sub_category = ? AND sub_category = ?)) AND id < ?",
			wantArgs: []interface{}{"music", "jpop", "hiphop", 3},
			wantErr:  false,
		},
		{
			name: "invalid left",
			c: &Conjection{
				Left:     nil,
				Combined: "AND",
				Right:    &Condition{},
			},
			want:     "",
			wantArgs: []interface{}{},
			wantErr:  true,
		},
		{
			name: "invalid left conjection",
			c: &Conjection{
				Left:     &Conjection{},
				Combined: "AND",
				Right:    &Condition{},
			},
			want:     "",
			wantArgs: []interface{}{},
			wantErr:  true,
		},
		{
			name: "invalid left write",
			c: &Conjection{
				Left: &ExprMock{
					WriteMock: func(Builder) error {
						return errors.New("error")
					},
				},
				Combined: "AND",
				Right:    &Condition{},
			},
			want:     "",
			wantArgs: []interface{}{},
			wantErr:  true,
		},
		{
			name: "valid unset combined",
			c: &Conjection{
				Left: &Condition{
					Column: "category",
					Compare: &CompOp{
						Op:    "=",
						Value: "music",
					},
				},
				Combined: "",
				Right:    &Condition{},
			},
			want:     "category = ?",
			wantArgs: []interface{}{"music"},
			wantErr:  false,
		},
		{
			name: "invalid right",
			c: &Conjection{
				Left: &Condition{
					Column: "category",
					Compare: &CompOp{
						Op:    "=",
						Value: "music",
					},
				},
				Combined: "AND",
				Right:    nil,
			},
			want:     "",
			wantArgs: []interface{}{},
			wantErr:  true,
		},
		{
			name: "invalid right conjection",
			c: &Conjection{
				Left: &Condition{
					Column: "category",
					Compare: &CompOp{
						Op:    "=",
						Value: "music",
					},
				},
				Combined: "AND",
				Right:    &Conjection{},
			},
			want:     "",
			wantArgs: []interface{}{},
			wantErr:  true,
		},
		{
			name: "invalid right write",
			c: &Conjection{
				Left: &Condition{
					Column: "category",
					Compare: &CompOp{
						Op:    "=",
						Value: "music",
					},
				},
				Combined: "AND",
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
			if err := tt.c.Write(b); (err != nil) != tt.wantErr {
				t.Errorf("Conjection.Write() error = %v, wantErr %v", err, tt.wantErr)
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

func TestCondition_Write(t *testing.T) {
	tests := []struct {
		name     string
		c        *Condition
		want     string
		wantArgs []interface{}
		wantErr  bool
	}{
		{
			name: "valid",
			c: &Condition{
				Column: "category",
				Compare: &CompOp{
					Op:    "=",
					Value: 1,
				},
			},
			want:     "category = ?",
			wantArgs: []interface{}{1},
			wantErr:  false,
		},
		{
			name: "invalid",
			c: &Condition{
				Column:  "category",
				Compare: nil,
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
			if err := tt.c.Write(b); (err != nil) != tt.wantErr {
				t.Errorf("Condition.Write() error = %v, wantErr %v", err, tt.wantErr)
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

func Test_makePlaceholders(t *testing.T) {
	makeArgs := func(i int) []interface{} {
		ret := make([]interface{}, i)
		for idx := range ret {
			ret[idx] = 0
		}
		return ret
	}
	tests := []struct {
		name    string
		args    int
		want    string
		wantErr bool
	}{
		{
			name:    "invalid 0 arguments",
			args:    0,
			wantErr: true,
		},
		{
			name:    "valid 1 arguments",
			args:    1,
			want:    "?",
			wantErr: false,
		},
		{
			name:    "valid 2 arguments",
			args:    2,
			want:    "?, ?",
			wantErr: false,
		},
		{
			name:    "valid 5 arguments",
			args:    5,
			want:    "?, ?, ?, ?, ?",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BuildCapture{
				buf:  strings.Builder{},
				Args: []interface{}{},
			}
			args := makeArgs(tt.args)
			if err := makePlaceholders(b, args); (err != nil) != tt.wantErr {
				t.Errorf("makePlaceholders() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got := b.buf.String(); tt.want != got {
				t.Errorf("\nwant: %q\ngot: %q", tt.want, got)
			}
		})
	}
}
