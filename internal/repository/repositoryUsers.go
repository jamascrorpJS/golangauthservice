package repository

import (
	"context"

	"github.com/Microsoft/go-winio/pkg/guid"
	"github.com/jamascrorpJS/auth/internal/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RepositoryUsers struct {
	Db *mongo.Collection
}

func NewRepositoryUsers(db mongo.Database) *RepositoryUsers {
	return &RepositoryUsers{Db: db.Collection("Users")}
}

func (r *RepositoryUsers) CreateUser(ctx context.Context, user core.User) (string, error) {
	res, err := r.Db.InsertOne(ctx, user)
	return res.InsertedID.(string), err
}

func (r *RepositoryUsers) ExistUser(ctx context.Context, user guid.GUID) (core.User, error) {

	u := core.User{}
	err := r.Db.FindOne(ctx, bson.M{"_id": user}).Decode(&u)
	return u, err
}
