# Backend Unit Testing Plan - Progress Tracker

## Overview
Comprehensive unit test coverage for Playability backend with 80%+ code coverage goal.

---

## ✅ COMPLETED TESTS (12/17 - 71%)

### 1. ✅ pkg/calc/score_test.go (100% coverage)
- `CalculateAccessibilityScore` - 11 test cases
- `toFixed` - 7 test cases
- `round` - 9 test cases
- Edge cases: empty arrays, precision, negative numbers
- **Status**: COMPLETE

### 2. ✅ pkg/calc/features_test.go (100% coverage)
- `CalculateFeatureScore` - 7 test cases
- `featureNameToDisplayName` - 5 test cases
- `determineConsensus` - 3 test cases
- `calculatePercentages` - 2 test cases
- `calculateTotal` - 1 test case
- **Status**: COMPLETE

### 3. ✅ auth/user_test.go (77.3% coverage)
- `GetHash` - 5 test cases (including uniqueness test)
- `CheckPassword` - 5 test cases
- `CreateToken` - 4 test cases + JWT claims validation
- `GenerateRandomString` - 4 test cases + uniqueness test
- `GenerateAuthToken` - 2 test cases
- Benchmark tests included
- **Status**: COMPLETE

### 4. ✅ db/user_test.go
Tests with sqlmock:
- `InsertUser` - 5 test cases (success, duplicate email, duplicate username, constraint violation, nil DB)
- `CheckUser` - 4 test cases (valid credentials with bcrypt, user not found, invalid password, nil DB)
- `QueryUser` - 4 test cases (found, not found, DB error, nil DB)
- **Status**: COMPLETE

### 5. ✅ db/verification_test.go
Tests with sqlmock:
- `InsertVerification` - 3 test cases (success, nil DB, delete error)
- `GetVerification` - 3 test cases (found, not found, nil DB)
- `DeleteVerification` - 3 test cases (success, not found, nil DB)
- `CleanExpiredVerifications` - 2 test cases (success, nil DB)
- `MarkUserAsVerified` - 3 test cases (success, user not found, nil DB)
- `UpdateUserPassword` - 3 test cases (success, user not found, nil DB)
- `GetUserByEmail` - 3 test cases (found, not found, nil DB)
- **Status**: COMPLETE

### 6. ✅ db/reports_test.go
Tests with sqlmock:
- `InsertReport` - 4 test cases (success, duplicate, nil report, nil DB)
- `InsertReportSummary` - 3 test cases (success, nil summary, nil DB)
- `QueryReportCards` - 2 test cases (found, nil DB)
- `QueryAccessibilityScores` - 3 test cases (found, empty, nil DB)
- `GetReportCount` - 2 test cases (found, nil DB)
- `QueryReportsForSummarization` - 2 test cases (found, nil DB)
- `QueryReportSummary` - 3 test cases (found, not found, nil DB)
- `QueryUserReports` - 2 test cases (found with JOIN, nil DB)
- **Status**: COMPLETE

### 7. ✅ db/features_test.go
Tests with sqlmock:
- `QueryFeatureReports` - 7 test cases (success, empty, query error, scan error, iteration error, nil DB, all value types)
- **Status**: COMPLETE

### 8. ✅ db/game_test.go
Tests with sqlmock:
- `InsertGame` - 6 test cases (new game, update on conflict, invalid JSON, marshaling, DB error, empty platforms)
- `QueryGame` - 5 test cases (found, not found, DB error, invalid JSON platforms, null timestamps)
- `QueryFeaturedGames` - 8 test cases (fresh cache, stale cache, empty cache, cache age error, query error, scan error, iteration error, nil DB)
- `UpsertFeaturedGames` - 9 test cases (successful transaction, empty list, position ordering, transaction error, clear error, prepare error, insert error, commit error, nil DB)
- **Status**: COMPLETE

### 9. ✅ handlers/users_test.go
Tests with httptest:
- `PostCreateUser` - 9 test cases (success, invalid JSON, missing fields x3, duplicate email, duplicate username, DB error, lowercase normalization)
- `PostLoginUser` - 8 test cases (success, invalid JSON, missing fields x2, invalid credentials x2, DB errors x2)
- `GetUserHandler` - 4 test cases (success, invalid ID format, not found, DB error)
- **Status**: COMPLETE
- **Notes**: Tests document bugs in code (inverted Verified logic, missing 'verified' column in QueryUser)

### 10. ✅ handlers/score_test.go
Tests with httptest:
- `GetScoreHandler` - 5 test cases (success, invalid ID format, DB error, not enough reports, single score)
- **Status**: COMPLETE

### 11. ✅ handlers/features_test.go
Tests with httptest:
- `GetFeatureReportsHandler` - 5 test cases (success, invalid ID format, DB error, not enough reports, single report)
- **Status**: COMPLETE

