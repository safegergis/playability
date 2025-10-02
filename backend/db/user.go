package db

import (
	"database/sql"
	"errors"
	"log"
	"playability/auth"
	"playability/types"
	"strconv"
	"strings"

	"github.com/lib/pq"
)

// InsertUser inserts a new user into the database
// It checks for existing email and username before insertion
func (m DatabaseModel) InsertUser(user types.UserRegister) error {
	// Check if database connection is valid
	if m.DB == nil {
		return errors.New("database connection is nil")
	}

	// Check if email or username already exists
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

	// Hash the password
	hash, err := auth.GetHash(user.Password)
	if err != nil {
		log.Println("Error hashing password:", err)
		return errors.New("internal server error")
	}

	// Insert the new user into the database
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

// CheckUser verifies user credentials and returns user ID if valid
func (m DatabaseModel) CheckUser(email string, password string) (string, bool, error) {
	// Check if database connection is valid
	if m.DB == nil {
		return "", false, errors.New("database connection is nil")
	}

	var hash string
	email = strings.ToLower(email)

	// Retrieve the hash for the given email
	hashQuery := `SELECT hash FROM users WHERE email = $1`
	err := m.DB.QueryRow(hashQuery, email).Scan(&hash)
	if err != nil {
		return "", false, nil
	}

	// Check if the provided password matches the stored hash
	err = auth.CheckPassword(password, hash)
	if err != nil {	
		return "", false, nil
	}

	// Retrieve the user ID
	idQuery := `SELECT id FROM users WHERE email = $1`
	var id int
	err = m.DB.QueryRow(idQuery, email).Scan(&id)
	if err != nil {
		log.Println("Error checking user:", err)
		return "", false, err
	}

	return strconv.Itoa(id), true, nil
}

// QueryUser retrieves user information from the database
func (m DatabaseModel) QueryUser(userID int) (types.UserRow, error) {
	// Check if database connection is valid
	if m.DB == nil {
		return types.UserRow{}, errors.New("database connection is nil")
	}

	query := `SELECT * FROM users WHERE id = $1`

	var user types.UserRow
	err := m.DB.QueryRow(query, userID).Scan(&user.ID, &user.Username, &user.Email, &user.Hash, &user.NumOfReports)
	if err != nil {
		return types.UserRow{}, err
	}

	return user, nil
}
