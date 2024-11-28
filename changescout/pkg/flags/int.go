package flags

import (
	"github.com/urfave/cli/v2"
	"time"
)

// FlagBuilder is a generic interface for building different types of flags
type FlagBuilder[T any] struct {
	name        string
	defaultFlag defaultFlag[T]
	buildFn     func(name string, def defaultFlag[T]) cli.Flag
}

// NewFlag creates a new flag builder for the given type
func NewFlag[T any](name string, buildFn func(name string, def defaultFlag[T]) cli.Flag, opts ...OptionFunc[T]) *FlagBuilder[T] {
	f := &FlagBuilder[T]{
		name:    name,
		buildFn: buildFn,
	}

	for _, opt := range opts {
		opt(&f.defaultFlag)
	}

	return f
}

// Build creates the concrete cli.Flag
func (f *FlagBuilder[T]) Build() cli.Flag {
	return f.buildFn(f.name, f.defaultFlag)
}

func (f *FlagBuilder[T]) Name() string {
	return f.name
}

// Flag types
type intFlag struct {
	name string
	opts []OptionFunc[int]
	flag *FlagBuilder[int]
}

type stringFlag struct {
	name string
	opts []OptionFunc[string]
	flag *FlagBuilder[string]
}

type boolFlag struct {
	name string
	opts []OptionFunc[bool]
	flag *FlagBuilder[bool]
}

type durationFlag struct {
	name string
	opts []OptionFunc[time.Duration]
	flag *FlagBuilder[time.Duration]
}

// Int Flag
func NewIntFlag(name string, opts ...OptionFunc[int]) *intFlag {
	return &intFlag{
		name: name,
		opts: opts,
	}
}

func (i *intFlag) Get(ctx *cli.Context) int {
	return ctx.Int(i.name)
}

func (i *intFlag) Build() cli.Flag {
	i.flag = NewFlag[int](i.name, func(name string, def defaultFlag[int]) cli.Flag {
		return &cli.IntFlag{
			Name:        name,
			Required:    def.required,
			DefaultText: def.description,
			Usage:       def.usage,
			EnvVars:     def.envVar,
			Aliases:     []string{def.alias},
			Hidden:      def.hidden,
			Value:       defaultValue(def.defaultValue),
			Category:    def.category,
		}
	}, i.opts...)
	return i.flag.Build()
}

// String Flag
func NewStringFlag(name string, opts ...OptionFunc[string]) *stringFlag {
	return &stringFlag{
		name: name,
		opts: opts,
	}
}

func (s *stringFlag) Get(ctx *cli.Context) string {
	return ctx.String(s.name)
}

func (s *stringFlag) Build() cli.Flag {
	s.flag = NewFlag[string](s.name, func(name string, def defaultFlag[string]) cli.Flag {
		return &cli.StringFlag{
			Name:        name,
			Required:    def.required,
			DefaultText: def.description,
			Usage:       def.usage,
			EnvVars:     def.envVar,
			Aliases:     []string{def.alias},
			Hidden:      def.hidden,
			Value:       defaultValue(def.defaultValue),
			Category:    def.category,
		}
	}, s.opts...)
	return s.flag.Build()
}

// Bool Flag
func NewBoolFlag(name string, opts ...OptionFunc[bool]) *boolFlag {
	return &boolFlag{
		name: name,
		opts: opts,
	}
}

func (b *boolFlag) Get(ctx *cli.Context) bool {
	return ctx.Bool(b.name)
}

func (b *boolFlag) Build() cli.Flag {
	b.flag = NewFlag[bool](b.name, func(name string, def defaultFlag[bool]) cli.Flag {
		var alias []string
		if def.alias != "" {
			alias = []string{
				def.alias,
			}
		}
		return &cli.BoolFlag{
			Name:        name,
			Required:    def.required,
			DefaultText: def.description,
			Usage:       def.usage,
			EnvVars:     def.envVar,
			Aliases:     alias,
			Hidden:      def.hidden,
			Value:       true,
			Category:    def.category,
		}
	}, b.opts...)
	return b.flag.Build()
}

// Duration Flag
func NewDurationFlag(name string, opts ...OptionFunc[time.Duration]) *durationFlag {
	return &durationFlag{
		name: name,
		opts: opts,
	}
}

func (d *durationFlag) Get(ctx *cli.Context) time.Duration {
	return ctx.Duration(d.name)
}

func (d *durationFlag) Build() cli.Flag {
	d.flag = NewFlag[time.Duration](d.name, func(name string, def defaultFlag[time.Duration]) cli.Flag {
		return &cli.DurationFlag{
			Name:        name,
			Required:    def.required,
			DefaultText: def.description,
			Usage:       def.usage,
			EnvVars:     def.envVar,
			Aliases:     []string{def.alias},
			Hidden:      def.hidden,
			Value:       defaultValue(def.defaultValue),
			Category:    def.category,
		}
	}, d.opts...)
	return d.flag.Build()
}

// Helper function to handle nil default values
func defaultValue[T any](v *T) T {
	if v == nil {
		var zero T
		return zero
	}
	return *v
}
