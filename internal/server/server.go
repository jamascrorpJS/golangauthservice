package server

import (
	"net/http"
)

type Server struct {
	Server http.Server
}

func NewServer(addr string, handler http.Handler) *Server {
	return &Server{
		Server: http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
}

func (s *Server) StartServer() error {
	return s.Server.ListenAndServe()
}
