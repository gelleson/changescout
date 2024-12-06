package transform

import (
	"testing"
)

func TestMapObject(t *testing.T) {
	type testCase struct {
		name   string
		input  int
		fn     func(int) int
		expect int
	}

	tests := []testCase{
		{"identity", 1, func(x int) int { return x }, 1},
		{"square", 2, func(x int) int { return x * x }, 4},
		{"double", 3, func(x int) int { return x * 2 }, 6},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := MapObject(tc.input, tc.fn)
			if result != tc.expect {
				t.Errorf("expected %v, got %v", tc.expect, result)
			}
		})
	}
}

func TestMapObjects(t *testing.T) {
	type testCase struct {
		name   string
		input  []int
		fn     func(int) int
		expect []int
	}

	tests := []testCase{
		{"empty", []int{}, func(x int) int { return x }, []int{}},
		{"identity", []int{1, 2}, func(x int) int { return x }, []int{1, 2}},
		{"square", []int{1, 2, 3}, func(x int) int { return x * x }, []int{1, 4, 9}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := MapObjects(tc.input, tc.fn)
			for i, v := range result {
				if v != tc.expect[i] {
					t.Errorf("expected %v, got %v", tc.expect, result)
					break
				}
			}
		})
	}
}

func TestForEach(t *testing.T) {
	cases := []struct {
		name  string
		input []int
		fn    func(int) int
		sum   int
	}{
		{"empty", []int{}, func(x int) int { return x }, 0},
		{"increment", []int{1, 2, 3}, func(x int) int { return x + 1 }, 9},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sum := 0
			ForEach(tc.input, func(n int) {
				sum += tc.fn(n)
			})
			if sum != tc.sum {
				t.Errorf("expected %v, got %v", tc.sum, sum)
			}
		})
	}
}

func TestToPtr(t *testing.T) {
	cases := []struct {
		name   string
		input  int
		expect *int
	}{
		{"zero", 0, func() *int { v := 0; return &v }()},
		{"positive", 5, func() *int { v := 5; return &v }()},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := ToPtr(tc.input)
			if *result != *tc.expect {
				t.Errorf("expected %v, got %v", *tc.expect, *result)
			}
		})
	}
}

func TestUnwrap(t *testing.T) {
	cases := []struct {
		name   string
		input  *int
		expect int
	}{
		{"nil", nil, 0},
		{"non-nil", func() *int { v := 5; return &v }(), 5},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := Unwrap(tc.input)
			if result != tc.expect {
				t.Errorf("expected %v, got %v", tc.expect, result)
			}
		})
	}
}

func TestToPtrOrNil(t *testing.T) {
	cases := []struct {
		name    string
		input   int
		checker func(int) bool
		expect  *int
	}{
		{"checker true", 5, func(x int) bool { return x > 0 }, func() *int { v := 5; return &v }()},
		{"checker false", 0, func(x int) bool { return x > 0 }, nil},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := ToPtrOrNil(tc.input, tc.checker)
			if (result == nil) != (tc.expect == nil) || (result != nil && *result != *tc.expect) {
				t.Errorf("expected %v, got %v", tc.expect, result)
			}
		})
	}
}

func TestToValueOrDefault(t *testing.T) {
	defaultVal := 10
	cases := []struct {
		name         string
		input        *int
		defaultValue int
		expect       int
	}{
		{"nil input", nil, defaultVal, 10},
		{"non-nil input", func() *int { v := 5; return &v }(), defaultVal, 5},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := ToValueOrDefault(tc.input, tc.defaultValue)
			if result != tc.expect {
				t.Errorf("expected %v, got %v", tc.expect, result)
			}
		})
	}
}

func TestIsZero(t *testing.T) {
	isZero := IsZero[int]()
	cases := []struct {
		name   string
		value  int
		expect bool
	}{
		{"zero", 0, true},
		{"non-zero", 1, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := isZero(tc.value)
			if result != tc.expect {
				t.Errorf("expected %v, got %v", tc.expect, result)
			}
		})
	}
}

func TestPipe(t *testing.T) {
	cases := []struct {
		name   string
		input  int
		fns    []func(int) int
		expect int
	}{
		{"identity", 1, []func(int) int{func(x int) int { return x }}, 1},
		{"increment", 1, []func(int) int{func(x int) int { return x + 1 }}, 2},
		{"double then add", 1, []func(int) int{func(x int) int { return x * 2 }, func(x int) int { return x + 3 }}, 5},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := Pipe(tc.input, tc.fns...)
			if result != tc.expect {
				t.Errorf("expected %v, got %v", tc.expect, result)
			}
		})
	}
}
