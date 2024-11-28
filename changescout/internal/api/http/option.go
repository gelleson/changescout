package http

type option struct {
	port string
}

type Option func(o *option)

func WithPort(port string) Option {
	return func(o *option) {
		o.port = port
	}
}
