package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"playability/auth"
	"playability/data"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	_ "github.com/lib/pq"
)

type Server struct {
	Router    *chi.Mux
	DB        *sql.DB
	Authtoken *jwtauth.JWTAuth
}

func CreateServer() *Server {
	server := &Server{
		Router:    chi.NewRouter(),
		DB:        data.InitDB(),
		Authtoken: auth.GenerateAuthToken(),
	}
	return server
}
func (server *Server) MountMiddleware() {
	server.Router.Use(middleware.Logger)
	server.Router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
	}))
}
func (server *Server) MountHandlers() {
	server.Router.Get("/search", server.searchHandler)
	server.Router.Get("/games", server.gamesHandler)

	server.Router.Route("/user", func(r chi.Router) {
		//r.Post("/login", server.LoginUser)
		r.Post("/register", server.CreateUser)

		//r.Group(func(r chi.Router) {
		//r.Use(jwtauth.Verifier(server.Authtoken))
		//r.Post("/report", server.CreateReport)
		//})
	})

}
func main() {

	// Set up the router and define routes
	server := CreateServer()

	server.MountMiddleware()
	server.MountHandlers()

	// Start the server on port 8080
	http.ListenAndServe(":8080", server.Router)
}

// searchHandler handles search requests for games
func (server *Server) searchHandler(w http.ResponseWriter, r *http.Request) {
	// Extract search term from query parameters
	searchTerm := r.URL.Query().Get("search")

	// Call getSearch function (not shown) to perform the search
	body, err := data.GetSearch(searchTerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// Log and send the search results
	log.Println("Search results:", string(body))
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

// gamesHandler handles requests for specific game details
func (server *Server) gamesHandler(w http.ResponseWriter, r *http.Request) {
	// Extract game ID from query parameters
	gameID := r.URL.Query().Get("id")
	// Try to query the game from the database
	body, err, found := data.QueryGame(server.DB, gameID)
	if !found {
		// If not found in DB, fetch from external source (getGame function not shown)
		body, err = data.GetGame(gameID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		// Insert the fetched game into the database
		data.InsertGame(server.DB, body)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// Log and send the game details
	log.Println("game details:", string(body))
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
func (server *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user auth.UserRegister
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Convert all fields of User to lowercase
	user.Username = strings.ToLower(user.Username)
	user.Email = strings.ToLower(user.Email)
	// Note: We don't convert the password to lowercase as it should remain case-sensitive



	err := data.InsertUser(server.DB, user)
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
