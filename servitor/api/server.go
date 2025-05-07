package api

import (
	"context"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
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
			if err := s.s.Serve(l); err != nil {
				log.WithError(err).Error("error serving http")
			}
		}(listener)
	}
}

func (s *Server) Shutdown(ctx context.Context) {
	if err := s.s.Shutdown(ctx); err != nil {
		log.WithError(err).Error("error closing http server")
	}
}
