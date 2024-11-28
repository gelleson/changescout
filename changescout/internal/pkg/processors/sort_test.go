package processors_test

import (
	"testing"

	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/gelleson/changescout/changescout/internal/pkg/processors"
	"github.com/stretchr/testify/assert"
)

func TestSortProcessor_Process(t *testing.T) {
	tests := []struct {
		name     string
		conf     domain.Setting
		input    []byte
		expected []byte
	}{
		{
			name: "With sorting",
			conf: domain.Setting{
				Sort: true,
			},
			input:    []byte("banana\napple\ncherry\n"),
			expected: []byte("apple\nbanana\ncherry\n"),
		},
		{
			name: "Without sorting",
			conf: domain.Setting{
				Sort: false,
			},
			input:    []byte("banana\napple\ncherry\n"),
			expected: []byte("banana\napple\ncherry\n"),
		},
		{
			name: "Empty input",
			conf: domain.Setting{
				Sort: true,
			},
			input:    []byte{},
			expected: []byte{},
		},
		{
			name: "Single line input",
			conf: domain.Setting{
				Sort: true,
			},
			input:    []byte("single\n"),
			expected: []byte("single\n"),
		},
		{
			name: "Already sorted input",
			conf: domain.Setting{
				Sort: true,
			},
			input:    []byte("apple\nbanana\ncherry\n"),
			expected: []byte("apple\nbanana\ncherry\n"),
		},
		{
			name: "Duplicated lines",
			conf: domain.Setting{
				Sort: true,
			},
			input:    []byte("banana\napple\nbanana\ncherry\n"),
			expected: []byte("apple\nbanana\nbanana\ncherry\n"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := processors.NewSortProcessor(test.conf)
			actual := p.Process(test.input)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestSortProcessor_Skip(t *testing.T) {
	tests := []struct {
		name     string
		conf     domain.Setting
		expected bool
	}{
		{
			name: "With sorting",
			conf: domain.Setting{
				Sort: true,
			},
			expected: false,
		},
		{
			name: "Without sorting",
			conf: domain.Setting{
				Sort: false,
			},
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := processors.NewSortProcessor(test.conf)
			assert.Equal(t, test.expected, p.Skip())
		})
	}
}
