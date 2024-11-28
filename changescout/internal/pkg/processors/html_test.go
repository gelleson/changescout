package processors_test

import (
	"testing"

	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/gelleson/changescout/changescout/internal/pkg/processors"
	"github.com/stretchr/testify/assert"
)

func TestHTMLProcessor_Process(t *testing.T) {
	tests := []struct {
		name     string
		conf     domain.Setting
		input    []byte
		expected []byte
	}{
		{
			name: "With selectors",
			conf: domain.Setting{
				Selectors: []string{"h1", "p"},
			},
			input: []byte(`
				<html>
					<body>
						<h1>Heading 1</h1>
						<p>Paragraph 1</p>
						<p>Paragraph 2</p>
					</body>
				</html>
			`),
			expected: []byte("Heading 1Paragraph 1Paragraph 2"),
		},
		{
			name: "Without selectors",
			conf: domain.Setting{
				Selectors: nil,
			},
			input: []byte(`
				<html>
					<body>
						<h1>Heading 1</h1>
						<p>Paragraph 1</p>
						<p>Paragraph 2</p>
					</body>
				</html>
			`),
			expected: []byte(`
				<html>
					<body>
						<h1>Heading 1</h1>
						<p>Paragraph 1</p>
						<p>Paragraph 2</p>
					</body>
				</html>
			`),
		},
		{
			name:     "Invalid HTML",
			conf:     domain.Setting{},
			input:    []byte("This is not valid HTML"),
			expected: []byte("This is not valid HTML"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := processors.NewHTMLProcessor(test.conf)
			actual := p.Process(test.input)
			assert.Equal(t, test.expected, actual)
		})
	}
}
func TestHTMLProcessor_Skip(t *testing.T) {
	tests := []struct {
		name     string
		conf     domain.Setting
		expected bool
	}{
		{
			name: "With selectors",
			conf: domain.Setting{
				Selectors: []string{"p", "h1"},
			},
			expected: false,
		},
		{
			name: "Without selectors",
			conf: domain.Setting{
				Selectors: nil,
			},
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := processors.NewHTMLProcessor(test.conf)
			assert.Equal(t, test.expected, p.Skip())
		})
	}
}
