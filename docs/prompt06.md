# Prompt 06 — Fiber: Config + Model + Repository + Service + Handler + Middleware

## Role

Bạn là một **Go Backend Engineer** thành thạo Fiber framework. Bạn port Gin implementation sang Fiber, giữ nguyên logic business, chỉ đổi transport layer.

---

## Context

Phase 3. Pattern đã rõ từ Gin (Prompt 03–05). Prompt này implement **toàn bộ** Fiber layers.

Điểm khác biệt Fiber vs Gin:
- `*fiber.Ctx` thay `*gin.Context`
- `c.Locals()` thay `c.Set()/c.Get()`
- Cookie: `c.Cookie(&fiber.Cookie{...})`
- Template: `gofiber/template/html/v2`

**Business logic (service, repository, model) giống hệt Gin.**

---

## Dependencies (Prompt phụ thuộc)

| Prompt | Đầu ra cần thiết |
|--------|-----------------|
| 01 | `shared/templates/*.html`, `fiber/` folder skeleton |
| 02 | MongoDB Atlas đã seed (users + sessions) |
| 03–05 | Gin implementation hoàn chỉnh (reference) |

---

## Yêu cầu

### 1. fiber/go.mod

- Module: `github.com/yourusername/go-auth-frameworks/fiber`
- Go version: 1.22
- Dependencies:
  - `github.com/gofiber/fiber/v2 v2.52.0`
  - `github.com/gofiber/template/html/v2 v2.1.1`
  - `go.mongodb.org/mongo-driver v1.15.0`
  - `golang.org/x/crypto v0.22.0`
  - `github.com/golang-jwt/jwt/v5 v5.2.1`
  - `github.com/joho/godotenv v1.5.1`

### 2. fiber/config/config.go

- `Framework` hardcode = `"Fiber"`, `Port` default = `"8082"`

### 3–5. db/mongo.go + model/ + repository/

- **Copy từ Gin** — giữ nguyên (chỉ đổi package name/import path)
- Model: User (với Phone, json:"-") + Session
- Repository: UserRepository (6 methods) + SessionRepository (4 methods)

### 6. fiber/internal/service/auth.go

- **Copy từ Gin** — logic giống hệt (5 methods)

### 7. fiber/internal/handler/auth.go

Fiber API:

- Form: `c.FormValue("email")`, `c.FormValue("password")`
- Cookie:
  ```go
  // Access token — 30 phút
  c.Cookie(&fiber.Cookie{
      Name: "access_token", Value: accessToken,
      MaxAge: 1800, Path: "/",
      HTTPOnly: true, Secure: true, SameSite: "Lax",
  })
  // Refresh token — 7 ngày
  c.Cookie(&fiber.Cookie{
      Name: "refresh_token", Value: refreshToken,
      MaxAge: 604800, Path: "/",
      HTTPOnly: true, Secure: true, SameSite: "Lax",
  })
  ```
- Redirect theo role: `c.Redirect("/home", 303)` hoặc `c.Redirect("/dashboard", 303)`
- Logout (POST): xóa 2 cookies (MaxAge: -1) + xóa session DB

### 8. fiber/internal/handler/home.go

- `c.Locals("user").(*model.User)` → render `home` template

### 9. fiber/internal/handler/dashboard.go

- Check `user.Role != "admin"` → redirect `/home`
- Render `dashboard` template

### 10. fiber/internal/handler/api.go

- `GET /api/me` → `c.JSON(200, fiber.Map{"user": user})`

### 11. fiber/internal/middleware/auth.go

Giống logic Gin nhưng dùng Fiber API:
- `c.Cookies("access_token")`, `c.Cookies("refresh_token")`
- `c.Locals("user", user)` thay `c.Set("user", user)`
- Auto-refresh logic giống hệt Gin
- HTML redirect vs API 401: check `strings.HasPrefix(c.Path(), "/api/")`

### 12. fiber/.env.example

```
MONGO_URI=mongodb+srv://<user>:<pass>@cluster.mongodb.net/?retryWrites=true&w=majority
MONGO_DB=goauth
JWT_SECRET=your-fiber-jwt-secret-key-change-me
PORT=8082
```

---

## Anti-Patterns (KHÔNG được làm)

❌ Không thay đổi logic business — chỉ đổi transport layer
❌ Không quên set **2 cookies** (access + refresh)
❌ Không dùng GET cho logout
❌ Không dùng 302 cho POST redirect — dùng 303

---

## Acceptance Criteria

1. `cd fiber && go build ./...` pass
2. Service/Repository logic giống hệt Gin
3. Login set 2 cookies, redirect theo role
4. Middleware auto-refresh hoạt động

---

## Checklist hoàn thành

| # | Yêu cầu | Trạng thái | Ghi chú |
|---|---------|------------|---------|
| 1 | go.mod đúng module + dependencies | 🔲 | |
| 2 | Config.Framework = "Fiber", Port = "8082" | 🔲 | |
| 3 | Model: User (có Phone, json:"-") + Session | 🔲 | |
| 4 | UserRepository 6 methods + SessionRepository 4 methods | 🔲 | |
| 5 | AuthService 5 methods (giống Gin) | 🔲 | |
| 6 | Handler set 2 cookies (access + refresh) | 🔲 | |
| 7 | Login redirect theo role (303) | 🔲 | |
| 8 | Logout POST: xóa 2 cookies + session DB | 🔲 | |
| 9 | ShowHome render home.html | 🔲 | |
| 10 | ShowDashboard check admin role | 🔲 | |
| 11 | GetMe trả JSON (no password) | 🔲 | |
| 12 | Middleware auto-refresh | 🔲 | |
| 13 | Middleware HTML redirect vs API 401 | 🔲 | |

---

## Report

> Điền sau khi hoàn thành
