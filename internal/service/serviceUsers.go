package service

import (
	"context"

	"github.com/Microsoft/go-winio/pkg/guid"
	"github.com/jamascrorpJS/auth/internal/core"
	"github.com/jamascrorpJS/auth/internal/repository"
	"github.com/jamascrorpJS/auth/pkg/hash"
)

type ServiceUser struct {
	repository repository.RepositoryUsersI
}

func NewUserService(repository *repository.Repository) *ServiceUser {
	return &ServiceUser{repository: repository.RepositoryUsers}
}
func (s *ServiceUser) CreateUser(ctx context.Context, user core.User) (string, error) {
	passwordEncoded, _ := hash.Encode(user.Password)
	user.Password = passwordEncoded
	return s.repository.CreateUser(ctx, user)
}
func (s *ServiceUser) ExistUser(ctx context.Context, id guid.GUID) (core.User, error) {
	return s.repository.ExistUser(ctx, id)
}
