package day2

import (
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestMap(t *testing.T) {
	t.Run("int to int", func(t *testing.T) {
		tests := []struct {
			name string
			in   []int
			want []int
		}{
			{"double values", []int{1, 2, 3}, []int{2, 4, 6}},
			{"empty slice", []int{}, []int{}},
			{"nil slice", nil, []int{}},
			{"single element", []int{7}, []int{14}},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				got := Map(tc.in, func(n int) int { return n * 2 })
				if !reflect.DeepEqual(got, tc.want) {
					t.Errorf("got %v, want %v", got, tc.want)
				}
			})
		}
	})

	t.Run("int to string", func(t *testing.T) {
		got := Map([]int{1, 2, 3}, strconv.Itoa)
		want := []string{"1", "2", "3"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("string to int (length)", func(t *testing.T) {
		got := Map([]string{"go", "gin", "mongo"}, func(s string) int { return len(s) })
		want := []int{2, 3, 5}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestFilter(t *testing.T) {
	isEven := func(n int) bool { return n%2 == 0 }

	tests := []struct {
		name string
		in   []int
		keep func(int) bool
		want []int
	}{
		{"keep evens", []int{1, 2, 3, 4, 5, 6}, isEven, []int{2, 4, 6}},
		{"none match", []int{1, 3, 5}, isEven, []int{}},
		{"all match", []int{2, 4}, isEven, []int{2, 4}},
		{"empty slice", []int{}, isEven, []int{}},
		{"nil slice", nil, isEven, []int{}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Filter(tc.in, tc.keep)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}

	t.Run("strings with prefix", func(t *testing.T) {
		got := Filter([]string{"gin", "go", "mongo", "gorm"}, func(s string) bool {
			return strings.HasPrefix(s, "g")
		})
		want := []string{"gin", "go", "gorm"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
