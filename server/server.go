package server

import (
	"github.com/Pomog/real-time-forum-V2/gorouter"
	"github.com/Pomog/real-time-forum-V2/internal/config"
	"log"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(conf *config.Conf, router *gorouter.Router) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + conf.API.Port,
			Handler: router,
		},
	}
}

func (s *Server) Run() error {
	log.Printf("API server is starting at %v", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}
