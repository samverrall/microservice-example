package handler

import (
	"context"

	"github.com/samverrall/user-service/internal/app/user"
)

type Service struct {
	User *user.Service
}

// Repo holds any repositories that the gRPC handlers may utilise.
// Ideally only Reader repositories from the core application should be used here.
// This is because simple "get" gRPC methods do not require business logic and therefore
// can retrieve entities directly from the repository. Any write operations however
// should go through the core application layer.
type Repo struct {
	User user.Reader
}

// Handler acts as a base for any fields a gRPC handler may require.
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
	return nil
}
