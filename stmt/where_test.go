package stmt

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

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

func TestWhere_Write(t *testing.T) {
	tests := []struct {
		name     string
		where    *Where
		want     string
		wantArgs []interface{}
		wantErr  bool
	}{
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
			if got := b.buf.String(); tt.want != got {
				t.Errorf("\nwant: %q\ngot: %q", tt.want, got)
			}
			if diff := cmp.Diff(tt.wantArgs, b.Args); diff != "" {
				t.Errorf("args (-want, +got)\n%s", diff)
			}
		})
	}
}
