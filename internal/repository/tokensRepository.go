package repository

import (
	"context"

	"github.com/jamascrorpJS/auth/internal/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TokensRepository struct {
	Db *mongo.Collection
}

func NewTokensRepository(db mongo.Database) *TokensRepository {
	return &TokensRepository{
		Db: db.Collection("Tokens"),
	}
}

func (t *TokensRepository) CreateTokens(ctx context.Context, tokens core.RefreshToken) (string, error) {
	res, err := t.Db.InsertOne(ctx, tokens)
	return res.InsertedID.(string), err
}

func (t *TokensRepository) GetTokens(ctx context.Context, id string) (core.RefreshToken, error) {
	token := core.RefreshToken{}
	err := t.Db.FindOne(ctx, bson.M{"_id": id}).Decode(&token)
	return token, err
}

func (t *TokensRepository) UpdateToken(ctx context.Context, id string) {
	t.Db.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"is_use": true}})
}
