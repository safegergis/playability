package db

import (
	"database/sql"
	"errors"
	"log"
	"playability/types"
	"time"
)

// InsertVerification inserts a new verification code into the database
func (m DatabaseModel) InsertVerification(verification types.VerificationRow) error {
	// Check if database connection is valid
	if m.DB == nil {
		log.Printf("[InsertVerification] Database connection is nil")
		return errors.New("database connection is nil")
	}

	log.Printf("[InsertVerification] Inserting verification for email: %s (type: %d)", verification.Email, verification.Type)

	// Delete any existing verification codes for this email and type
	deleteQuery := `DELETE FROM verifications WHERE email = $1 AND type = $2`
	_, err := m.DB.Exec(deleteQuery, verification.Email, verification.Type)
	if err != nil {
		log.Printf("[InsertVerification] Error deleting old verifications for %s: %v", verification.Email, err)
		return errors.New("internal server error")
	}

	// Insert the new verification code
	insertQuery := `INSERT INTO verifications (email, code, expires_at, type) VALUES ($1, $2, $3, $4)`
	_, err = m.DB.Exec(insertQuery, verification.Email, verification.Code, verification.ExpiresAt, verification.Type)
	if err != nil {
		log.Printf("[InsertVerification] Error inserting verification for %s: %v", verification.Email, err)
		return errors.New("internal server error")
	}

	log.Printf("[InsertVerification] Successfully inserted verification for: %s", verification.Email)
	return nil
}

// GetVerification retrieves a verification code from the database
func (m DatabaseModel) GetVerification(email string, mailType int) (types.VerificationRow, error) {
	// Check if database connection is valid
	if m.DB == nil {
		log.Printf("[GetVerification] Database connection is nil")
		return types.VerificationRow{}, errors.New("database connection is nil")
	}

	log.Printf("[GetVerification] Fetching verification for email: %s (type: %d)", email, mailType)

	query := `SELECT email, code, expires_at, type FROM verifications WHERE email = $1 AND type = $2`

	var verification types.VerificationRow
	err := m.DB.QueryRow(query, email, mailType).Scan(
		&verification.Email,
		&verification.Code,
		&verification.ExpiresAt,
		&verification.Type,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("[GetVerification] No verification found for email: %s (type: %d)", email, mailType)
			return types.VerificationRow{}, errors.New("verification not found")
		}
		log.Printf("[GetVerification] Error querying verification for %s: %v", email, err)
		return types.VerificationRow{}, errors.New("internal server error")
	}

	log.Printf("[GetVerification] Successfully fetched verification for: %s", email)
	return verification, nil
}

// DeleteVerification removes a verification code from the database
func (m DatabaseModel) DeleteVerification(email string, mailType int) error {
	// Check if database connection is valid
	if m.DB == nil {
		log.Printf("[DeleteVerification] Database connection is nil")
		return errors.New("database connection is nil")
	}

	log.Printf("[DeleteVerification] Deleting verification for email: %s (type: %d)", email, mailType)

	deleteQuery := `DELETE FROM verifications WHERE email = $1 AND type = $2`
	result, err := m.DB.Exec(deleteQuery, email, mailType)
	if err != nil {
		log.Printf("[DeleteVerification] Error deleting verification for %s: %v", email, err)
		return errors.New("internal server error")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("[DeleteVerification] Error checking rows affected: %v", err)
		return errors.New("internal server error")
	}

	if rowsAffected == 0 {
		log.Printf("[DeleteVerification] No verification found to delete for email: %s (type: %d)", email, mailType)
		return errors.New("verification not found")
	}

	log.Printf("[DeleteVerification] Successfully deleted verification for: %s", email)
	return nil
}

// CleanExpiredVerifications removes all expired verification codes from the database
func (m DatabaseModel) CleanExpiredVerifications() error {
	// Check if database connection is valid
	if m.DB == nil {
		log.Printf("[CleanExpiredVerifications] Database connection is nil")
		return errors.New("database connection is nil")
	}

	log.Printf("[CleanExpiredVerifications] Cleaning expired verifications")

	deleteQuery := `DELETE FROM verifications WHERE expires_at < $1`
	result, err := m.DB.Exec(deleteQuery, time.Now())
	if err != nil {
		log.Printf("[CleanExpiredVerifications] Error deleting expired verifications: %v", err)
		return errors.New("internal server error")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("[CleanExpiredVerifications] Error checking rows affected: %v", err)
		return errors.New("internal server error")
	}

	log.Printf("[CleanExpiredVerifications] Successfully deleted %d expired verifications", rowsAffected)
	return nil
}

// MarkUserAsVerified sets the verified flag to true for a user
func (m DatabaseModel) MarkUserAsVerified(email string) error {
	// Check if database connection is valid
	if m.DB == nil {
		log.Printf("[MarkUserAsVerified] Database connection is nil")
		return errors.New("database connection is nil")
	}

	log.Printf("[MarkUserAsVerified] Marking user as verified: %s", email)

	updateQuery := `UPDATE users SET verified = true WHERE email = $1`
	result, err := m.DB.Exec(updateQuery, email)
	if err != nil {
		log.Printf("[MarkUserAsVerified] Error updating user: %v", err)
		return errors.New("internal server error")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("[MarkUserAsVerified] Error checking rows affected: %v", err)
		return errors.New("internal server error")
	}

	if rowsAffected == 0 {
		log.Printf("[MarkUserAsVerified] No user found with email: %s", email)
		return errors.New("user not found")
	}

	log.Printf("[MarkUserAsVerified] Successfully marked user as verified: %s", email)
	return nil
}

// UpdateUserPassword updates the password hash for a user
func (m DatabaseModel) UpdateUserPassword(email string, newHash string) error {
	// Check if database connection is valid
	if m.DB == nil {
		log.Printf("[UpdateUserPassword] Database connection is nil")
		return errors.New("database connection is nil")
	}

	log.Printf("[UpdateUserPassword] Updating password for user: %s", email)

	updateQuery := `UPDATE users SET hash = $1 WHERE email = $2`
	result, err := m.DB.Exec(updateQuery, newHash, email)
	if err != nil {
		log.Printf("[UpdateUserPassword] Error updating password: %v", err)
		return errors.New("internal server error")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("[UpdateUserPassword] Error checking rows affected: %v", err)
		return errors.New("internal server error")
	}

	if rowsAffected == 0 {
		log.Printf("[UpdateUserPassword] No user found with email: %s", email)
		return errors.New("user not found")
	}

	log.Printf("[UpdateUserPassword] Successfully updated password for: %s", email)
	return nil
}

// GetUserByEmail retrieves user information from the database by email
func (m DatabaseModel) GetUserByEmail(email string) (types.UserRow, error) {
	// Check if database connection is valid
	if m.DB == nil {
		log.Printf("[GetUserByEmail] Database connection is nil")
		return types.UserRow{}, errors.New("database connection is nil")
	}

	log.Printf("[GetUserByEmail] Fetching user with email: %s", email)

	query := `SELECT id, username, email, hash, verified, num_reports FROM users WHERE email = $1`

	var user types.UserRow
	err := m.DB.QueryRow(query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Hash, &user.Verified, &user.NumOfReports)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("[GetUserByEmail] No user found with email: %s", email)
			return types.UserRow{}, errors.New("user not found")
		}
		log.Printf("[GetUserByEmail] Error querying user with email %s: %v", email, err)
		return types.UserRow{}, errors.New("internal server error")
	}

	log.Printf("[GetUserByEmail] Successfully fetched user: %s (email: %s)", user.Username, email)
	return user, nil
}
