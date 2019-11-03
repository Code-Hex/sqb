package stmt

import (
	"errors"
	"testing"
)

func TestBuildError_Unwrap(t *testing.T) {
	e := &BuildError{
		Err: errors.New("error"),
	}
	if got := e.Unwrap(); e.Err != got {
		t.Fatalf("want %v, but got %v", e.Err, got)
	}
}

func TestBuildError_Error(t *testing.T) {
	tests := []struct {
		name string
		b    *BuildError
		want string
	}{
		{
			name: "nil",
			b:    nil,
			want: "<nil>",
		},
		{
			name: "nil",
			b: &BuildError{
				Op:  "ope",
				Err: errors.New("error"),
			},
			want: "ope: error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Error(); got != tt.want {
				t.Errorf("BuildError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
