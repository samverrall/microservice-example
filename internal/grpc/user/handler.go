package userhandler

import (
	"context"
	"log/slog"

	"github.com/samverrall/microservice-example/internal/grpc/handler"
	"github.com/samverrall/microservice-example/pkg/proto"
)

type Handler struct {
	proto.UnimplementedUserServer

	*handler.Handler
}

func NewHandler(handler *handler.Handler) *Handler {
	return &Handler{
		Handler: handler,
	}
}

func (h *Handler) Signup(ctx context.Context, p *proto.SignUpRequest) (*proto.SignupResponse, error) {
	user, err := h.Svc.User.CreateAccount(ctx, p.Email, p.Password)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create account", "error", err)

		return nil, h.Error(ctx, err)
	}

	return &proto.SignupResponse{
		UserId: user.ID,
	}, nil
}
