// build_test.go
package ent

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"

	_ "github.com/glebarez/go-sqlite"
)

func TestBuild(t *testing.T) {
	tests := []struct {
		name        string
		config      *BuildConfig
		expectError bool
	}{
		{name: "ValidConfig", config: &BuildConfig{DBEngine: "sqlite", DBURL: ":memory:?_pragma=foreign_keys(1)"}, expectError: false},
		{name: "InvalidEngine", config: &BuildConfig{DBEngine: "invalid", DBURL: ":memory:"}, expectError: true},
		{name: "InvalidURL", config: &BuildConfig{DBEngine: "sqlite", DBURL: "invalid_url"}, expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := Build(context.Background(), tt.config)
			if tt.expectError && err != nil {
				assert.Nil(t, client, "Client should be nil on error")
			} else {
				assert.NotNil(t, client, "Client should not be nil")
			}
		})
	}
}
