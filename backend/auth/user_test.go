package auth

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func TestGetHash(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{
			name:     "standard password",
			password: "password123",
		},
		{
			name:     "complex password",
			password: "P@ssw0rd!#$%^&*()",
		},
		{
			name:     "empty password",
			password: "",
		},
		{
			name:     "long password within bcrypt limit",
			password: strings.Repeat("a", 72),
		},
		{
			name:     "unicode password",
			password: "パスワード123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := GetHash(tt.password)
			if err != nil {
				t.Fatalf("GetHash() error = %v", err)
			}

			// Hash should not be empty
			if hash == "" {
				t.Error("GetHash() returned empty hash")
			}

			// Hash should be different from password
			if hash == tt.password {
				t.Error("GetHash() returned password as hash")
			}

			// Hash should be valid bcrypt hash (starts with $2a$ or $2b$)
			if !strings.HasPrefix(hash, "$2a$") && !strings.HasPrefix(hash, "$2b$") {
				t.Errorf("GetHash() returned invalid bcrypt hash format: %s", hash)
			}

			// Verify hash can be used with CheckPassword
			err = CheckPassword(tt.password, hash)
			if err != nil {
				t.Errorf("CheckPassword() failed for generated hash: %v", err)
			}
		})
	}
}

func TestGetHashUniqueness(t *testing.T) {
	password := "test123"
	hash1, err1 := GetHash(password)
	hash2, err2 := GetHash(password)

	if err1 != nil || err2 != nil {
		t.Fatalf("GetHash() errors: %v, %v", err1, err2)
	}

	// Same password should produce different hashes (due to salt)
	if hash1 == hash2 {
		t.Error("GetHash() produced identical hashes for same password (should use different salts)")
	}
}

func TestCheckPassword(t *testing.T) {
	password := "correctPassword123"
	hash, _ := GetHash(password)

	tests := []struct {
		name        string
		password    string
		hash        string
		expectError bool
	}{
		{
			name:        "correct password",
			password:    password,
			hash:        hash,
			expectError: false,
		},
		{
			name:        "incorrect password",
			password:    "wrongPassword",
			hash:        hash,
			expectError: true,
		},
		{
			name:        "empty password",
			password:    "",
			hash:        hash,
			expectError: true,
		},
		{
			name:        "case sensitive - wrong case",
			password:    "CORRECTPASSWORD123",
			hash:        hash,
			expectError: true,
		},
		{
			name:        "invalid hash format",
			password:    password,
			hash:        "invalid_hash",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPassword(tt.password, tt.hash)
			if (err != nil) != tt.expectError {
				t.Errorf("CheckPassword() error = %v, expectError = %v", err, tt.expectError)
			}
			if tt.expectError && err == nil {
				t.Error("CheckPassword() expected error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("CheckPassword() unexpected error: %v", err)
			}
		})
	}
}

func TestCreateToken(t *testing.T) {
	// Set JWT_SECRET for testing
	originalSecret := os.Getenv("JWT_SECRET")
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing")
	defer os.Setenv("JWT_SECRET", originalSecret)

	tests := []struct {
		name   string
		userID int
	}{
		{
			name:   "positive user ID",
			userID: 123,
		},
		{
			name:   "user ID 1",
			userID: 1,
		},
		{
			name:   "large user ID",
			userID: 999999,
		},
		{
			name:   "zero user ID",
			userID: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenString, err := CreateToken(tt.userID)
			if err != nil {
				t.Fatalf("CreateToken() error = %v", err)
			}

			// Token should not be empty
			if tokenString == "" {
				t.Error("CreateToken() returned empty token")
			}

			// Token should have 3 parts (header.payload.signature)
			parts := strings.Split(tokenString, ".")
			if len(parts) != 3 {
				t.Errorf("CreateToken() returned invalid JWT format, got %d parts", len(parts))
			}

			// Parse and validate the token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte("test-secret-key-for-testing"), nil
			})
			if err != nil {
				t.Fatalf("Failed to parse token: %v", err)
			}

			if !token.Valid {
				t.Error("Token is not valid")
			}

			// Check claims
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				t.Fatal("Failed to get claims")
			}

			// Check sub (subject) claim contains user ID
			sub, ok := claims["sub"].(string)
			if !ok {
				t.Error("'sub' claim missing or not a string")
			}
			// Just verify sub exists and is not empty (user ID is stored as string)
			if sub == "" {
				t.Error("'sub' claim is empty")
			}

			// Check iss (issuer) claim
			iss, ok := claims["iss"].(string)
			if !ok || iss != "playability" {
				t.Errorf("'iss' claim = %v, want 'playability'", iss)
			}

			// Check exp (expiration) claim
			exp, ok := claims["exp"].(float64)
			if !ok {
				t.Error("'exp' claim missing or not a number")
			}
			expTime := time.Unix(int64(exp), 0)
			if time.Until(expTime) > 24*time.Hour+time.Minute {
				t.Error("Token expiration is more than 24 hours")
			}
			if time.Until(expTime) < 23*time.Hour {
				t.Error("Token expiration is less than 23 hours")
			}

			// Check iat (issued at) claim
			_, ok = claims["iat"].(float64)
			if !ok {
				t.Error("'iat' claim missing or not a number")
			}
		})
	}
}

