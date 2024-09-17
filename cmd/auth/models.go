package auth

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

type Server struct {
	Router    *chi.Mux
	DB        *sql.DB
	Authtoken *jwtauth.JWTAuth
}
type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Hash         string `json:"hash"`
	NumOfReports int    `json:"num_of_reports"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegister struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Response struct {
	ID int `json:"id"`
}
