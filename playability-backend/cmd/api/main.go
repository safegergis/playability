package main

import (
	"net/http"
	"playability/auth"
	"playability/db"
	"playability/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	_ "github.com/lib/pq"
)

// Env struct holds the application's environment, including handlers, router, and auth token
type Env struct {
	handlers  *handlers.Env
	router    *chi.Mux
	authtoken *jwtauth.JWTAuth
}

// MountMiddleware sets up middleware for the application
func (env *Env) MountMiddleware() {
	// Add logging middleware
	env.router.Use(middleware.Logger)
	// Add CORS middleware with specific options
	env.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	}))
}

// MountHandlers sets up the routes and their corresponding handlers
func (env *Env) MountHandlers() {
	// Set up routes for search and games
	env.router.Get("/search", handlers.GetSearchHandler)
	env.router.Get("/games", env.handlers.GetGamesHandler)

	// Set up routes for reports
	env.router.Route("/reports", func(r chi.Router) {
		r.Get("/cards/{game}", env.handlers.GetReportCardsHandler)
		// Commented out route for accessibility reports
		r.Get("/features/{game}", env.handlers.GetFeatureReportsHandler)
	})

	// Set up routes for user-related actions
	env.router.Route("/user", func(r chi.Router) {
		r.Post("/login", env.handlers.PostLoginUser)
		r.Post("/register", env.handlers.PostCreateUser)
		r.Get("/{id}", env.handlers.GetUserHandler)
		// Commented out route for getting user reports
		// r.Get("/reports/{id}", env.handlers.GetReportHandler)

		// Group of routes that require JWT authentication
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(env.authtoken))
			r.Use(jwtauth.Authenticator(env.authtoken))

			r.Post("/report", env.handlers.PostReportHandler)
		})
	})
}

func main() {
	// Initialize the database
	database := db.InitDB()
	// Generate the authentication token
	authtoken := auth.GenerateAuthToken()

	// Set up the application environment
	env := &Env{
		router: chi.NewRouter(),
		handlers: &handlers.Env{
			DB: db.DatabaseModel{DB: database},
		},
		authtoken: authtoken,
	}

	// Mount middleware and handlers
	env.MountMiddleware()
	env.MountHandlers()

	// Start the server on port 8080
	http.ListenAndServe(":8080", env.router)
}
