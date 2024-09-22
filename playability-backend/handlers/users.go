// Package handlers contains HTTP request handlers for the application
package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"playability/auth"
	"playability/types"

	"github.com/go-chi/chi/v5"
)

// PostCreateUser handles the creation of a new user
func (env *Env) PostCreateUser(w http.ResponseWriter, r *http.Request) {
	var user types.UserRegister
	// Decode the request body into a UserRegister struct
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Convert username and email to lowercase
	user.Username = strings.ToLower(user.Username)
	user.Email = strings.ToLower(user.Email)
	// Note: Password remains case-sensitive

	// Attempt to insert the new user into the database
	err := env.DB.InsertUser(user)
	if err != nil {
		// Handle specific error cases
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

	// If successful, return 201 Created status
	w.WriteHeader(http.StatusCreated)
}

// PostLoginUser handles user login attempts
func (env *Env) PostLoginUser(w http.ResponseWriter, r *http.Request) {
	var user types.UserLogin
	// Decode the request body into a UserLogin struct
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the user credentials are valid
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
	// Create a JWT token for the authenticated user
	token, err := auth.CreateToken(id)
	if err != nil {
		log.Println("Error creating token:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Token:", token)

	// Return the token to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(token))
}

// GetUserHandler retrieves user information based on the provided ID
func (env *Env) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the user ID from the URL parameters
	id := chi.URLParam(r, "id")
	idInt, _ := strconv.Atoi(id)
	// Query the database for the user
	user, err := env.DB.QueryUser(idInt)
	if err != nil {
		log.Println("Error getting user:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Encode and return the user information as JSON
	json.NewEncoder(w).Encode(user)
}
