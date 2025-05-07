package api

import (
	"context"
	"net"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	s *http.Server
}

func NewServer(handler http.Handler) *Server {
	server := &Server{
		s: &http.Server{
			Handler: handler,
		},
	}
	return server
}

func (s *Server) Start(listeners ...net.Listener) {
	for _, listener := range listeners {
		go func(l net.Listener) {
			if err := s.s.Serve(l); err != nil && err != http.ErrServerClosed {
				log.WithError(err).Error("error serving http")
			}
		}(listener)
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.s.Shutdown(ctx); err != nil {
		log.WithError(err).Error("error closing http server")
		return err
	}
	return nil
}
