package validation

import (
	"github.com/arthurpessoa/user-service/internal/usersvc/user"
	"errors"
	"time"
)


var ErrEmptyEmail = errors.New("empty email")

type Service interface {
	Validate(email string) error
}

type service struct {
	domains user.Repository
}

func NewService(domains user.Repository) Service {

	return &service{
		domains: domains,
	}
}

func (s *service) Validate(email string) error {

	if email == "" {
		return ErrEmptyEmail
	}

	return nil
}