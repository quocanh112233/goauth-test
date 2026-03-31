# Prompt 10 — Echo: Config + Model + Repository + Service + Handler + Middleware

## Role

Bạn là một **Go Backend Engineer** thành thạo Echo framework v4. Bạn hiểu Echo cần custom `echo.Renderer`, handler return `error`, và cookie dùng `*http.Cookie`.

---

## Context

Phase 5. Echo — framework cuối cùng. Pattern rõ từ Gin/Fiber/stdlib.

Điểm khác biệt Echo:
- Handler: `func(c echo.Context) error` (return error)
- Template: implement `echo.Renderer` interface
- Cookie: `c.SetCookie(&http.Cookie{...})`, `c.Cookie("name")`
- Context value: `c.Set()`, `c.Get()`

**Business logic giống hệt Gin.**

---

## Dependencies (Prompt phụ thuộc)

| Prompt | Đầu ra cần thiết |
|--------|-----------------|
| 01 | `shared/templates/*.html`, `echo/` folder skeleton |
| 02 | MongoDB Atlas đã seed |
| 03–05 | Gin implementation hoàn chỉnh (reference) |

---

## Yêu cầu

### 1. echo/go.mod

- Dependencies:
  - `github.com/labstack/echo/v4 v4.12.0`
  - `go.mongodb.org/mongo-driver v1.15.0`
  - `golang.org/x/crypto v0.22.0`
  - `github.com/golang-jwt/jwt/v5 v5.2.1`
  - `github.com/joho/godotenv v1.5.1`

### 2. Config: `Framework` = `"Echo"`, `Port` = `"8084"`

### 3–5. db + model + repository

- Copy từ Gin (User + Session, 6+4 methods)

### 6. Service: copy từ Gin (5 methods)

### 7. echo/internal/renderer/renderer.go

```go
type TemplateRenderer struct {
    templates map[string]*template.Template
}

func NewRenderer(templateDir string) *TemplateRenderer {
    // Parse 5 pages: login, signup, home, dashboard, error
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    return tmpl.ExecuteTemplate(w, "base", data)
}
```

### 8. Handler (Echo API)

- Cookie:
  ```go
  c.SetCookie(&http.Cookie{
      Name: "access_token", Value: accessToken,
      MaxAge: 1800, Path: "/",
      HttpOnly: true, Secure: true, SameSite: http.SameSiteLaxMode,
  })
  c.SetCookie(&http.Cookie{
      Name: "refresh_token", Value: refreshToken,
      MaxAge: 604800, Path: "/",
      HttpOnly: true, Secure: true, SameSite: http.SameSiteLaxMode,
  })
  ```
- Redirect: `return c.Redirect(http.StatusSeeOther, "/home")` hoặc `/dashboard`
- Logout POST: xóa 2 cookies + session DB
- Handler **return error** — return `nil` sau redirect

### 9. home.go + dashboard.go + api.go

- Home: `c.Get("user").(*model.User)` → render `home`
- Dashboard: check admin → render `dashboard`
- API: `c.JSON(200, map[string]interface{}{"user": user})`

### 10. Middleware

Echo signature: `func(next echo.HandlerFunc) echo.HandlerFunc`
- Auto-refresh logic giống Gin
- `c.Set("user", user)` / `c.Get("user")`
- HTML redirect vs API 401

### 11. echo/.env.example

`PORT=8084`

---

## Anti-Patterns (KHÔNG được làm)

❌ Không quên implement `echo.Renderer`
❌ Không quên return error trong handler
❌ Không quên 2 cookies
❌ Không dùng GET cho logout

---

## Acceptance Criteria

1. `cd echo && go build ./...` pass
2. TemplateRenderer implement `echo.Renderer`
3. Service/Repository giống hệt Gin
4. Login set 2 cookies, redirect theo role

---

## Checklist hoàn thành

| # | Yêu cầu | Trạng thái | Ghi chú |
|---|---------|------------|---------|
| 1 | go.mod echo/v4 + dependencies | 🔲 | |
| 2 | Config.Framework = "Echo", Port = "8084" | 🔲 | |
| 3 | Model: User (Phone, json:"-") + Session | 🔲 | |
| 4 | UserRepo 6 + SessionRepo 4 methods | 🔲 | |
| 5 | AuthService 5 methods | 🔲 | |
| 6 | TemplateRenderer parse 5 pages | 🔲 | |
| 7 | Handler set 2 cookies | 🔲 | |
| 8 | Login redirect theo role | 🔲 | |
| 9 | Logout POST | 🔲 | |
| 10 | Dashboard check admin | 🔲 | |
| 11 | GetMe JSON (no password) | 🔲 | |
| 12 | Middleware auto-refresh | 🔲 | |
| 13 | Handler return error | 🔲 | |

---

## Report

> Điền sau khi hoàn thành
