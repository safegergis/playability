package db

import (
	"database/sql"
	"errors"
	"log"
	"playability/auth"
	"playability/types"
	"strings"

	"github.com/lib/pq"
)

// InsertUser inserts a new user into the database
// It checks for existing email and username before insertion
func (m DatabaseModel) InsertUser(user types.UserRegister) error {
	// Check if database connection is valid
	if m.DB == nil {
		log.Printf("[InsertUser] Database connection is nil")
		return errors.New("database connection is nil")
	}

	log.Printf("[InsertUser] Attempting to insert user: %s (email: %s)", user.Username, user.Email)

	// Check if email or username already exists
	var existingID int
	checkQuery := `SELECT id FROM users WHERE email = $1 OR username = $2`
	err := m.DB.QueryRow(checkQuery, user.Email, user.Username).Scan(&existingID)

	if err != sql.ErrNoRows {
		if err != nil {
			log.Printf("[InsertUser] Error checking existing email/username for %s: %v", user.Email, err)
			return errors.New("internal server error")
		}
		// If we found a matching record, determine which field caused the conflict
		if existingID != 0 {
			checkEmailQuery := `SELECT id FROM users WHERE email = $1`
			err := m.DB.QueryRow(checkEmailQuery, user.Email).Scan(&existingID)
			if err == nil {
				log.Printf("[InsertUser] Email already in use: %s", user.Email)
				return errors.New("email is already in use")
			} else {
				log.Printf("[InsertUser] Username already in use: %s", user.Username)
				return errors.New("username is already in use")
			}
		}
	}

	// Proceed to create the user since both email and username are unique

	// Hash the password
	hash, err := auth.GetHash(user.Password)
	if err != nil {
		log.Printf("[InsertUser] Error hashing password for user %s: %v", user.Username, err)
		return errors.New("internal server error")
	}

	// Insert the new user into the database
	insertQuery := `INSERT INTO users (username, email, hash, num_reports) VALUES ($1, $2, $3, $4)`
	_, err = m.DB.Exec(insertQuery, user.Username, user.Email, hash, 0)
	if err != nil {
		// Check if the error is due to a unique constraint violation
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" { // unique_violation
			log.Printf("[InsertUser] Unique constraint violation: %s (constraint: %s)", pqErr.Code, pqErr.Constraint)
			if pqErr.Constraint == "unique_email" {
				return errors.New("email is already in use")
			} else if pqErr.Constraint == "unique_username" {
				return errors.New("username is already in use")
			} else {
				return errors.New("user already exists")
			}
		}
		log.Printf("[InsertUser] Error inserting user %s: %v", user.Username, err)
		return errors.New("internal server error")
	}

	log.Printf("[InsertUser] Successfully created user: %s (email: %s)", user.Username, user.Email)
	return nil
}

// CheckUser verifies user credentials and returns user ID if valid
func (m DatabaseModel) CheckUser(email string, password string) (int, bool, error) {
	// Check if database connection is valid
	if m.DB == nil {
		log.Printf("[CheckUser] Database connection is nil")
		return 0, false, errors.New("database connection is nil")
	}

	var hash string
	email = strings.ToLower(email)

	log.Printf("[CheckUser] Attempting login for email: %s", email)

	// Retrieve the hash for the given email
	hashQuery := `SELECT hash FROM users WHERE email = $1`
	err := m.DB.QueryRow(hashQuery, email).Scan(&hash)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("[CheckUser] No user found with email: %s", email)
		} else {
			log.Printf("[CheckUser] Error querying user by email %s: %v", email, err)
		}
		return 0, false, nil
	}

	// Check if the provided password matches the stored hash
	err = auth.CheckPassword(password, hash)
	if err != nil {
		log.Printf("[CheckUser] Invalid password for email: %s", email)
		return 0, false, nil
	}

	// Retrieve the user ID
	idQuery := `SELECT id FROM users WHERE email = $1`
	var id int
	err = m.DB.QueryRow(idQuery, email).Scan(&id)
	if err != nil {
		log.Printf("[CheckUser] Error retrieving user ID for %s: %v", email, err)
		return 0, false, err
	}

	log.Printf("[CheckUser] Successful login for user ID: %d", id)
	return id, true, nil
}

// QueryUser retrieves user information from the database
func (m DatabaseModel) QueryUser(userID int) (types.UserRow, error) {
	// Check if database connection is valid
	if m.DB == nil {
		log.Printf("[QueryUser] Database connection is nil")
		return types.UserRow{}, errors.New("database connection is nil")
	}

	log.Printf("[QueryUser] Fetching user with ID: %d", userID)

	query := `SELECT id, username, email, hash, num_reports FROM users WHERE id = $1`

	var user types.UserRow
	err := m.DB.QueryRow(query, userID).Scan(&user.ID, &user.Username, &user.Email, &user.Hash, &user.NumOfReports)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("[QueryUser] No user found with ID: %d", userID)
		} else {
			log.Printf("[QueryUser] Error querying user ID %d: %v", userID, err)
		}
		return types.UserRow{}, err
	}

	log.Printf("[QueryUser] Successfully fetched user: %s (ID: %d)", user.Username, userID)
	return user, nil
}
