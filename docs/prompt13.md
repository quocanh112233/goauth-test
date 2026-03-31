# Prompt 13 — Google OAuth: Shared Setup + Gin Implementation

## Role

Bạn là một **Go Security Engineer** chuyên OAuth 2.0 / OpenID Connect. Bạn implement Google OAuth đúng chuẩn, tích hợp với hệ thống Access Token + Refresh Token + Sessions hiện có.

---

## Context

Phase 7. Email/password + session auth đã hoàn chỉnh. Thêm "Login with Google".

Flow:
```
Click "Login with Google"
→ GET /auth/google (tạo state cookie, redirect Google)
→ Google login
→ GET /auth/google/callback (validate state, exchange code, get user info)
→ Upsert user (google_id → email → tạo mới)
→ Tạo session (như login thường)
→ Set 2 cookies (access_token + refresh_token)
→ Redirect /home hoặc /dashboard (theo role)
```

---

## Dependencies (Prompt phụ thuộc)

| Prompt | Đầu ra cần thiết |
|--------|-----------------|
| 03 | UserRepo có FindByGoogleID, FindByEmail, UpdateByID |
| 03 | SessionRepo có Create |
| 04 | AuthService logic (JWT generation, session creation) |
| 05 | Gin app hoàn chỉnh đang chạy |

---

## Yêu cầu

### 1. docs/GOOGLE_OAUTH_SETUP.md

Hướng dẫn setup Google OAuth Console:
- Tạo project, enable API
- Tạo OAuth Client ID
- Redirect URIs:
  - `http://localhost:808{1,2,3,4}/auth/google/callback`
  - `https://goauth-{gin,fiber,stdlib,echo}.fly.dev/auth/google/callback`

### 2. Config cập nhật

Thêm: `GoogleClientID`, `GoogleClientSecret`, `GoogleRedirectURL` (optional — chỉ warning nếu thiếu)

### 3. gin/go.mod — thêm `golang.org/x/oauth2 v0.19.0`

### 4. gin/internal/service/oauth.go

```go
type OAuthService interface {
    GetAuthURL(state string) string
    ExchangeToken(ctx context.Context, code string) (*GoogleUser, error)
    UpsertUser(ctx context.Context, googleUser *GoogleUser) (*model.User, error)
}
```

- **UpsertUser** logic:
  1. Tìm `google_id` → update name/updated_at → return
  2. Tìm `email` → merge account (set google_id, provider="google") → return
  3. Không thấy → tạo mới (role="user", provider="google", password="")

> **Sau UpsertUser**: dùng AuthService logic tạo session + 2 tokens (không duplicate code)

### 5. gin/internal/handler/oauth.go

- `InitiateGoogleLogin` — GET /auth/google:
  - State: `crypto/rand` 16 bytes hex
  - Cookie `oauth_state`: MaxAge=300, HttpOnly, Secure, SameSite=Lax
  - Redirect → Google OAuth URL (307)

- `GoogleCallback` — GET /auth/google/callback:
  - Validate state cookie == query param state
  - Exchange code → get Google user info
  - Upsert user
  - **Tạo session + 2 cookies** (giống POST /login flow):
    - Access token JWT (30 phút)
    - Refresh token (random, 7 ngày) + session DB
  - Xóa cookie `oauth_state`
  - Redirect theo role: admin→/dashboard, user→/home (303)

### 6. Cập nhật login.html

Thêm nút Google OAuth (thay TODO comment hiện tại):
```html
<a href="/auth/google" class="google-btn">Đăng nhập với Google</a>
```

### 7. Router — thêm 2 routes (public)

```go
r.GET("/auth/google", oauthHandler.InitiateGoogleLogin)
r.GET("/auth/google/callback", oauthHandler.GoogleCallback)
```

---

## Anti-Patterns (KHÔNG được làm)

❌ Không hardcode state
❌ Không bỏ qua state validation
❌ Không tạo duplicate user khi merge
❌ Không quên tạo session + 2 cookies trong OAuth callback (giống login flow)

---

## Acceptance Criteria

1. `cd gin && go build ./...` pass
2. Google OAuth flow hoạt động end-to-end
3. Login Google → set 2 cookies + tạo session
4. Login Google lần 2 → không duplicate
5. Local user login Google cùng email → merge

---

## Checklist hoàn thành

| # | Yêu cầu | Trạng thái | Ghi chú |
|---|---------|------------|---------|
| 1 | Config thêm 3 Google fields | 🔲 | |
| 2 | go.mod thêm oauth2 | 🔲 | |
| 3 | OAuthService 3 methods | 🔲 | |
| 4 | UpsertUser: google_id → email → tạo mới | 🔲 | |
| 5 | Callback tạo session + set 2 cookies | 🔲 | |
| 6 | Callback redirect theo role | 🔲 | |
| 7 | State crypto/rand 16 bytes | 🔲 | |
| 8 | State cookie HttpOnly+Secure MaxAge=300 | 🔲 | |
| 9 | Validate state trong callback | 🔲 | |
| 10 | Xóa oauth_state cookie sau callback | 🔲 | |
| 11 | login.html có nút Google | 🔲 | |
| 12 | 2 routes OAuth (public) | 🔲 | |
| 13 | GOOGLE_OAUTH_SETUP.md | 🔲 | |

---

## Report

> Điền sau khi hoàn thành
