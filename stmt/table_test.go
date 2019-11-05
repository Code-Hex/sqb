package stmt

import (
	"strings"
	"testing"
)

func TestTable_Write(t *testing.T) {
	tests := []struct {
		name    string
		t       Table
		want    string
		wantErr bool
	}{
		{
			name:    "valid",
			t:       Table("table"),
			want:    "table",
			wantErr: false,
		},
		{
			name:    "invalid",
			t:       "",
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BuildCapture{
				buf:  strings.Builder{},
				Args: []interface{}{},
			}
			if err := tt.t.Write(b); (err != nil) != tt.wantErr {
				t.Errorf("Table.Write() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if got := b.buf.String(); tt.want != got {
					t.Errorf("\nwant: %q\ngot: %q", tt.want, got)
				}
			}
		})
	}
}

func TestLimit_Write(t *testing.T) {
	tests := []struct {
		name string
		l    Limit
		want string
	}{
		{
			name: "valid",
			l:    Limit(10),
			want: "LIMIT 10",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BuildCapture{
				buf:  strings.Builder{},
				Args: []interface{}{},
			}
			if err := tt.l.Write(b); err != nil {
				t.Fatalf("Limit.Write() unexpected error: %v", err)
			}
			if got := b.buf.String(); tt.want != got {
				t.Errorf("\nwant: %q\ngot: %q", tt.want, got)
			}
		})
	}
}

func TestOffset_Write(t *testing.T) {
	tests := []struct {
		name string
		o    Offset
		want string
	}{
		{
			name: "valid",
			o:    Offset(10),
			want: "OFFSET 10",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BuildCapture{
				buf:  strings.Builder{},
				Args: []interface{}{},
			}
			if err := tt.o.Write(b); err != nil {
				t.Fatalf("Offset.Write() unexpected error: %v", err)
			}
			if got := b.buf.String(); tt.want != got {
				t.Errorf("\nwant: %q\ngot: %q", tt.want, got)
			}
		})
	}
}

func TestOrderBy_Write(t *testing.T) {
	tests := []struct {
		name string
		o    *OrderBy
		want string
	}{
		{
			name: "valid unary ASC",
			o: &OrderBy{
				Column: "column",
				Desc:   false,
				Next:   nil,
			},
			want: "column",
		},
		{
			name: "valid unary DESC",
			o: &OrderBy{
				Column: "column",
				Desc:   true,
				Next:   nil,
			},
			want: "column DESC",
		},
		{
			name: "valid has Next",
			o: &OrderBy{
				Column: "column1",
				Desc:   false,
				Next: &OrderBy{
					Column: "column2",
					Desc:   true,
					Next:   nil,
				},
			},
			want: "column1, column2 DESC",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BuildCapture{
				buf:  strings.Builder{},
				Args: []interface{}{},
			}
			if err := tt.o.Write(b); err != nil {
				t.Fatalf("From.Write() unexpected error: %v", err)
			}
			if got := b.buf.String(); tt.want != got {
				t.Errorf("\nwant: %q\ngot: %q", tt.want, got)
			}
		})
	}
}
