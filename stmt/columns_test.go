package stmt

import (
	"strings"
	"testing"
)

func TestColumns_Write(t *testing.T) {
	tests := []struct {
		name    string
		c       Columns
		want    string
		wantErr bool
	}{
		{
			name:    "invalid no column",
			c:       Columns{},
			want:    "",
			wantErr: true,
		},
		{
			name:    "valid a column",
			c:       Columns{"col"},
			want:    "col",
			wantErr: false,
		},
		{
			name:    "valid some columns",
			c:       Columns{"col1", "col2", "col3"},
			want:    "col1, col2, col3",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BuildCapture{
				buf:  strings.Builder{},
				Args: []interface{}{},
			}
			if err := tt.c.Write(b); (err != nil) != tt.wantErr {
				t.Errorf("Columns.Write() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if got := b.buf.String(); tt.want != got {
					t.Errorf("\nwant: %q\ngot: %q", tt.want, got)
				}
			}
		})
	}
}
