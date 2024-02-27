package repository

import (
	"context"

	"github.com/Microsoft/go-winio/pkg/guid"
	"github.com/jamascrorpJS/auth/internal/core"
	"go.mongodb.org/mongo-driver/mongo"
)

type RepositoryTokenI interface {
	CreateTokens(ctx context.Context, tokens core.RefreshToken) (string, error)
	GetTokens(ctx context.Context, id string) (core.RefreshToken, error)
	UpdateToken(ctx context.Context, id string)
}

type RepositoryUsersI interface {
	CreateUser(ctx context.Context, user core.User) (string, error)
	ExistUser(ctx context.Context, user guid.GUID) (core.User, error)
}

type Repository struct {
	RepositoryToken RepositoryTokenI
	RepositoryUsers RepositoryUsersI
}

func NewRepository(db mongo.Database) *Repository {
	return &Repository{
		RepositoryToken: NewTokensRepository(db),
		RepositoryUsers: NewRepositoryUsers(db),
	}
}
