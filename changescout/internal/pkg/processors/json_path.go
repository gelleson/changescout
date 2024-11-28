package processors

import (
	"encoding/json"
	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/oliveagle/jsonpath"
)

type JSONPathProcessor struct {
	conf domain.Setting
}

func NewJSONPathProcessor(conf domain.Setting) *JSONPathProcessor {
	return &JSONPathProcessor{
		conf: conf,
	}
}

func (p *JSONPathProcessor) Skip() bool {
	return len(p.conf.JSONPath) == 0
}

func (p *JSONPathProcessor) Process(body []byte) []byte {
	// If JSONPath is empty, return the original body
	if p.Skip() {
		return body
	}

	// If the body is empty, return an empty slice
	if len(body) == 0 {
		return []byte{}
	}

	// Parse the JSON input first
	var jsonData interface{}
	if err := json.Unmarshal(body, &jsonData); err != nil {
		return body // Return original body if JSON is invalid
	}

	// Use a map to store the evaluated JSONPath results
	result := make(map[string]interface{})
	foundValidPath := false

	// Evaluate the JSONPath expressions
	for _, path := range p.conf.JSONPath {
		// Perform the JSONPath lookup on the parsed JSON
		value, err := jsonpath.JsonPathLookup(jsonData, path)
		if err == nil {
			result[path] = value
			foundValidPath = true
		}
	}

	// If we found valid paths, return the result
	if foundValidPath {
		output, err := json.Marshal(result)
		if err != nil {
			return []byte(`{}`)
		}
		return output
	}

	// If no valid paths were found and we have JSONPaths configured
	if len(p.conf.JSONPath) > 0 {
		return []byte(`{}`)
	}

	// If no JSONPaths were configured, return the original JSON
	return body
}
