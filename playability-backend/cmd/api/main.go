package main

import (
	"fmt"
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

type Env struct {
	handlers *handlers.Env
	router   *chi.Mux
	authtoken *jwtauth.JWTAuth
}

func (env *Env) MountMiddleware() {
	env.router.Use(middleware.Logger)
	env.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	}))
}
func (env *Env) MountHandlers() {
	env.router.Get("/search", handlers.GetSearchHandler)
	env.router.Get("/games", env.handlers.GetGamesHandler)

	env.router.Route("/user", func(r chi.Router) {
		r.Post("/login", env.handlers.PostLoginUser)
		r.Post("/register", env.handlers.PostCreateUser)

		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(env.authtoken))
			r.Use(jwtauth.Authenticator(env.authtoken))


			r.Get("/report", func(w http.ResponseWriter, r *http.Request) {
				_, claims, err := jwtauth.FromContext(r.Context())
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				userID := claims["sub"].(string)
				fmt.Fprintf(w, "Authenticated user ID: %v", userID)
			})
		})
	})

}
func main() {

	database := db.InitDB()
	authtoken := auth.GenerateAuthToken()
	// Set up the env
	env := &Env{
		router: chi.NewRouter(),
		handlers: &handlers.Env{
			DB: db.DatabaseModel{DB: database},
		},
		authtoken: authtoken,
	}

	env.MountMiddleware()
	env.MountHandlers()

	// Start the server on port 8080
	http.ListenAndServe(":8080", env.router)
}
