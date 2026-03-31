# Prompt 06 — Fiber: Config + Model + Repository + Service + Handler + Middleware

## Role

Bạn là một **Go Backend Engineer** thành thạo Fiber framework. Port Gin → Fiber.

---

## Context

Phase 3. Logic business giống hệt Gin, chỉ đổi transport layer. Quy tắc chung xem `docs/conventions.md`.

**Điểm khác biệt Fiber vs Gin:**
- Cookie: `c.Cookie(&fiber.Cookie{...})`
- Context: `c.Locals()` thay `c.Set()/c.Get()`
- Template: `gofiber/template/html/v2`

---

## Dependencies

| Prompt | Đầu ra cần thiết |
|--------|-----------------|
| 01 | `shared/templates/*.html` |
| 02 | MongoDB Atlas đã seed |
| 03–05 | Gin implementation (reference) |

---

## Yêu cầu

### 1. fiber/go.mod

- `github.com/gofiber/fiber/v2 v2.52.0`
- `github.com/gofiber/template/html/v2 v2.1.1`
- + mongodb, crypto, jwt, godotenv (cùng version Gin)

### 2. fiber/config/config.go

- `Framework = "Fiber"`, `Port = "8082"`
- **Có `TemplateDir` + `IsProduction`** (giống Gin, xem `conventions.md`)

### 3–6. db + model + repository + service

Copy từ Gin (không thay đổi logic).

### 7. fiber/internal/handler/

Cookie dùng `cfg.IsProduction`:
```go
c.Cookie(&fiber.Cookie{
    Name: "access_token", Value: accessToken,
    MaxAge: 1800, Path: "/",
    HTTPOnly: true,
    Secure:   cfg.IsProduction,   // ← KHÔNG hardcode true
    SameSite: "Lax",
})
c.Cookie(&fiber.Cookie{
    Name: "refresh_token", Value: refreshToken,
    MaxAge: 604800, Path: "/",
    HTTPOnly: true,
    Secure:   cfg.IsProduction,
    SameSite: "Lax",
})
```

### 8. Middleware

Auto-refresh, HTML/API phân biệt — giống Gin logic.

### 9. fiber/.env.example

```
MONGO_URI=...
MONGO_DB=goauth
JWT_SECRET=...
PORT=8082
# TEMPLATE_DIR=../shared/templates
# APP_ENV=development
```

---

## Anti-Patterns

❌ Không hardcode `Secure: true` — dùng `cfg.IsProduction`
❌ Không thay đổi business logic

---

## Acceptance Criteria

1. `cd fiber && go build ./...` pass
2. Cookie Secure = cfg.IsProduction
3. Service/Repository giống Gin

---

## Checklist hoàn thành

| # | Yêu cầu | Trạng thái | Ghi chú |
|---|---------|------------|---------|
| 1 | Config: Framework="Fiber", Port="8082" | 🔲 | |
| 2 | Config: TemplateDir + IsProduction | 🔲 | |
| 3 | Model/Repo/Service giống Gin | 🔲 | |
| 4 | Cookie Secure = cfg.IsProduction | 🔲 | |
| 5 | Middleware auto-refresh | 🔲 | |
| 6 | Login redirect theo role | 🔲 | |
| 7 | Logout POST | 🔲 | |

---

## Report

> Điền sau khi hoàn thành
