package health

import "github.com/samverrall/user-service/internal/grpc/handler"

type Handler struct {
	*handler.Handler
}

func NewHandler(h *handler.Handler) *Handler {
	return &Handler{
		Handler: h,
	}
}
