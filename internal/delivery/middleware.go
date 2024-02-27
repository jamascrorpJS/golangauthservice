package delivery

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jamascrorpJS/auth/pkg/token"
)

type ID struct {
	TokenID string
	UserID  string
}

var ErrAccess = ""

func AccessToken(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader, err := r.Cookie("access_token")
		token, err := token.ParseAccessToken(authHeader.Value)
		if errors.Is(err, jwt.ErrTokenExpired) {
			id := ID{
				TokenID: token["TokenID"].(string),
				UserID:  token["ID"].(string),
			}
			ctx := context.WithValue(context.Background(), "id", id)
			r = r.WithContext(ctx)
			next(w, r)
			return
		} else if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error:%v", errors.New(ErrAccess))
			return
		}
		id := ID{
			TokenID: token["TokenID"].(string),
			UserID:  token["ID"].(string),
		}
		ctx := context.WithValue(context.Background(), "id", id)
		r = r.WithContext(ctx)
		next(w, r)
	}
}
