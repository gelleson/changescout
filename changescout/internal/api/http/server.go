package http

import (
	"github.com/brpaz/echozap"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"net/http"
	"sync"
)

type Server struct {
	logger  *zap.Logger
	router  *echo.Echo
	options *option
	once    sync.Once
}

func runPipeline[T any](pipeline []func(*T), input *T) *T {
	for _, p := range pipeline {
		p(input)
	}
	return input
}

func New(logger *zap.Logger, options ...func(o *option)) *Server {
	return &Server{
		logger:  logger,
		router:  echo.New(),
		options: runPipeline[option](options, &option{}),
	}
}

func (s *Server) WithMiddlewares(middlewares ...echo.MiddlewareFunc) *Server {
	s.once.Do(func() {
		s.router.Use(
			middleware.Gzip(),
			middleware.Recover(),
			middleware.RequestID(),
			echozap.ZapLogger(s.logger),
			WithCtx(),
		)
	})

	s.router.Use(middlewares...)
	return s
}

func (s *Server) Register(method string, pattern string, handler echo.HandlerFunc, middlewares ...echo.MiddlewareFunc) {
	s.router.Add(method, pattern, handler, middlewares...)
}

func (s *Server) Start() error {
	s.logger.Info("Starting server")
	return http.ListenAndServe(":3311", s.router)
}
