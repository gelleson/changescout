package processors_test

import (
	"testing"

	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/gelleson/changescout/changescout/internal/pkg/processors"
	"github.com/stretchr/testify/assert"
)

func TestDeduplicationProcessor_Process(t *testing.T) {
	tests := []struct {
		name     string
		conf     domain.Setting
		input    []byte
		expected []byte
	}{
		{
			name: "With deduplication",
			conf: domain.Setting{
				Deduplication: true,
			},
			input:    []byte("line1\nline2\nline1\nline3\nline2\n"),
			expected: []byte("line1\nline2\nline3\n"),
		},
		{
			name: "Without deduplication",
			conf: domain.Setting{
				Deduplication: false,
			},
			input:    []byte("line1\nline2\nline1\nline3\nline2\n"),
			expected: []byte("line1\nline2\nline1\nline3\nline2\n"),
		},
		{
			name: "Empty input",
			conf: domain.Setting{
				Deduplication: true,
			},
			input:    []byte{},
			expected: []byte{},
		},
		{
			name: "Input with no newlines",
			conf: domain.Setting{
				Deduplication: true,
			},
			input:    []byte("line1line2line1line3line2"),
			expected: []byte("line1\nline2\nline3\n"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := processors.NewDeduplicationProcessor(test.conf)
			actual := p.Process(test.input)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestDeduplicationProcessor_Skip(t *testing.T) {
	tests := []struct {
		name     string
		conf     domain.Setting
		expected bool
	}{
		{
			name: "With deduplication",
			conf: domain.Setting{
				Deduplication: true,
			},
			expected: false,
		},
		{
			name: "Without deduplication",
			conf: domain.Setting{
				Deduplication: false,
			},
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := processors.NewDeduplicationProcessor(test.conf)
			assert.Equal(t, test.expected, p.Skip())
		})
	}
}
