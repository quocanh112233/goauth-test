# Prompt 14 — Google OAuth: Shared Setup + Gin Implementation

## Role

Bạn là một **Go Security Engineer** chuyên OAuth 2.0. Cookie conventions xem `docs/conventions.md`.

---

## Context

Phase 7. Email/password auth đã hoàn chỉnh. Thêm "Login with Google".

Flow: Click → `/auth/google` → Google → `/auth/google/callback` → Upsert user → Tạo session + 2 cookies → Redirect theo role.

---

## Dependencies

| Prompt | Đầu ra cần thiết |
|--------|-----------------|
| 03 | UserRepo (FindByGoogleID, FindByEmail, UpdateByID) |
| 04 | AuthService (JWT generation, session creation) |
| 05 | Gin app hoàn chỉnh |

---

## Yêu cầu

### 1. docs/GOOGLE_OAUTH_SETUP.md

### 2. Config — thêm GoogleClientID, GoogleClientSecret, GoogleRedirectURL

### 3. gin/go.mod — thêm `golang.org/x/oauth2 v0.19.0`

### 4. gin/internal/service/oauth.go

`OAuthService`: `GetAuthURL`, `ExchangeToken`, `UpsertUser`

- UpsertUser: google_id → email → tạo mới

### 5. gin/internal/handler/oauth.go

- **InitiateGoogleLogin** (GET /auth/google):
  - State: `crypto/rand` 16 bytes hex
  - Cookie `oauth_state`: MaxAge=300, HttpOnly=true, **Secure=cfg.IsProduction**, SameSite=Lax
  - Redirect → Google (307)

- **GoogleCallback** (GET /auth/google/callback):
  - Validate state
  - Upsert user
  - **Tạo session + set 2 cookies** (giống POST /login):
    ```go
    c.SetCookie("access_token", accessToken, 1800, "/", "", cfg.IsProduction, true)
    c.SetCookie("refresh_token", refreshToken, 604800, "/", "", cfg.IsProduction, true)
    ```
  - Xóa oauth_state cookie
  - Redirect theo role (303)

### 6. Router — thêm 2 routes public

---

## Anti-Patterns

❌ Không hardcode state
❌ Không hardcode `Secure: true` — dùng `cfg.IsProduction`
❌ Không quên tạo session + 2 cookies trong callback

---

## Acceptance Criteria

1. Google OAuth flow end-to-end
2. Cookie Secure = cfg.IsProduction
3. Callback tạo session + 2 cookies

---

## Checklist hoàn thành

| # | Yêu cầu | Trạng thái | Ghi chú |
|---|---------|------------|---------|
| 1 | Config 3 Google fields | 🔲 | |
| 2 | OAuthService 3 methods | 🔲 | |
| 3 | UpsertUser logic | 🔲 | |
| 4 | State crypto/rand | 🔲 | |
| 5 | oauth_state cookie Secure=cfg.IsProduction | 🔲 | |
| 6 | Callback: session + 2 cookies | 🔲 | |
| 7 | Redirect theo role | 🔲 | |
| 8 | 2 routes public | 🔲 | |
| 9 | GOOGLE_OAUTH_SETUP.md | 🔲 | |

---

## Report

> Điền sau khi hoàn thành
