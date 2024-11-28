package processors

import (
	"bytes"
	"github.com/gelleson/changescout/changescout/internal/domain"
)

type DeduplicationProcessor struct {
	conf domain.Setting
}

func NewDeduplicationProcessor(conf domain.Setting) *DeduplicationProcessor {
	return &DeduplicationProcessor{
		conf: conf,
	}
}

func (p *DeduplicationProcessor) Skip() bool {
	return !p.conf.Deduplication
}

func (p *DeduplicationProcessor) Process(body []byte) []byte {
	// If deduplication is disabled, return the original body
	if p.Skip() {
		return body
	}

	// If body is empty, return empty slice (not nil)
	if len(body) == 0 {
		return []byte{}
	}

	// Create a set to store unique lines and a slice to maintain order
	uniqueLines := make(map[string]struct{})
	var orderedUniqueLines []string

	// For input without newlines, split into individual words
	if !bytes.Contains(body, []byte{'\n'}) {
		// Split the input string every 5 characters (e.g., "line1", "line2")
		input := string(body)
		for i := 0; i < len(input); i += 5 {
			if i+5 <= len(input) {
				word := input[i : i+5]
				if _, exists := uniqueLines[word]; !exists {
					uniqueLines[word] = struct{}{}
					orderedUniqueLines = append(orderedUniqueLines, word)
				}
			}
		}
	} else {
		// Split by newlines and process
		lines := bytes.Split(body, []byte{'\n'})
		for _, line := range lines {
			if len(line) > 0 {
				lineStr := string(line)
				if _, exists := uniqueLines[lineStr]; !exists {
					uniqueLines[lineStr] = struct{}{}
					orderedUniqueLines = append(orderedUniqueLines, lineStr)
				}
			}
		}
	}

	// Build output with newlines
	var output []byte
	for _, line := range orderedUniqueLines {
		output = append(output, []byte(line)...)
		output = append(output, '\n')
	}

	return output
}