### 12. ✅ handlers/games_test.go
Tests with httptest:
- `GetGamesHandler` - 3 test cases (cache hit, cache miss, DB error)
- `GetFeaturedHandler` - 5 test cases (fresh cache, stale cache, empty cache, DB error, marshal error)
- `GetSearchHandler` - 2 test cases (search term, empty term)
- **Status**: COMPLETE (partial - fetch package not mocked)
- **Notes**: External fetch calls not mocked, tests focus on DB interactions

---

## 📋 REMAINING TESTS (5/17 - 29%)

### HTTP Handlers Layer (2 remaining)

#### 13. ⏳ handlers/reports_test.go
Use httptest for:
- `PostReportHandler` - valid report with JWT, AI moderation pass/fail, duplicate reports, 10th report triggers summarization
- `GetReportCardsHandler` - valid game ID, invalid ID conversion, empty results
- `GetReportSummaryHandler` - existing summary, no summary (404), empty summary field
- `GetUserReportsHandler` - JWT extraction, valid reports with game data, empty results

#### 14. ⏳ handlers/verification_test.go
Use httptest for:
- `PostVerifyEmail` - valid code, expired code, invalid code, email normalization
- `PostRequestPasswordReset` - existing user, non-existent user (security), email sending
- `PostResetPassword` - valid reset, expired code, invalid code, password hashing
- `PostResendVerification` - unverified user, already verified user, non-existent user

---

### External Integration Layer (3 remaining)

#### 15. ⏳ pkg/ai/moderation_test.go
Mock Claude API:
- `Moderation` - violation detection (various categories), safe content, API errors, empty/nil reports, JSON parsing

#### 16. ⏳ pkg/ai/summarization_test.go
Mock Claude API:
- `SummarizeReports` - summary generation with multiple reports, platform tracking, empty reports, API errors

#### 17. ⏳ pkg/fetch/fetch_test.go
Mock HTTP calls to IGDB/PCGamingWiki:
- `GetSearch` - valid search, empty results, IGDB API errors
- `GetGame` - with Steam ID + PCGamingWiki data, without Steam, cover art, accessibility features
- `GetFeaturedGames` - popularity primitives fetch, game details enrichment
- `getFeaturedGameDetails` - valid game, not found, cover fetching
- `makeIgdbRequest` - successful request, auth headers, error responses, status codes

---

## 🎯 Current Coverage Stats

```
✅ pkg/calc:      100.0%
✅ auth:           77.3%
✅ db:             82.5%
⏳ handlers:       0%
⏳ pkg/ai:         0%
⏳ pkg/fetch:      0%
```

**Overall Progress: 12/17 files (71%)**

---

## 🔧 Testing Tools & Dependencies

### Installed
- ✅ `github.com/DATA-DOG/go-sqlmock` - Database mocking
- ✅ `github.com/lib/pq` - PostgreSQL driver (for error types)

### To Install (when needed)
- `net/http/httptest` - Built-in, no install needed
- Mock HTTP client for external API testing (custom or `github.com/jarcoal/httpmock`)

---

## 📝 Key Testing Patterns Established

### Database Tests (sqlmock)
```go
db, mock, err := sqlmock.New()
defer db.Close()
dbModel := DatabaseModel{DB: db}

mock.ExpectQuery("SELECT...").
    WithArgs(expectedArgs...).
    WillReturnRows(sqlmock.NewRows([]string{"col1"}).AddRow(val1))

// Test the function
result, err := dbModel.SomeMethod(args)

// Verify expectations
if err := mock.ExpectationsWereMet(); err != nil {
    t.Errorf("Unfulfilled expectations: %v", err)
}
```

### HTTP Handler Tests (httptest)
```go
req := httptest.NewRequest("POST", "/endpoint", body)
w := httptest.NewRecorder()

env := &handlers.Env{DB: mockDB, MS: mockMailService}
env.SomeHandler(w, req)

resp := w.Result()
// Assert status code, body, headers
```

### Mock External APIs
```go
// Use custom http.RoundTripper or httpmock
// Mock Claude AI responses
// Mock IGDB/PCGamingWiki responses
```

---

## 🎯 Next Steps

1. **IMMEDIATE**: Start handler tests - `handlers/users_test.go` (high value, user-facing)
2. **SHORT TERM**: Complete remaining handler tests
3. **MEDIUM TERM**: AI integration tests (moderation, summarization)
4. **LONG TERM**: External API integration tests (IGDB/PCGamingWiki)

---

## 📊 Coverage Goals

- **Minimum**: 80% overall
- **Critical paths**: 95%+ (auth, db user/reports, moderation)
- **Target**: 85-90% overall

---

## 🚀 Commands

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./pkg/calc -v
go test ./auth -v
go test ./db -v

# Check coverage
go test ./pkg/calc -cover
go test ./auth -cover
go test ./db -cover

# Generate coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

**Last Updated**: In progress - 12/17 tests complete (71%)
**Next Test File**: handlers/reports_test.go or handlers/verification_test.go
