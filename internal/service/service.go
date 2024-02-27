package service

import (
	"context"

	"github.com/Microsoft/go-winio/pkg/guid"
	"github.com/jamascrorpJS/auth/internal/core"
	"github.com/jamascrorpJS/auth/internal/repository"
)

type ServiceUserI interface {
	CreateUser(ctx context.Context, user core.User) (string, error)
	ExistUser(ctx context.Context, id guid.GUID) (core.User, error)
}

type ServiceTokensI interface {
	CreateTokens(ctx context.Context, tokens core.RefreshToken) (string, error)
	GetTokens(ctx context.Context, id string) (core.RefreshToken, error)
	UpdateToken(ctx context.Context, id string)
	CreateTokenPairs(id guid.GUID) (core.TokenPair, error)
}
type ServiceRepository struct {
	ServiceUser   ServiceUserI
	ServiceTokens ServiceTokensI
}

func NewService(repository *repository.Repository) *ServiceRepository {

	return &ServiceRepository{
		ServiceUser:   NewUserService(repository),
		ServiceTokens: NewServiceTokens(repository),
	}
}
