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
	logger      *zap.Logger
	router      *echo.Echo
	options     *option
	once        sync.Once
	middlewares []echo.MiddlewareFunc
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
		middlewares: []echo.MiddlewareFunc{
			middleware.Gzip(),
			middleware.CORSWithConfig(middleware.CORSConfig{
				AllowOrigins: []string{"*"},
				AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch},
				AllowHeaders: []string{"Content-Type", "Authorization", "X-Requested-With"},
			}),
			middleware.Recover(),
			middleware.RequestID(),
			echozap.ZapLogger(logger),
			WithCtx(),
		},
	}
}

func (s *Server) WithMiddlewares(middlewares ...echo.MiddlewareFunc) *Server {
	s.middlewares = append(s.middlewares, middlewares...)
	return s
}

func (s *Server) Register(method string, pattern string, handler echo.HandlerFunc, middlewares ...echo.MiddlewareFunc) {
	s.router.Add(method, pattern, handler, middlewares...)
}

func (s *Server) Start() error {
	s.router.Use(s.middlewares...)
	RegisterHandlers(s.router)

	s.logger.Info("Starting server")
	return http.ListenAndServe(":3311", s.router)
}
