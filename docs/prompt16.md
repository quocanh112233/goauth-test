# Prompt 16 — Unit Tests: Tất cả frameworks

## Role

Bạn là một **Go Test Engineer** chuyên viết unit test và mock dependencies.

---

## Context

Phase 8 (cuối cùng). Tests cho auth flow mới (Access Token + Refresh Token + Sessions + Auto-refresh). Viết cho Gin trước, port sang 3 approach còn lại.

---

## Dependencies

| Prompt | Đầu ra cần thiết |
|--------|-----------------|
| 05, 07, 10, 12 | 4 apps hoàn chỉnh |
| 14–15 | Google OAuth tất cả |

---

## Yêu cầu

### 1. Test dependency: `github.com/stretchr/testify v1.9.0`

### 2. Mock Repositories

- `mock_user.go`: mock 6 methods UserRepository
- `mock_session.go`: mock 4 methods SessionRepository

### 3. Service Tests — auth_test.go

```go
// Signup
TestSignup_Success
TestSignup_DuplicateEmail
TestSignup_DuplicatePhone
TestSignup_BcryptCost12

// Login (kiểm tra 3 return values + session tạo)
TestLogin_Success
TestLogin_UserNotFound        // error = "email hoặc mật khẩu không đúng"
TestLogin_WrongPassword       // CÙNG error message
TestLogin_AccessTokenExpiry30Min
TestLogin_RefreshTokenRandom  // len >= 64

// Logout
TestLogout_Success           // sessionRepo.DeleteByRefreshToken called
TestLogout_InvalidToken      // graceful

// RefreshAccessToken
TestRefreshAccessToken_Success
TestRefreshAccessToken_SessionNotFound
TestRefreshAccessToken_UserDeleted  // session deleted
```

### 4. OAuth Tests — oauth_test.go

```go
TestUpsertUser_NewUser
TestUpsertUser_ExistingGoogleUser
TestUpsertUser_MergeLocalAccount
```

### 5. Middleware Tests — auth_test.go

```go
TestRequireAuth_ValidAccessToken
TestRequireAuth_NoTokens
TestRequireAuth_ExpiredAccess_ValidRefresh  // auto-refresh
TestRequireAuth_ExpiredAccess_InvalidRefresh
TestRequireAuth_APIRoute_401JSON           // /api/me → 401
```

### 6. Port tests: Fiber/stdlib/Echo

### 7. Makefile targets: `test-gin`, `test-all`, `test-coverage`

---

## Acceptance Criteria

1. `go test ./... -v` pass cho cả 4
2. Coverage service layer ≥ 80%
3. `make test-all` pass

---

## Checklist hoàn thành

### Service Tests

| # | Test | Trạng thái |
|---|------|------------|
| 1 | TestSignup_Success | 🔲 |
| 2 | TestSignup_DuplicateEmail | 🔲 |
| 3 | TestSignup_DuplicatePhone | 🔲 |
| 4 | TestLogin_Success | 🔲 |
| 5 | TestLogin_UserNotFound | 🔲 |
| 6 | TestLogin_WrongPassword | 🔲 |
| 7 | TestLogout_Success | 🔲 |
| 8 | TestRefreshAccessToken_Success | 🔲 |
| 9 | TestRefreshAccessToken_SessionNotFound | 🔲 |
| 10 | TestRefreshAccessToken_UserDeleted | 🔲 |
| 11 | TestUpsertUser_NewUser | 🔲 |
| 12 | TestUpsertUser_MergeLocalAccount | 🔲 |

### Middleware Tests

| # | Test | Trạng thái |
|---|------|------------|
| 13 | TestRequireAuth_ValidAccessToken | 🔲 |
| 14 | TestRequireAuth_NoTokens | 🔲 |
| 15 | TestRequireAuth_AutoRefresh | 🔲 |
| 16 | TestRequireAuth_APIRoute_401 | 🔲 |

### Port

| # | Approach | Trạng thái |
|---|----------|------------|
| 17 | Fiber tests pass | 🔲 |
| 18 | stdlib tests pass | 🔲 |
| 19 | Echo tests pass | 🔲 |
| 20 | `make test-all` pass | 🔲 |
| 21 | Coverage ≥ 80% service | 🔲 |

---

## Report

> Điền sau khi hoàn thành
