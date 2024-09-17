package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"playability/types"
)

func (env *Env) PostUser(w http.ResponseWriter, r *http.Request) {
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
