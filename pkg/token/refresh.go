package token

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

type Token struct {
	ExpiresAt int64
}

const TokenTTLs = time.Minute

func CreateRefershToken() string {
	t := Token{
		ExpiresAt: time.Now().Add(TokenTTLs).Unix(),
	}
	json, err := json.Marshal(t)
	if err != nil {

	}
	encode := base64.StdEncoding.EncodeToString(json)
	return encode
}
