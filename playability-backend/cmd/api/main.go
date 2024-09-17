package main

import (
	"net/http"
	"playability/db"
	"playability/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	_ "github.com/lib/pq"
)

type Env struct {
	handlers *handlers.Env
	router   *chi.Mux
}

func (env *Env) MountMiddleware() {
	env.router.Use(middleware.Logger)
	env.router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
	}))
}
func (env *Env) MountHandlers() {
	env.router.Get("/search", handlers.GetSearchHandler)
	env.router.Get("/games", env.handlers.GetGamesHandler)

	env.router.Route("/user", func(r chi.Router) {
		//r.Post("/login", server.LoginUser)
		r.Post("/register", env.handlers.PostUser)

		//r.Group(func(r chi.Router) {
		//r.Use(jwtauth.Verifier(server.Authtoken))
		//r.Post("/report", server.CreateReport)
		//})
	})

}
func main() {

	database := db.InitDB()

	// Set up the env
	env := &Env{
		router: chi.NewRouter(),
		handlers: &handlers.Env{
			DB: db.DatabaseModel{DB: database},
		},
	}

	env.MountMiddleware()
	env.MountHandlers()

	// Start the server on port 8080
	http.ListenAndServe(":8080", env.router)
}
