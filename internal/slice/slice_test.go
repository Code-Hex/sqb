package slice

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestFlatten(t *testing.T) {
	tests := []struct {
		name string
		args []interface{}
		want []interface{}
	}{
		{
			name: "normal",
			args: []interface{}{
				"hello",
				1234,
				time.Date(2019, 11, 4, 0, 0, 0, 0, time.UTC),
			},
			want: []interface{}{
				"hello",
				1234,
				time.Date(2019, 11, 4, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "contain []byte",
			args: []interface{}{
				"hello",
				1234,
				[]byte("hello"),
				time.Date(2019, 11, 4, 0, 0, 0, 0, time.UTC),
			},
			want: []interface{}{
				"hello",
				1234,
				[]byte("hello"),
				time.Date(2019, 11, 4, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "contain []int8",
			args: []interface{}{
				"hello",
				1234,
				[]uint8("hello"),
				time.Date(2019, 11, 4, 0, 0, 0, 0, time.UTC),
			},
			want: []interface{}{
				"hello",
				1234,
				[]uint8("hello"),
				time.Date(2019, 11, 4, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "contain []int",
			args: []interface{}{
				"hello",
				1234,
				[]int{1, 2, 3},
				time.Date(2019, 11, 4, 0, 0, 0, 0, time.UTC),
			},
			want: []interface{}{
				"hello",
				1234,
				1, 2, 3,
				time.Date(2019, 11, 4, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "contain []string",
			args: []interface{}{
				"hello",
				1234,
				[]string{"a", "b", "c"},
				time.Date(2019, 11, 4, 0, 0, 0, 0, time.UTC),
			},
			want: []interface{}{
				"hello",
				1234,
				"a", "b", "c",
				time.Date(2019, 11, 4, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "contain []int64",
			args: []interface{}{
				"hello",
				1234,
				[]int64{1, 2, 3},
				time.Date(2019, 11, 4, 0, 0, 0, 0, time.UTC),
			},
			want: []interface{}{
				"hello",
				1234,
				int64(1), int64(2), int64(3),
				time.Date(2019, 11, 4, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "contain []interface{}",
			args: []interface{}{
				"hello",
				1234,
				[]interface{}{"a", 1, "b", 2},
				time.Date(2019, 11, 4, 0, 0, 0, 0, time.UTC),
			},
			want: []interface{}{
				"hello",
				1234,
				"a", 1, "b", 2,
				time.Date(2019, 11, 4, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "contain [][]int",
			args: []interface{}{
				"hello",
				1234,
				[][]int{
					[]int{1, 2, 3},
					[]int{5, 6, 7},
				},
				time.Date(2019, 11, 4, 0, 0, 0, 0, time.UTC),
			},
			want: []interface{}{
				"hello",
				1234,
				1, 2, 3,
				5, 6, 7,
				time.Date(2019, 11, 4, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "contain [][]interface{}",
			args: []interface{}{
				"hello",
				1234,
				[][]interface{}{
					[]interface{}{"a", 1},
					[]interface{}{"b", 2},
				},
				time.Date(2019, 11, 4, 0, 0, 0, 0, time.UTC),
			},
			want: []interface{}{
				"hello",
				1234,
				"a", 1, "b", 2,
				time.Date(2019, 11, 4, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "complex",
			args: []interface{}{
				interface{}(1),
				[][]interface{}{
					[]interface{}{
						100,
						[]interface{}{
							"hello",
							[]byte("world"),
						},
						200,
					},
					[]interface{}{"b", 2},
				},
				[2]int{1, 2},
				time.Date(2019, 11, 4, 0, 0, 0, 0, time.UTC),
			},
			want: []interface{}{
				1,
				100,
				"hello",
				[]byte("world"),
				200,
				"b", 2,
				1, 2,
				time.Date(2019, 11, 4, 0, 0, 0, 0, time.UTC),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Flatten(tt.args)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("(-want, +got)\n%s\n%#v", diff, got)
			}
		})
	}
}
