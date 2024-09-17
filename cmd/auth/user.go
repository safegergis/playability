package auth

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"golang.org/x/crypto/bcrypt"
)

func GenerateAuthToken() *jwtauth.JWTAuth {
	return jwtauth.New("HS256", []byte("secret"), nil)
}

func getHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (server *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user UserRegister
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hash, err := getHash(user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3)`
	result, err := server.DB.Exec(query, user.Username, user.Email, hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	recordId, _ := result.LastInsertId()
	response := Response{
		ID: int(recordId),
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
