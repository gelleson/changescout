package browser

type options struct {
	managedInstanceURL *string
}

type optionFunc func(*options) *options

func WithManagedInstanceURL(url *string) optionFunc {
	return func(o *options) *options {
		if url == nil {
			return o
		}
		o.managedInstanceURL = url

		return o
	}
}

func Noop() optionFunc {
	return func(o *options) *options {
		return o
	}
}
