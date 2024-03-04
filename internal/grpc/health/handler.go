package health

import (
	"context"

	"github.com/samverrall/microservice-example/pkg/proto"
)

type Handler struct {
	proto.UnimplementedHealthServer
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Check(ctx context.Context, req *proto.HealthCheckRequest) (*proto.HealthCheckResponse, error) {
	return &proto.HealthCheckResponse{
		Status: proto.HealthCheckResponse_SERVING,
	}, nil
}

func (s *Handler) Watch(req *proto.HealthCheckRequest, stream proto.Health_WatchServer) error {
	return nil
}
