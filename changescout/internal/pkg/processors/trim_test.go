package processors_test

import (
	"testing"

	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/gelleson/changescout/changescout/internal/pkg/processors"
	"github.com/stretchr/testify/assert"
)

func TestTrimProcessor_Process(t *testing.T) {
	tests := []struct {
		name     string
		conf     domain.Setting
		input    []byte
		expected []byte
	}{
		{
			name: "With trimming",
			conf: domain.Setting{
				Trim: true,
			},
			input:    []byte("   Hello, World!   \n"),
			expected: []byte("Hello, World!"),
		},
		{
			name: "With trimming on empty input",
			conf: domain.Setting{
				Trim: true,
			},
			input:    []byte("   \n"),
			expected: []byte(""),
		},
		{
			name: "Without trimming",
			conf: domain.Setting{
				Trim: false,
			},
			input:    []byte("   Hello, World!   \n"),
			expected: []byte("   Hello, World!   \n"),
		},
		{
			name: "With trimming with newline only",
			conf: domain.Setting{
				Trim: true,
			},
			input:    []byte("\n"),
			expected: []byte(""),
		},
		{
			name: "With trimming no spaces",
			conf: domain.Setting{
				Trim: true,
			},
			input:    []byte("NoSpaces"),
			expected: []byte("NoSpaces"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := processors.NewTrimProcessor(test.conf)
			actual := p.Process(test.input)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestTrimProcessor_Skip(t *testing.T) {
	tests := []struct {
		name     string
		conf     domain.Setting
		expected bool
	}{
		{
			name: "With trimming",
			conf: domain.Setting{
				Trim: true,
			},
			expected: false,
		},
		{
			name: "Without trimming",
			conf: domain.Setting{
				Trim: false,
			},
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := processors.NewTrimProcessor(test.conf)
			assert.Equal(t, test.expected, p.Skip())
		})
	}
}
