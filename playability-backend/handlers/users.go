package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"playability/auth"
	"playability/types"
)

func (env *Env) PostCreateUser(w http.ResponseWriter, r *http.Request) {
	var user types.UserRegister
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Convert all fields of User to lowercase
	user.Username = strings.ToLower(user.Username)
	user.Email = strings.ToLower(user.Email)
	// Note: We don't convert the password to lowercase as it should remain case-sensitive

	err := env.DB.InsertUser(user)
	if err != nil {
		if err.Error() == "email is already in use" {
			http.Error(w, "email is already in use", http.StatusConflict)
			return
		} else if err.Error() == "username is already in use" {
			http.Error(w, "username is already in use", http.StatusConflict)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
}

func (env *Env) PostLoginUser(w http.ResponseWriter, r *http.Request) {
	var user types.UserLogin
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, valid, err := env.DB.CheckUser(user.Email, user.Password)
	if err != nil {
		log.Println("Error checking user:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !valid {
		log.Println("Invalid email or password")
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}
	token, err := auth.CreateToken(id)
	if err != nil {
		log.Println("Error creating token:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Token:", token)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(token))

}
