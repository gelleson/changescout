package flags

import "github.com/urfave/cli/v2"

type Flag[T any] interface {
	Get(ctx *cli.Context) T
	Name() string
	Build() cli.Flag
	setDefault(v T)
	modifyDefault(fn func(*defaultFlag[T]))
}

type defaultFlag[T any] struct {
	required     bool
	description  string
	envVar       []string
	usage        string
	alias        string
	hidden       bool
	category     string
	defaultValue *T
}

type OptionFunc[T any] func(*defaultFlag[T])

func WithDescription[T any](description string) OptionFunc[T] {
	return func(f *defaultFlag[T]) {
		f.description = description
	}
}

func WithEnvVars[T any](envVars ...string) OptionFunc[T] {
	return func(f *defaultFlag[T]) {
		f.envVar = append(f.envVar, envVars...)
	}
}

func WithRequired[T any](required bool) OptionFunc[T] {
	return func(f *defaultFlag[T]) {
		f.required = required
	}
}

func WithUsage[T any](usage string) OptionFunc[T] {
	return func(f *defaultFlag[T]) {
		f.usage = usage
	}
}

func WithAlias[T any](alias string) OptionFunc[T] {
	return func(f *defaultFlag[T]) {
		f.alias = alias
	}
}
func Hidden[T any]() OptionFunc[T] {
	return func(f *defaultFlag[T]) {
		f.hidden = true
	}
}

func WithDefaultValue[T any](v T) OptionFunc[T] {
	return func(d *defaultFlag[T]) {
		d.defaultValue = &v
	}
}

type Builder interface {
	Build() cli.Flag
}

func Build(flags ...Builder) []cli.Flag {
	var result []cli.Flag
	for _, flag := range flags {
		result = append(result, flag.Build())
	}
	return result
}

func WithCategory[T any](category string) OptionFunc[T] {
	return func(f *defaultFlag[T]) {
		f.category = category
	}
}
