package grpc

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/samverrall/microservice-example/internal/grpc/handler"
	"google.golang.org/grpc"
)

type ServerOption func(*Server)

func WithPort(port int) ServerOption {
	return func(s *Server) {
		s.port = port
	}
}

type Server struct {
	grpc        *grpc.Server
	port        int
	baseHandler *handler.Handler
}

func (s *Server) Server() *grpc.Server {
	return s.grpc
}

func (s *Server) Serve() error {
	addr := fmt.Sprintf(":%d", s.port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to create listener: %v", err)
	}
	slog.Info("grpc server listening", "addr", addr)
	return s.grpc.Serve(l)
}

func NewServer(h *handler.Handler, opts ...ServerOption) (*Server, error) {
	grpcServer := grpc.NewServer()

	server := &Server{
		port:        8000, // Default port
		grpc:        grpcServer,
		baseHandler: h,
	}

	// Apply server options
	for _, opt := range opts {
		opt(server)
	}

	return server, nil
}
