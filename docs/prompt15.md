# Prompt 15 — Unit Tests: Tất cả frameworks

## Role

Bạn là một **Go Test Engineer** chuyên viết unit test, mock dependencies, và đảm bảo behavior consistency giữa 4 implementations.

---

## Context

Phase 8 (cuối cùng). Viết tests để verify logic business và auth flow mới (Access Token + Refresh Token + Sessions + Auto-refresh).

Strategy: Tests cho **Gin trước**, copy + adapt cho 3 approach còn lại.

---

## Dependencies (Prompt phụ thuộc)

| Prompt | Đầu ra cần thiết |
|--------|-----------------|
| 05, 07, 09, 11 | 4 apps hoàn chỉnh |
| 13–14 | Google OAuth tất cả |

---

## Yêu cầu

### 1. Test Dependencies

- `github.com/stretchr/testify v1.9.0`

### 2. Mock Repositories

`internal/repository/mock_user.go` — mock 6 methods `UserRepository`
`internal/repository/mock_session.go` — mock 4 methods `SessionRepository`

### 3. Service Tests — auth_test.go

**Signup tests:**
```go
TestSignup_Success
TestSignup_DuplicateEmail
TestSignup_DuplicatePhone
TestSignup_BcryptCost12
```

**Login tests:**
```go
TestLogin_Success
// Assert return 3 values: accessToken, refreshToken, role
// Assert session tạo trong DB (sessionRepo.Create called)

TestLogin_UserNotFound
// Assert error = "email hoặc mật khẩu không đúng"

TestLogin_WrongPassword
// Assert CÙNG error message (chống user enumeration)

TestLogin_AccessTokenExpiry30Min
// Parse JWT → assert exp = iat + 30 phút

TestLogin_RefreshTokenRandom
// Assert refreshToken != "" && len >= 64 (32 bytes hex)
```

**Logout tests:**
```go
TestLogout_Success
// Assert sessionRepo.DeleteByRefreshToken called

TestLogout_InvalidRefreshToken
// Assert không panic, graceful handling
```

**RefreshAccessToken tests:**
```go
TestRefreshAccessToken_Success
// Mock: FindByRefreshToken → valid session
// Assert: return new access token
// Assert: new access token has user's claims

TestRefreshAccessToken_SessionNotFound
// Assert error

TestRefreshAccessToken_SessionExpired
// Mock: session with ExpiredAt < now
// Assert error

TestRefreshAccessToken_UserDeleted
// Mock: valid session but user not found
// Assert error + session deleted
```

### 4. OAuth Service Tests — oauth_test.go

```go
TestUpsertUser_NewUser
// Assert Create called, provider="google"

TestUpsertUser_ExistingGoogleUser
// Assert UpdateByID called (update name)

TestUpsertUser_MergeLocalAccount
// Assert UpdateByID called (set google_id, provider="google")
```

### 5. Middleware Tests — auth_test.go

```go
TestRequireAuth_ValidAccessToken
// Assert next handler called + user in context

TestRequireAuth_NoTokens
// Assert redirect /login (HTML) hoặc 401 (API)

TestRequireAuth_ExpiredAccessToken_ValidRefresh
// Assert auto-refresh: new access_token cookie set
// Assert next handler called

TestRequireAuth_ExpiredAccessToken_InvalidRefresh
// Assert both cookies cleared + redirect /login

TestRequireAuth_APIRoute_Unauthorized
// Request to /api/me without token
// Assert 401 JSON (not redirect)
```

### 6. Port tests sang Fiber, stdlib, Echo

- Service tests: copy nguyên (logic giống hệt)
- Middleware tests: adapt framework API
  - Fiber: `app.Test()`
  - stdlib: `httptest.NewRecorder()`
  - Echo: `echo.New()` + `httptest`

### 7. Makefile — targets đã có ở Prompt 12

`test-gin`, `test-fiber`, `test-stdlib`, `test-echo`, `test-all`, `test-coverage`

---

## Anti-Patterns (KHÔNG được làm)

❌ Không kết nối MongoDB thật — chỉ mock
❌ Không gộp nhiều test cases vào 1 function
❌ Không skip error/edge cases — chúng quan trọng hơn happy path

---

## Acceptance Criteria

1. `go test ./... -v` pass cho cả 4 approaches
2. Coverage cho service layer ≥ 80%
3. `make test-all` pass

---

## Checklist hoàn thành

### Service Tests (Gin reference)

| # | Test | Trạng thái |
|---|------|------------|
| 1 | TestSignup_Success | 🔲 |
| 2 | TestSignup_DuplicateEmail | 🔲 |
| 3 | TestSignup_DuplicatePhone | 🔲 |
| 4 | TestSignup_BcryptCost12 | 🔲 |
| 5 | TestLogin_Success (3 return values + session created) | 🔲 |
| 6 | TestLogin_UserNotFound | 🔲 |
| 7 | TestLogin_WrongPassword (same error) | 🔲 |
| 8 | TestLogin_AccessTokenExpiry30Min | 🔲 |
| 9 | TestLogin_RefreshTokenRandom | 🔲 |
| 10 | TestLogout_Success | 🔲 |
| 11 | TestRefreshAccessToken_Success | 🔲 |
| 12 | TestRefreshAccessToken_SessionNotFound | 🔲 |
| 13 | TestRefreshAccessToken_UserDeleted | 🔲 |
| 14 | TestUpsertUser_NewUser | 🔲 |
| 15 | TestUpsertUser_ExistingGoogleUser | 🔲 |
| 16 | TestUpsertUser_MergeLocalAccount | 🔲 |

### Middleware Tests (Gin reference)

| # | Test | Trạng thái |
|---|------|------------|
| 17 | TestRequireAuth_ValidAccessToken | 🔲 |
| 18 | TestRequireAuth_NoTokens | 🔲 |
| 19 | TestRequireAuth_ExpiredAccessToken_ValidRefresh | 🔲 |
| 20 | TestRequireAuth_ExpiredAccessToken_InvalidRefresh | 🔲 |
| 21 | TestRequireAuth_APIRoute_Unauthorized | 🔲 |

### Port

| # | Approach | Trạng thái |
|---|----------|------------|
| 22 | Fiber tests pass | 🔲 |
| 23 | stdlib tests pass | 🔲 |
| 24 | Echo tests pass | 🔲 |
| 25 | `make test-all` pass | 🔲 |
| 26 | Coverage ≥ 80% service | 🔲 |

---

## Report

> Điền sau khi hoàn thành
