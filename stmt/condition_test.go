package stmt

import (
	"errors"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

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
				Column: "name",
				Compare: &CompOp{
					Op:    "=",
					Value: "taro",
				},
			},
			want:     "name = ?",
			wantArgs: []interface{}{"taro"},
			wantErr:  false,
		},
		{
			name: "invalid nil compare",
			c: &Condition{
				Column:  "name",
				Compare: nil,
			},
			want:     "",
			wantArgs: []interface{}{},
			wantErr:  true,
		},
		{
			name: "invalid write compare",
			c: &Condition{
				Column: "name",
				Compare: &ComparisonerMock{
					WriteComparisonMock: func(Builder) error {
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
