package core

import (
	"time"
)

type TokenPair struct {
	AccessToken  string
	RefreshToken RefreshToken
}

type RefreshToken struct {
	ID        string    `json:"-" bson:"_id"`
	Token     string    `bson:"token"`
	ExpiresAt time.Time `json:"-" bson:"expires_at"`
	IsUse     bool      `json:"-" bson:"is_use"`
}
