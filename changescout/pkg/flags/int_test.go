package flags

import (
	"github.com/urfave/cli/v2"
	"testing"
	"time"
)

func TestIntFlag(t *testing.T) {
	tests := []struct {
		name     string
		opts     []OptionFunc[int]
		expected int
	}{
		{
			name:     "basic",
			opts:     nil,
			expected: 0,
		},
		{
			name:     "default",
			opts:     []OptionFunc[int]{WithDefaultValue(42)},
			expected: 42,
		},
		{
			name:     "envVar",
			opts:     []OptionFunc[int]{WithEnvVars[int]("INT_ENV_VAR")},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iFlag := NewIntFlag(tt.name, tt.opts...)
			flag := iFlag.Build().(*cli.IntFlag) // (*cli.IntFlag)

			app := &cli.App{
				Flags: []cli.Flag{flag},
				Action: func(ctx *cli.Context) error {
					val := iFlag.Get(ctx)
					if val != tt.expected {
						t.Errorf("expected %d, got %d", tt.expected, val)
					}
					return nil
				},
			}

			err := app.Run([]string{"app"})
			if err != nil {
				t.Fatalf("failed to run app: %v", err)
			}
		})
	}
}

func TestStringFlag(t *testing.T) {
	tests := []struct {
		name     string
		opts     []OptionFunc[string]
		expected string
	}{
		{
			name:     "basic",
			opts:     nil,
			expected: "",
		},
		{
			name:     "default",
			opts:     []OptionFunc[string]{WithDefaultValue("hello")},
			expected: "hello",
		},
		{
			name:     "envVar",
			opts:     []OptionFunc[string]{WithEnvVars[string]("STRING_ENV_VAR")},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sFlag := NewStringFlag(tt.name, tt.opts...)
			flag := sFlag.Build().(*cli.StringFlag) // (*cli.StringFlag)

			app := &cli.App{
				Flags: []cli.Flag{flag},
				Action: func(ctx *cli.Context) error {
					val := sFlag.Get(ctx)
					if val != tt.expected {
						t.Errorf("expected %s, got %s", tt.expected, val)
					}
					return nil
				},
			}

			err := app.Run([]string{"app"})
			if err != nil {
				t.Fatalf("failed to run app: %v", err)
			}
		})
	}
}

func TestBoolFlag(t *testing.T) {
	tests := []struct {
		name     string
		opts     []OptionFunc[bool]
		expected bool
	}{
		{
			name:     "default",
			opts:     []OptionFunc[bool]{WithDefaultValue(true)},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bFlag := NewBoolFlag(tt.name, tt.opts...)
			flag := bFlag.Build().(*cli.BoolFlag) // (*cli.BoolFlag)

			app := &cli.App{
				Flags: []cli.Flag{flag},
				Action: func(ctx *cli.Context) error {
					val := bFlag.Get(ctx)
					if val != tt.expected {
						t.Errorf("expected %v, got %v", tt.expected, val)
					}
					return nil
				},
			}

			err := app.Run([]string{"app"})
			if err != nil {
				t.Fatalf("failed to run app: %v", err)
			}
		})
	}
}

func TestDurationFlag(t *testing.T) {
	tests := []struct {
		name     string
		opts     []OptionFunc[time.Duration]
		expected time.Duration
	}{
		{
			name:     "basic",
			opts:     nil,
			expected: 0,
		},
		{
			name:     "default",
			opts:     []OptionFunc[time.Duration]{WithDefaultValue(10 * time.Second)},
			expected: 10 * time.Second,
		},
		{
			name:     "envVar",
			opts:     []OptionFunc[time.Duration]{WithEnvVars[time.Duration]("DURATION_ENV_VAR")},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dFlag := NewDurationFlag(tt.name, tt.opts...)
			flag := dFlag.Build().(*cli.DurationFlag) // (*cli.DurationFlag)

			app := &cli.App{
				Flags: []cli.Flag{flag},
				Action: func(ctx *cli.Context) error {
					val := dFlag.Get(ctx)
					if val != tt.expected {
						t.Errorf("expected %v, got %v", tt.expected, val)
					}
					return nil
				},
			}

			err := app.Run([]string{"app"})
			if err != nil {
				t.Fatalf("failed to run app: %v", err)
			}
		})
	}
}
