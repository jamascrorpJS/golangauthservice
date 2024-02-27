package hash

import (
	"golang.org/x/crypto/bcrypt"
)

func Encode(password string) (string, error) {

	byte, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(byte), err
}

func Check(secret string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(secret), []byte(password))
	return err == nil
}
