package service

import (
	"context"
	"time"

	"github.com/Microsoft/go-winio/pkg/guid"
	"github.com/jamascrorpJS/auth/internal/core"
	"github.com/jamascrorpJS/auth/internal/repository"
	"github.com/jamascrorpJS/auth/pkg/hash"
	"github.com/jamascrorpJS/auth/pkg/token"
)

type ServiceTokens struct {
	repository repository.RepositoryTokenI
}

func NewServiceTokens(repository *repository.Repository) *ServiceTokens {
	return &ServiceTokens{repository: repository.RepositoryToken}
}

func (s *ServiceTokens) CreateTokens(ctx context.Context, token core.RefreshToken) (string, error) {
	encodedToken, err := hash.Encode(token.Token)
	if err != nil {
		return "", err
	}
	token.Token = encodedToken
	return s.repository.CreateTokens(ctx, token)
}

func (s *ServiceTokens) GetTokens(ctx context.Context, id string) (core.RefreshToken, error) {
	return s.repository.GetTokens(ctx, id)
}

func (s *ServiceTokens) CreateTokenPairs(id guid.GUID) (core.TokenPair, error) {
	tokenID, err := guid.NewV4()
	if err != nil {
		return core.TokenPair{}, err
	}

	accessToken, err := token.CreateAccessToken(id, tokenID)
	if err != nil {
		return core.TokenPair{}, err
	}
	ttl := time.Now().Local().Add(token.TokenTTLs).UTC()
	refreshToken := core.RefreshToken{
		ID:        tokenID.String(),
		Token:     token.CreateRefershToken(),
		ExpiresAt: ttl,
		IsUse:     false,
	}
	access := core.TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}
	return access, nil
}

func (s *ServiceTokens) UpdateToken(ctx context.Context, id string) {
	s.repository.UpdateToken(ctx, id)
}
