package auth

import (
	"log"
	"os"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// GenerateAuthToken creates a new JWT authentication token
func GenerateAuthToken() *jwtauth.JWTAuth {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	return jwtauth.New("HS256", []byte(secret), nil)
}

// GetHash generates a bcrypt hash from a password
func GetHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPassword compares a password with a bcrypt hash
func CheckPassword(password string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// CreateToken generates a new JWT token for a user
func CreateToken(userid int) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userid,
		"iss": "playability",
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
	})
	tokenString, err := claims.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
