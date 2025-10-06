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
	log.Printf("[PostCreateUser] Received user registration request from %s", r.RemoteAddr)

	var user types.UserRegister
	// Decode the request body into a UserRegister struct
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("[PostCreateUser] Invalid JSON in request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if user.Username == "" || user.Email == "" || user.Password == "" {
		log.Printf("[PostCreateUser] Missing required fields")
		http.Error(w, "username, email, and password are required", http.StatusBadRequest)
		return
	}

	// Convert username and email to lowercase
	user.Username = strings.ToLower(user.Username)
	user.Email = strings.ToLower(user.Email)
	// Note: Password remains case-sensitive

	log.Printf("[PostCreateUser] Attempting to create user: %s (email: %s)", user.Username, user.Email)

	// Attempt to insert the new user into the database
	err := env.DB.InsertUser(user)
	if err != nil {
		// Handle specific error cases
		if err.Error() == "email is already in use" {
			log.Printf("[PostCreateUser] Email already in use: %s", user.Email)
			http.Error(w, "email is already in use", http.StatusConflict)
			return
		} else if err.Error() == "username is already in use" {
			log.Printf("[PostCreateUser] Username already in use: %s", user.Username)
			http.Error(w, "username is already in use", http.StatusConflict)
			return
		} else {
			log.Printf("[PostCreateUser] Error creating user %s: %v", user.Username, err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	// Get the newly created user to send verification email
	// userRow, err := env.DB.GetUserByEmail(user.Email)
	if err != nil {
		log.Printf("[PostCreateUser] Error getting user for verification email: %v", err)
		// User is created, so we can still return success
		w.WriteHeader(http.StatusCreated)
		return
	}

	// Send verification email
	// verification, err := env.SendVerifyEmail(userRow)
	if err != nil {
		log.Printf("[PostCreateUser] Error sending verification email: %v", err)
		// User is created, so we can still return success
		w.WriteHeader(http.StatusCreated)
		return
	}

	// Store verification in database
	// err = env.DB.InsertVerification(verification)
	if err != nil {
		log.Printf("[PostCreateUser] Error storing verification: %v", err)
		// User is created, so we can still return success
		w.WriteHeader(http.StatusCreated)
		return
	}

	// If successful, return 201 Created status
	log.Printf("[PostCreateUser] Successfully created user: %s (email: %s) and sent verification email", user.Username, user.Email)
	w.WriteHeader(http.StatusCreated)
}

// PostLoginUser handles user login attempts
func (env *Env) PostLoginUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("[PostLoginUser] Login attempt from %s", r.RemoteAddr)

	var user types.UserLogin
	// Decode the request body into a UserLogin struct
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("[PostLoginUser] Invalid JSON in request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if user.Email == "" || user.Password == "" {
		log.Printf("[PostLoginUser] Missing email or password")
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	log.Printf("[PostLoginUser] Login attempt for email: %s", user.Email)

	// Check if the user credentials are valid
	id, valid, err := env.DB.CheckUser(user.Email, user.Password)
	if err != nil {
		log.Printf("[PostLoginUser] Database error checking credentials for %s: %v", user.Email, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if !valid {
		log.Printf("[PostLoginUser] Invalid credentials for email: %s", user.Email)
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	// Get user details to check verification status
	userRow, err := env.DB.QueryUser(id)
	if err != nil {
		log.Printf("[PostLoginUser] Error querying for user info %d: %v", id, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Check if user is verified
	if userRow.Verified {
		log.Printf("[PostLoginUser] User not verified: %s", user.Email)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "email not verified",
			"email": user.Email,
		})
		return
	}

	// Create a JWT token for the authenticated user
	token, err := auth.CreateToken(id)
	if err != nil {
		log.Printf("[PostLoginUser] Error creating token for user ID %d: %v", id, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	userInfo := types.UserInfo{
		ID:           userRow.ID,
		Username:     userRow.Username,
		NumOfReports: userRow.NumOfReports,
	}

	loginResponse := types.LoginResponse{
		Token: token,
		User:  userInfo,
	}
	// Return the token to the client
	log.Printf("[PostLoginUser] Successfully authenticated user ID: %d", id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(loginResponse)
}

// GetUserHandler retrieves user information based on the provided ID
func (env *Env) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the user ID from the URL parameters
	id := chi.URLParam(r, "id")
	log.Printf("[GetUserHandler] Request for user ID: %s from %s", id, r.RemoteAddr)

	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("[GetUserHandler] Invalid user ID format: %s", id)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Query the database for the user
	user, err := env.DB.QueryUser(idInt)
	if err != nil {
		log.Printf("[GetUserHandler] Error querying user ID %d: %v", idInt, err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Encode and return the user information as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("[GetUserHandler] Error encoding user ID %d to JSON: %v", idInt, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("[GetUserHandler] Successfully returned user: %s (ID: %d)", user.Username, idInt)
}
