package flags

import (
	"github.com/urfave/cli/v2"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithDescription(t *testing.T) {
	tests := []struct {
		desc string
		d    *defaultFlag[string]
	}{
		{"Setting description", &defaultFlag[string]{description: ""}},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			f := WithDescription[string]("Test description")
			f(tt.d)
			assert.Equal(t, "Test description", tt.d.description)
		})
	}
}

func TestWithEnvVars(t *testing.T) {
	tests := []struct {
		desc     string
		original []string
		envVars  []string
		expected []string
	}{
		{"Adding env vars", []string{"VAR1"}, []string{"VAR2", "VAR3"}, []string{"VAR1", "VAR2", "VAR3"}},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			d := &defaultFlag[string]{envVar: tt.original}
			f := WithEnvVars[string](tt.envVars...)
			f(d)
			assert.Equal(t, tt.expected, d.envVar)
		})
	}
}

func TestWithRequired(t *testing.T) {
	tests := []struct {
		desc     string
		required bool
	}{
		{"Set required to true", true},
		{"Set required to false", false},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			d := &defaultFlag[string]{required: !tt.required}
			f := WithRequired[string](tt.required)
			f(d)
			assert.Equal(t, tt.required, d.required)
		})
	}
}

func TestWithUsage(t *testing.T) {
	tests := []struct {
		desc  string
		usage string
	}{
		{"Set usage", "Usage message"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			d := &defaultFlag[string]{usage: ""}
			f := WithUsage[string](tt.usage)
			f(d)
			assert.Equal(t, tt.usage, d.usage)
		})
	}
}

func TestWithAlias(t *testing.T) {
	tests := []struct {
		desc  string
		alias string
	}{
		{"Set alias", "test-alias"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			d := &defaultFlag[string]{alias: ""}
			f := WithAlias[string](tt.alias)
			f(d)
			assert.Equal(t, tt.alias, d.alias)
		})
	}
}

func TestHidden(t *testing.T) {
	tests := []struct {
		desc   string
		hidden bool
	}{
		{"Set hidden to true", true},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			d := &defaultFlag[string]{hidden: !tt.hidden}
			f := Hidden[string]()
			f(d)
			assert.Equal(t, tt.hidden, d.hidden)
		})
	}
}

func TestWithDefaultValue(t *testing.T) {
	tests := []struct {
		desc         string
		defaultValue string
	}{
		{"Set default value", "DefaultValue"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			d := &defaultFlag[string]{defaultValue: nil}
			f := WithDefaultValue[string](tt.defaultValue)
			f(d)
			assert.Equal(t, tt.defaultValue, *d.defaultValue)
		})
	}
}

func TestBuild(t *testing.T) {
	tests := []struct {
		desc      string
		flags     []Builder
		expLength int
	}{
		{"No flags", []Builder{}, 0},
		{"One flag", []Builder{mockBuilder{}}, 1},
		{"Two flags", []Builder{mockBuilder{}, mockBuilder{}}, 2},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			result := Build(tt.flags...)
			assert.Equal(t, tt.expLength, len(result))
		})
	}
}

func TestWithCategory(t *testing.T) {
	tests := []struct {
		desc     string
		category string
	}{
		{"Set category", "Category"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			d := &defaultFlag[string]{category: ""}
			f := WithCategory[string](tt.category)
			f(d)
			assert.Equal(t, tt.category, d.category)
		})
	}
}

type mockBuilder struct{}

func (b mockBuilder) Build() cli.Flag {
	return &cli.StringFlag{}
}
