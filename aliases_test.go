package sqb_test

import (
	"testing"

	"github.com/Code-Hex/sqb"
)

func TestColumns(t *testing.T) {
	tests := []struct {
		name    string
		c       sqb.Columns
		want    string
		wantErr bool
	}{
		{
			name:    "valid",
			c:       sqb.Columns{"hello"},
			want:    "SELECT hello FROM table",
			wantErr: false,
		},
		{
			name:    "valid columns",
			c:       sqb.Columns{"hello", "world", "sqb"},
			want:    "SELECT hello, world, sqb FROM table",
			wantErr: false,
		},
		{
			name:    "invalid",
			c:       sqb.Columns{},
			want:    "",
			wantErr: true,
		},
	}
	const sqlstr = "SELECT ? FROM table"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := sqb.New(sqlstr).Bind(tt.c)
			got, _, err := builder.Build()
			if (err != nil) != tt.wantErr {
				t.Errorf("Columns error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.want != got {
				t.Errorf("\nwant: %q\ngot: %q", tt.want, got)
			}
		})
	}
}

func TestTable(t *testing.T) {
	tests := []struct {
		name    string
		t       sqb.Table
		want    string
		wantErr bool
	}{
		{
			name:    "valid",
			t:       sqb.Table("hello"),
			want:    "SELECT * FROM hello",
			wantErr: false,
		},
		{
			name:    "invalid",
			t:       sqb.Table(""),
			want:    "",
			wantErr: true,
		},
	}
	const sqlstr = "SELECT * FROM ?"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := sqb.New(sqlstr).Bind(tt.t)
			got, _, err := builder.Build()
			if (err != nil) != tt.wantErr {
				t.Errorf("From error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.want != got {
				t.Errorf("\nwant: %q\ngot: %q", tt.want, got)
			}
		})
	}
}

func TestLimit(t *testing.T) {
	tests := []struct {
		name    string
		l       sqb.Limit
		want    string
		wantErr bool
	}{
		{
			name:    "valid",
			l:       sqb.Limit(100),
			want:    "SELECT * FROM table LIMIT 100",
			wantErr: false,
		},
		{
			name:    "valid 0",
			l:       sqb.Limit(0),
			want:    "SELECT * FROM table LIMIT 0",
			wantErr: false,
		},
	}
	const sqlstr = "SELECT * FROM table ?"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := sqb.New(sqlstr).Bind(tt.l)
			got, _, err := builder.Build()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.want != got {
				t.Errorf("\nwant: %q\ngot: %q", tt.want, got)
			}
		})
	}
}

func TestOffset(t *testing.T) {
	tests := []struct {
		name    string
		o       sqb.Offset
		want    string
		wantErr bool
	}{
		{
			name:    "valid",
			o:       sqb.Offset(100),
			want:    "SELECT * FROM table LIMIT 1 OFFSET 100",
			wantErr: false,
		},
		{
			name:    "valid 0",
			o:       sqb.Offset(0),
			want:    "SELECT * FROM table LIMIT 1 OFFSET 0",
			wantErr: false,
		},
	}
	const sqlstr = "SELECT * FROM table LIMIT 1 ?"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := sqb.New(sqlstr).Bind(tt.o)
			got, _, err := builder.Build()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.want != got {
				t.Errorf("\nwant: %q\ngot: %q", tt.want, got)
			}
		})
	}
}
