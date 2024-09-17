package auth

import (
	"github.com/go-chi/jwtauth/v5"
	"golang.org/x/crypto/bcrypt"
)

func GenerateAuthToken() *jwtauth.JWTAuth {
	return jwtauth.New("HS256", []byte("secret"), nil)
}

func GetHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
