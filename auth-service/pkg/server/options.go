package server

import (
	"time"
)

type Option func(s *Server)

func WithReadTimeout(n time.Duration) Option {
	return func(s *Server) {
		s.httpSrv.ReadTimeout = n
	}
}

func WithWriteTimeout(n time.Duration) Option {
	return func(s *Server) {
		s.httpSrv.WriteTimeout = n
	}
}

func WithIdleTimeout(n time.Duration) Option {
	return func(s *Server) {
		s.httpSrv.IdleTimeout = n
	}
}

func WithShutdownTimeout(n time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = n
	}
}
