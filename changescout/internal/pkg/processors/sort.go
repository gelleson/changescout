package processors

import (
	"bytes"
	"github.com/gelleson/changescout/changescout/internal/domain"
)

type SortProcessor struct {
	conf domain.Setting
}

func NewSortProcessor(conf domain.Setting) *SortProcessor {
	return &SortProcessor{
		conf: conf,
	}
}

func (p *SortProcessor) Skip() bool {
	return !p.conf.Sort
}

func (p *SortProcessor) Process(body []byte) []byte {
	// If sorting is disabled, return the original body
	if p.Skip() {
		return body
	}

	// If body is empty, return empty slice (not nil)
	if len(body) == 0 {
		return []byte{}
	}

	// Split the input string by newlines
	lines := bytes.Split(body, []byte{'\n'})

	// Create a slice for sorting
	var sortedLines [][]byte
	for _, line := range lines {
		// We only want non-empty lines
		if len(line) > 0 {
			sortedLines = append(sortedLines, line)
		}
	}

	// Sort the lines
	sortBytes(sortedLines)

	// Join the sorted lines back into a single byte slice with newline
	var output []byte
	for _, line := range sortedLines {
		output = append(output, line...)
		output = append(output, '\n') // Append newline after each line
	}

	return output
}

func sortBytes(lines [][]byte) {
	for i := 0; i < len(lines); i++ {
		for j := i + 1; j < len(lines); j++ {
			if bytes.Compare(lines[i], lines[j]) > 0 {
				lines[i], lines[j] = lines[j], lines[i]
			}
		}
	}
}
