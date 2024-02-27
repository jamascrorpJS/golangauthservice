package token

import (
	"errors"
	"os"
	"time"

	"github.com/Microsoft/go-winio/pkg/guid"
	"github.com/golang-jwt/jwt/v5"
)

type Claim struct {
	ID      guid.GUID
	TokenID guid.GUID
	jwt.RegisteredClaims
}

const TokenTTL = time.Second * 19

var ErrToken = "TokenFault"
var ErrHash = "No true hash"

func CreateAccessToken(id guid.GUID, tokenID guid.GUID) (string, error) {
	secret := os.Getenv("SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, Claim{
		ID:      id,
		TokenID: tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})
	signToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", errors.New(ErrToken)
	}
	return signToken, nil
}

func ParseAccessToken(token string) (jwt.MapClaims, error) {
	secret := os.Getenv("SECRET")
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(ErrHash)
		}
		return []byte(secret), nil
	})
	jwt := t.Claims.(jwt.MapClaims)

	if err != nil {
		return jwt, err
	}

	return jwt, nil
}
