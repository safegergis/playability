package db

import (
	"database/sql"
	"errors"
	"log"
	"playability/auth"
	"playability/types"

	"github.com/lib/pq"
)

func (m DatabaseModel) InsertUser(user types.UserRegister) error {
	if m.DB == nil {
		return errors.New("database connection is nil")
	}
	var existingID int
	checkQuery := `SELECT id FROM users WHERE email = $1 OR username = $2`
	err := m.DB.QueryRow(checkQuery, user.Email, user.Username).Scan(&existingID)

	if err != sql.ErrNoRows {
		if err != nil {
			log.Println("Error checking email/username:", err)
			return errors.New("internal server error")
		}
		// If we found a matching record, determine which field caused the conflict
		if existingID != 0 {
			checkEmailQuery := `SELECT id FROM users WHERE email = $1`
			err := m.DB.QueryRow(checkEmailQuery, user.Email).Scan(&existingID)
			if err == nil {
				log.Println("email error")
				return errors.New("email is already in use")
			} else {
				log.Println("username error")
				return errors.New("username is already in use")
			}

		}
	}

	// Proceed to create the user since both email and username are unique

	hash, err := auth.GetHash(user.Password)
	if err != nil {
		log.Println("Error hashing password:", err)
		return errors.New("internal server error")
	}
	insertQuery := `INSERT INTO users (username, email, hash, num_reports) VALUES ($1, $2, $3, $4)`
	_, err = m.DB.Exec(insertQuery, user.Username, user.Email, hash, 0)
	if err != nil {
		// Check if the error is due to a unique constraint violation
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" { // unique_violation
			if pqErr.Constraint == "users_email_key" {
				return errors.New("email is already in use")
			} else if pqErr.Constraint == "users_username_key" {
				return errors.New("username is already in use")
			} else {
				return errors.New("user already exists")
			}
		}
		log.Println("Error inserting user:", err)
		return errors.New("internal server error")
	}

	return nil
}
