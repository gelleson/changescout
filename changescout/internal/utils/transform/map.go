package transform

func MapObject[T any, R any](input T, fn func(T) R) R {
	return fn(input)
}

func MapObjects[T any, R any](input []T, fn func(T) R) []R {
	var result []R
	for _, item := range input {
		result = append(result, MapObject(item, fn))
	}
	return result
}

func ForEach[T any](input []T, fn func(T)) {
	for _, item := range input {
		fn(item)
	}
}

func ToPtr[T any](value T) *T {
	return &value
}

func Unwrap[T any](v *T) T {
	if v == nil {
		var zero T
		return zero
	}
	return *v
}

func ToPtrOrNil[T any](value T, checker func(T) bool) *T {
	if checker(value) {
		return &value
	}
	return nil
}

func ToValueOrDefault[T any](value *T, defaultValue T) T {
	if value != nil {
		return *value
	}
	return defaultValue
}

func IsZero[T comparable]() func(value T) bool {
	return func(value T) bool {
		var zero T
		return value == zero
	}
}

func Pipe[T any](input T, fns ...func(T) T) T {
	for _, fn := range fns {
		input = fn(input)
	}
	return input
}
