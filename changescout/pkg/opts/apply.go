package opts

type OptionFunc[T any] func(*T)

func Apply[T any](flag *T, opts ...func(*T)) T {
	for _, opt := range opts {
		opt(flag)
	}

	return *flag
}
