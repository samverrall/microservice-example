package handler

import (
	"context"
	"errors"

	"github.com/samverrall/microservice-example/internal/app"
	"github.com/samverrall/microservice-example/internal/app/user"
	"github.com/samverrall/microservice-example/internal/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	User *user.Service
}

type Repo struct {
	User user.Reader
}

type Handler struct {
	Svc  *Service
	Repo *Repo
}

func NewHandler(svc *Service, repo *Repo) *Handler {
	return &Handler{
		Svc:  svc,
		Repo: repo,
	}
}

func (h *Handler) Error(ctx context.Context, err error) error {
	switch {
	case err == nil:
		return nil

	case errors.Is(err, repository.ErrNotFound):
		return status.Error(codes.NotFound, "The requested resource was not found.")

	case errors.Is(err, app.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, err.Error())

	case errors.Is(err, app.ErrForbidden):
		return status.Error(codes.PermissionDenied, "You do not have permission to perform this action.")

	default:
		return status.Error(codes.Internal, "An internal error occurred.")
	}
}
