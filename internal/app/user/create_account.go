package user

import (
	"context"
	"fmt"

	"github.com/samverrall/microservice-example/internal/app"
)

func (s *Service) CreateAccount(ctx context.Context, email, password string) (*Account, error) {
	var input struct {
		email    Email
		password Password
	}
	{
		var err error

		if input.email, err = NewEmail(email); err != nil {
			return nil, app.NewInvalidInputErr(err)
		}

		if input.password, err = NewPassword(password); err != nil {
			return nil, app.NewInvalidInputErr(err)
		}
	}

	account := NewAccount(input.email, input.password)
	account.Create()

	if err := s.repo.AddAccount(ctx, account); err != nil {
		return nil, fmt.Errorf("add account: %w", err)
	}

	return account, nil
}