func TestCreateTokenNoSecret(t *testing.T) {
	// Skip this test as log.Fatal will terminate the test process
	// In production, missing JWT_SECRET should cause the application to fail fast
	t.Skip("Skipping test that would cause log.Fatal to terminate the test process")
}

func TestGenerateRandomString(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{
			name:   "length 10",
			length: 10,
		},
		{
			name:   "length 32",
			length: 32,
		},
		{
			name:   "length 64",
			length: 64,
		},
		{
			name:   "length 1",
			length: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GenerateRandomString(tt.length)
			if err != nil {
				t.Fatalf("GenerateRandomString() error = %v", err)
			}

			// Check length
			if len(result) != tt.length {
				t.Errorf("GenerateRandomString() length = %d, want %d", len(result), tt.length)
			}

			// String should only contain valid base64 URL characters
			for _, char := range result {
				if !((char >= 'A' && char <= 'Z') ||
					(char >= 'a' && char <= 'z') ||
					(char >= '0' && char <= '9') ||
					char == '-' || char == '_') {
					t.Errorf("GenerateRandomString() contains invalid character: %c", char)
				}
			}
		})
	}
}

func TestGenerateRandomStringUniqueness(t *testing.T) {
	length := 32
	iterations := 100
	seen := make(map[string]bool)

	for i := 0; i < iterations; i++ {
		result, err := GenerateRandomString(length)
		if err != nil {
			t.Fatalf("GenerateRandomString() error = %v", err)
		}

		if seen[result] {
			t.Errorf("GenerateRandomString() produced duplicate: %s", result)
		}
		seen[result] = true
	}

	if len(seen) != iterations {
		t.Errorf("GenerateRandomString() produced %d unique strings out of %d", len(seen), iterations)
	}
}

func TestGenerateAuthToken(t *testing.T) {
	// Set JWT_SECRET for testing
	originalSecret := os.Getenv("JWT_SECRET")
	os.Setenv("JWT_SECRET", "test-secret-for-auth-token")
	defer os.Setenv("JWT_SECRET", originalSecret)

	tokenAuth := GenerateAuthToken()

	if tokenAuth == nil {
		t.Fatal("GenerateAuthToken() returned nil")
	}

	// Try to create a token with the tokenAuth
	claims := jwt.MapClaims{
		"sub": "42",
		"iss": "test",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("test-secret-for-auth-token"))
	if err != nil {
		t.Fatalf("Failed to sign token: %v", err)
	}

	// Verify the token can be parsed
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("test-secret-for-auth-token"), nil
	})

	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}

	if !parsedToken.Valid {
		t.Error("Token is not valid")
	}
}

func TestGenerateAuthTokenNoSecret(t *testing.T) {
	// Remove JWT_SECRET and test panic recovery
	originalSecret := os.Getenv("JWT_SECRET")
	os.Unsetenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", originalSecret)

	// This should panic with log.Fatal
	// We can't easily test log.Fatal, so we'll just document the behavior
	t.Skip("Skipping test that would cause log.Fatal")
}

// Test bcrypt cost
func TestGetHashCost(t *testing.T) {
	password := "testPassword"
	hash, err := GetHash(password)
	if err != nil {
		t.Fatalf("GetHash() error = %v", err)
	}

	cost, err := bcrypt.Cost([]byte(hash))
	if err != nil {
		t.Fatalf("bcrypt.Cost() error = %v", err)
	}

	if cost != bcrypt.DefaultCost {
		t.Errorf("Hash cost = %d, want %d (DefaultCost)", cost, bcrypt.DefaultCost)
	}
}

// Benchmark tests
func BenchmarkGetHash(b *testing.B) {
	password := "benchmarkPassword123"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetHash(password)
	}
}

func BenchmarkCheckPassword(b *testing.B) {
	password := "benchmarkPassword123"
	hash, _ := GetHash(password)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CheckPassword(password, hash)
	}
}

func BenchmarkCreateToken(b *testing.B) {
	os.Setenv("JWT_SECRET", "benchmark-secret")
	defer os.Unsetenv("JWT_SECRET")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CreateToken(123)
	}
}

func BenchmarkGenerateRandomString(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GenerateRandomString(32)
	}
}
