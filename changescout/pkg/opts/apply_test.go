package opts

import (
	"testing"
)

func TestApply(t *testing.T) {
	type testStruct struct {
		value int
	}

	addOne := func(t *testStruct) {
		t.value += 1
	}

	multiplyByTwo := func(t *testStruct) {
		t.value *= 2
	}

	reset := func(t *testStruct) {
		t.value = 0
	}

	cases := []struct {
		name     string
		input    *testStruct
		options  []func(*testStruct)
		expected int
	}{
		{"NoOptions", &testStruct{value: 5}, []func(*testStruct){}, 5},
		{"AddOne", &testStruct{value: 5}, []func(*testStruct){addOne}, 6},
		{"MultiplyByTwo", &testStruct{value: 5}, []func(*testStruct){multiplyByTwo}, 10},
		{"AddAndMultiply", &testStruct{value: 5}, []func(*testStruct){addOne, multiplyByTwo}, 12},
		{"Reset", &testStruct{value: 5}, []func(*testStruct){reset}, 0},
		{"NoInput", &testStruct{}, []func(*testStruct){addOne, multiplyByTwo}, 2},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := Apply(c.input, c.options...)
			if result.value != c.expected {
				t.Errorf("expected %v, got %v", c.expected, result.value)
			}
		})
	}
}
