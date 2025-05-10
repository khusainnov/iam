package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json2"
	"github.com/khusainnov/iam/app/config"
	"github.com/khusainnov/iam/app/handler"
)

const (
	HTTP_API_PREFIX = "/jsonrpc/v2"
)

// Server example
type Server struct {
	httpServer *http.Server
	cfg        *config.Server
}

func New(cfg *config.Server) *Server {
	return &Server{cfg: cfg}
}

func (s *Server) Init(handlers *handler.Handler) error {
	rpcServer := rpc.NewServer()
	rpcServer.RegisterCodec(json2.NewCodec(), "application/json")
	if err := rpcServer.RegisterService(handlers, ""); err != nil {
		return fmt.Errorf("failed to register service: %w", err)
	}

	router := mux.NewRouter()
	router.Handle(HTTP_API_PREFIX, rpcServer)

	s.httpServer = &http.Server{
		Addr:    s.cfg.Addr,
		Handler: router,
	}

	if err := http.ListenAndServe(s.cfg.Addr, router); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func (s *Server) Run() error {
	if err := s.httpServer.ListenAndServe(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if s.httpServer == nil {
		return fmt.Errorf("server is not running")
	}

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to gracefully shutdown the server: %w", err)
	}

	return nil
}

