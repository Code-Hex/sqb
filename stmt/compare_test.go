package stmt

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

// "=", ">=", ">", "<=", "<", "!=", "IS", "IS NOT"
func TestCompOp_WriteComparison(t *testing.T) {
	tests := []string{
		"=", ">=", ">", "<=", "<", "!=", "IS", "IS NOT",
	}
	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			c := &CompOp{
				Op:    tt,
				Value: 1,
			}
			b := &BuildCapture{
				buf:  strings.Builder{},
				Args: []interface{}{},
			}
			if err := c.WriteComparison(b); err != nil {
				t.Fatalf("CompOp.WriteComparison() error = %v", err)
			}
			got := b.buf.String()
			want := fmt.Sprintf("%s ?", tt) // = ?, >= ?, IS ?, etc...
			if want != got {
				t.Errorf("\nwant: %q\ngot: %q", want, got)
			}

			wantArgs := []interface{}{1}
			if diff := cmp.Diff(wantArgs, b.Args); diff != "" {
				t.Errorf("args (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestCompLike_WriteComparison(t *testing.T) {
	tests := []struct {
		name     string
		c        *CompLike
		want     string
		wantArgs []interface{}
	}{
		{
			name: "LIKE",
			c: &CompLike{
				Negative: false,
				Value:    "abc%",
			},
			want:     "LIKE ?",
			wantArgs: []interface{}{"abc%"},
		},
		{
			name: "NOT LIKE",
			c: &CompLike{
				Negative: true,
				Value:    "abc%",
			},
			want:     "NOT LIKE ?",
			wantArgs: []interface{}{"abc%"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BuildCapture{
				buf:  strings.Builder{},
				Args: []interface{}{},
			}
			if err := tt.c.WriteComparison(b); err != nil {
				t.Fatalf("CompLike.WriteComparison() error = %v", err)
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

func TestCompBetween_WriteComparison(t *testing.T) {
	tests := []struct {
		name     string
		c        *CompBetween
		want     string
		wantArgs []interface{}
		wantErr  bool
	}{
		{
			name: "valid BETWEEN",
			c: &CompBetween{
				Negative: false,
				Left:     100,
				Right:    200,
			},
			want:     "BETWEEN ? AND ?",
			wantArgs: []interface{}{100, 200},
			wantErr:  false,
		},
		{
			name: "valid NOT BETWEEN",
			c: &CompBetween{
				Negative: true,
				Left:     time.Date(2014, 2, 1, 0, 0, 0, 0, time.UTC),
				Right:    time.Date(2014, 2, 28, 0, 0, 0, 0, time.UTC),
			},
			want: "NOT BETWEEN ? AND ?",
			wantArgs: []interface{}{
				time.Date(2014, 2, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2014, 2, 28, 0, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name: "invalid Left value",
			c: &CompBetween{
				Negative: false,
				Left:     nil,
				Right:    200,
			},
			want:     "",
			wantArgs: []interface{}{},
			wantErr:  true,
		},
		{
			name: "invalid Right value",
			c: &CompBetween{
				Negative: true,
				Left:     100,
				Right:    nil,
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
			if err := tt.c.WriteComparison(b); (err != nil) != tt.wantErr {
				t.Errorf("CompBetween.WriteComparison() error = %v, wantErr %v", err, tt.wantErr)
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

func TestCompIn_WriteComparison(t *testing.T) {
	tests := []struct {
		name     string
		c        *CompIn
		want     string
		wantArgs []interface{}
		wantErr  bool
	}{
		{
			name: "valid IN",
			c: &CompIn{
				Negative: false,
				Values:   []interface{}{1},
			},
			want:     "IN (?)",
			wantArgs: []interface{}{1},
			wantErr:  false,
		},
		{
			name: "valid nested list",
			c: &CompIn{
				Negative: false,
				Values: []interface{}{
					[]int{1, 2},
					[]interface{}{
						100,
						[]interface{}{"hello", []byte("world")},
						500,
					},
				},
			},
			want: "IN (?, ?, ?, ?, ?, ?)",
			wantArgs: []interface{}{
				1, 2,
				100,
				"hello", []byte("world"),
				500,
			},
			wantErr: false,
		},
		{
			name: "valid NOT IN",
			c: &CompIn{
				Negative: true,
				Values:   []interface{}{1, 2, 3},
			},
			want:     "NOT IN (?, ?, ?)",
			wantArgs: []interface{}{1, 2, 3},
			wantErr:  false,
		},
		{
			name: "invalid",
			c: &CompIn{
				Negative: false,
				Values:   []interface{}{},
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
			if err := tt.c.WriteComparison(b); (err != nil) != tt.wantErr {
				t.Errorf("CompIn.WriteComparison() error = %v, wantErr %v", err, tt.wantErr)
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
