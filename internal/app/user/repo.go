package user

import "context"

type Reader interface {
	FindAccountByID(ctx context.Context, id string) (*Account, error)
}

type Writer interface {
	AddAccount(ctx context.Context, a *Account) error
}

type ReadWriter interface {
	Reader
	Writer
}
