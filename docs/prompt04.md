# Prompt 04 — Gin: Service + Handler + Middleware

## Role

Bạn là một **Go Backend Engineer** chuyên về authentication với JWT + session management.

---

## Context

Tiếp Prompt 03. Implement business logic + transport layer cho Gin. Auth flow chi tiết xem `docs/api-spec.md`, quy tắc chung xem `docs/conventions.md`.

---

## Dependencies

| Prompt | Đầu ra cần thiết |
|--------|-----------------|
| 01 | `shared/templates/*.html` (6 files) |
| 03 | Config (có TemplateDir + IsProduction), DB, Models, Repositories |

---

## Yêu cầu

### 1. gin/internal/service/auth.go

5 methods: `Signup`, `Login`, `Logout`, `RefreshAccessToken`, `GetUserByID`

Chi tiết xem prompt trước + `docs/api-spec.md`. Lưu ý:
- **Login** return 3 values: `accessToken, refreshToken, role`
- **Login** tạo session trong DB
- **Refresh token**: `crypto/rand` 32 bytes → hex (64 chars)
- **Error messages**: xem `docs/conventions.md` mục 4

### 2. gin/internal/handler/auth.go

- **ShowLogin/ShowSignup**: nếu đã authenticated → redirect `/home`
- **Login POST**: set 2 cookies (xem conventions.md mục 1), redirect theo role (conventions.md mục 9)
- **Signup POST**: validate confirm_password, redirect `/login` (303)
- **Logout POST**: xóa session DB + xóa 2 cookies, redirect `/login` (303)

**Cookie sử dụng `cfg.IsProduction`**:
```go
c.SetCookie("access_token", accessToken, 1800, "/", "", cfg.IsProduction, true)
//                                                       ^^^^^^^^^^^^^^^^
//                                                       Secure flag
c.SetCookie("refresh_token", refreshToken, 604800, "/", "", cfg.IsProduction, true)
```

> Xem `docs/conventions.md` mục 1 — **Không** hardcode `Secure: true` vì local dev dùng HTTP.

### 3. gin/internal/handler/home.go + dashboard.go + api.go

- **ShowHome** (`GET /home`): render home.html với User + Framework
- **ShowDashboard** (`GET /dashboard`): check admin → render; non-admin → redirect `/home`
- **GetMe** (`GET /api/me`): JSON response

### 4. gin/internal/middleware/auth.go

Auto-refresh middleware:
1. Check `access_token` cookie → parse JWT → valid → set user context → continue
2. Expired/missing → check `refresh_token` cookie
3. Có refresh → `service.RefreshAccessToken()` → set new cookie → continue
4. Không có refresh → redirect/401

```go
// Cookie mới khi auto-refresh
c.SetCookie("access_token", newToken, 1800, "/", "", cfg.IsProduction, true)
//                                                    ^^^^^^^^^^^^^^^^
```

**HTML vs API**: `strings.HasPrefix(path, "/api/")` → 401 JSON; otherwise → redirect `/login`

---

## Anti-Patterns

❌ Không hardcode `Secure: true` — dùng `cfg.IsProduction`
❌ Không chỉ set 1 cookie — phải cả access + refresh
❌ Không dùng GET cho logout
❌ Không dùng 302 cho POST redirect — dùng 303

---

## Acceptance Criteria

1. `cd gin && go build ./...` pass
2. Cookie Secure = `cfg.IsProduction` (không hardcode)
3. Login set 2 cookies, redirect theo role
4. Middleware auto-refresh, phân biệt HTML/API

---

## Checklist hoàn thành

| # | Yêu cầu | Trạng thái | Ghi chú |
|---|---------|------------|---------|
| 1 | AuthService 5 methods | 🔲 | |
| 2 | Signup check duplicate email + phone | 🔲 | |
| 3 | Login error chống user enumeration | 🔲 | |
| 4 | Login tạo session DB | 🔲 | |
| 5 | Cookie Secure = cfg.IsProduction | 🔲 | |
| 6 | Login set 2 cookies | 🔲 | |
| 7 | Login redirect theo role (303) | 🔲 | |
| 8 | Logout POST: xóa session + 2 cookies | 🔲 | |
| 9 | ShowHome render user info | 🔲 | |
| 10 | ShowDashboard check admin | 🔲 | |
| 11 | GetMe JSON (no password) | 🔲 | |
| 12 | Middleware auto-refresh | 🔲 | |
| 13 | Middleware HTML redirect vs API 401 | 🔲 | |

---

## Report

> Điền sau khi hoàn thành
