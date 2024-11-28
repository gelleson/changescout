package processors_test

import (
	"encoding/json"
	"testing"

	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/gelleson/changescout/changescout/internal/pkg/processors"
	"github.com/stretchr/testify/assert"
)

func TestJSONPathProcessor_Process(t *testing.T) {
	tests := []struct {
		name     string
		conf     domain.Setting
		input    string
		expected string
	}{
		{
			name: "With valid JSONPath",
			conf: domain.Setting{
				JSONPath: []string{"$.store.book[*].author"},
			},
			input: `{
				"store": {
					"book": [
						{"author": "John", "title": "Book 1"},
						{"author": "Jane", "title": "Book 2"}
					]
				}
			}`,
			expected: `{"$.store.book[*].author":["John","Jane"]}`,
		},
		{
			name: "With multiple JSONPaths",
			conf: domain.Setting{
				JSONPath: []string{
					"$.store.book[*].title",
					"$.store.book[*].author",
				},
			},
			input: `{
				"store": {
					"book": [
						{"author": "John", "title": "Book 1"},
						{"author": "Jane", "title": "Book 2"}
					]
				}
			}`,
			expected: `{
				"$.store.book[*].author": ["John","Jane"],
				"$.store.book[*].title": ["Book 1","Book 2"]
			}`,
		},
		{
			name: "With empty JSONPath",
			conf: domain.Setting{
				JSONPath: []string{},
			},
			input: `{
				"store": {
					"book": [
						{"author": "John", "title": "Book 1"}
					]
				}
			}`,
			expected: `{
				"store": {
					"book": [
						{"author": "John", "title": "Book 1"}
					]
				}
			}`,
		},
		{
			name: "With invalid JSONPath",
			conf: domain.Setting{
				JSONPath: []string{"$.invalid.path"},
			},
			input: `{
				"store": {
					"book": [
						{"author": "John", "title": "Book 1"}
					]
				}
			}`,
			expected: `{}`,
		},
		{
			name: "With empty input",
			conf: domain.Setting{
				JSONPath: []string{"$.store.book[*].author"},
			},
			input:    "",
			expected: "",
		},
		{
			name: "HTTP Bin",
			conf: domain.Setting{
				JSONPath: []string{"$.headers.User-Agent"},
			},
			input: `{
  "args": {}, 
  "headers": {
    "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7", 
    "Accept-Encoding": "gzip, deflate, br, zstd", 
    "Accept-Language": "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7", 
    "Dnt": "1", 
    "Host": "httpbin.org", 
    "Priority": "u=0, i", 
    "Sec-Ch-Ua": "\"Chromium\";v=\"131\", \"Not_A Brand\";v=\"24\"", 
    "Sec-Ch-Ua-Mobile": "?0", 
    "Sec-Ch-Ua-Platform": "\"macOS\"", 
    "Sec-Fetch-Dest": "document", 
    "Sec-Fetch-Mode": "navigate", 
    "Sec-Fetch-Site": "none", 
    "Sec-Fetch-User": "?1", 
    "Upgrade-Insecure-Requests": "1", 
    "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36", 
    "X-Amzn-Trace-Id": "Root=1-67484ae2-44c803aa34190bc030dc8d79"
  }, 
  "origin": "176.110.126.5", 
  "url": "https://httpbin.org/get"
}`,
			expected: `{"$.headers.User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := processors.NewJSONPathProcessor(test.conf)

			// Convert input string to bytes
			var input []byte
			if test.input != "" {
				// Normalize the JSON formatting
				var tmp interface{}
				err := json.Unmarshal([]byte(test.input), &tmp)
				assert.NoError(t, err)
				input, err = json.Marshal(tmp)
				assert.NoError(t, err)
			}

			// Get actual result
			actual := p.Process(input)

			// For non-empty expected results, normalize the JSON formatting
			var expectedBytes []byte
			if test.expected != "" {
				var tmp interface{}
				err := json.Unmarshal([]byte(test.expected), &tmp)
				assert.NoError(t, err)
				expectedBytes, err = json.Marshal(tmp)
				assert.NoError(t, err)
			}

			// Compare normalized JSON strings
			assert.Equal(t, string(expectedBytes), string(actual))
		})
	}
}

func TestJSONPathProcessor_Skip(t *testing.T) {
	tests := []struct {
		name     string
		conf     domain.Setting
		expected bool
	}{
		{
			name: "With JSONPath",
			conf: domain.Setting{
				JSONPath: []string{"$.store.book[*].author"},
			},
			expected: false,
		},
		{
			name: "Without JSONPath",
			conf: domain.Setting{
				JSONPath: []string{},
			},
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := processors.NewJSONPathProcessor(test.conf)
			assert.Equal(t, test.expected, p.Skip())
		})
	}
}
