# Prompt 08 — net/http (stdlib): Config + DB + Model + Repository + Service + Handler + Middleware

## Role

Bạn là một **Go Backend Engineer** thành thạo Go standard library. Bạn hiểu rõ `net/http` thuần và **Go 1.22+ enhanced ServeMux**. Bạn tự viết middleware chaining, template rendering, request handling chỉ với standard library.

---

## Context

Phase 4. Approach **zero-framework** — chỉ dùng Go standard library. Go 1.22+ ServeMux hỗ trợ method routing: `"GET /login"`, `"POST /login"`.

Điểm khác biệt:
- Handler: `func(w http.ResponseWriter, r *http.Request)`
- Router: `http.NewServeMux()` Go 1.22+
- Cookie: `http.SetCookie()`, `*http.Cookie`
- Context value: `context.WithValue()` + typed key
- Middleware: tự viết chaining pattern

**Business logic giống hệt Gin.**

---

## Dependencies (Prompt phụ thuộc)

| Prompt | Đầu ra cần thiết |
|--------|-----------------|
| 01 | `shared/templates/*.html`, `stdlib/` folder skeleton |
| 02 | MongoDB Atlas đã seed |
| 03–05 | Gin implementation hoàn chỉnh (reference) |

---

## Yêu cầu

### 1. stdlib/go.mod

- Module: `github.com/yourusername/go-auth-frameworks/stdlib`
- Go version: 1.22
- Dependencies (**không có framework**):
  - `go.mongodb.org/mongo-driver v1.15.0`
  - `golang.org/x/crypto v0.22.0`
  - `github.com/golang-jwt/jwt/v5 v5.2.1`
  - `github.com/joho/godotenv v1.5.1`

### 2. Config: `Framework` = `"net/http"`, `Port` = `"8083"`

### 3–5. db + model + repository

- Copy từ Gin (User + Session models, 6+4 repository methods)

### 6. Service: copy từ Gin (5 methods)

### 7. stdlib/internal/middleware/middleware.go

```go
type Middleware func(http.Handler) http.Handler

func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
    for i := len(middlewares) - 1; i >= 0; i-- {
        h = middlewares[i](h)
    }
    return h
}
```

### 8. stdlib/internal/middleware/auth.go

```go
type contextKey string
const UserKey contextKey = "user"

func RequireAuth(service service.AuthService, jwtSecret string) Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // 1. Đọc cookie access_token: r.Cookie("access_token")
            // 2. Parse JWT → nếu valid → lấy user → set context → continue
            // 3. Nếu expired → đọc refresh_token cookie
            // 4. Auto-refresh: service.RefreshAccessToken()
            // 5. Set access_token cookie mới: http.SetCookie(w, &http.Cookie{...})
            // 6. Set user vào context:
            //    ctx := context.WithValue(r.Context(), UserKey, user)
            //    next.ServeHTTP(w, r.WithContext(ctx))
            // 7. HTML redirect vs API 401 (check r.URL.Path prefix "/api/")
        })
    }
}
```

> **Typed context key** bắt buộc — KHÔNG dùng string literal.

### 9. stdlib/internal/handler/template.go

```go
type TemplateMap struct {
    templates map[string]*template.Template
}

func NewTemplateMap(templateDir string) *TemplateMap {
    // Parse: base.html + login.html → templates["login"]
    // Parse: base.html + signup.html → templates["signup"]
    // Parse: base.html + home.html → templates["home"]
    // Parse: base.html + dashboard.html → templates["dashboard"]
    // Parse: base.html + error.html → templates["error"]
}

func (tm *TemplateMap) Render(w http.ResponseWriter, name string, data interface{}) {
    tmpl.ExecuteTemplate(w, "base", data)
}
```

### 10. Handler (net/http)

- Form: `r.ParseForm()` + `r.FormValue("email")`
- Cookie:
  ```go
  http.SetCookie(w, &http.Cookie{
      Name: "access_token", Value: accessToken,
      MaxAge: 1800, Path: "/",
      HttpOnly: true, Secure: true, SameSite: http.SameSiteLaxMode,
  })
  http.SetCookie(w, &http.Cookie{
      Name: "refresh_token", Value: refreshToken,
      MaxAge: 604800, Path: "/",
      HttpOnly: true, Secure: true, SameSite: http.SameSiteLaxMode,
  })
  ```
- Redirect POST: `http.Redirect(w, r, "/home", http.StatusSeeOther)`
- Login redirect: admin→/dashboard, user→/home
- Logout POST: xóa 2 cookies (MaxAge=-1) + xóa session DB

### 11. home.go + dashboard.go + api.go

- Home: `r.Context().Value(middleware.UserKey).(*model.User)` → render home
- Dashboard: check admin → render dashboard
- API: `json.NewEncoder(w).Encode(map[string]interface{}{"user": user})`

### 12. stdlib/.env.example

`PORT=8083`

---

## Anti-Patterns (KHÔNG được làm)

❌ KHÔNG import third-party router
❌ Không dùng string literal làm context key
❌ Không parse template mỗi request
❌ Không quên 2 cookies

---

## Acceptance Criteria

1. `cd stdlib && go build ./...` pass
2. `go.mod` KHÔNG chứa web framework
3. Service/Repository giống hệt Gin
4. Middleware auto-refresh + typed context key

---

## Checklist hoàn thành

| # | Yêu cầu | Trạng thái | Ghi chú |
|---|---------|------------|---------|
| 1 | go.mod NO framework dependencies | 🔲 | |
| 2 | Config.Framework = "net/http", Port = "8083" | 🔲 | |
| 3 | Model: User (Phone, json:"-") + Session | 🔲 | |
| 4 | UserRepo 6 methods + SessionRepo 4 methods | 🔲 | |
| 5 | AuthService 5 methods | 🔲 | |
| 6 | Middleware Chain function | 🔲 | |
| 7 | Middleware RequireAuth typed context key | 🔲 | |
| 8 | Middleware auto-refresh | 🔲 | |
| 9 | TemplateMap parse 5 pages riêng | 🔲 | |
| 10 | Login set 2 cookies | 🔲 | |
| 11 | Login redirect theo role | 🔲 | |
| 12 | Logout POST: xóa 2 cookies + session | 🔲 | |
| 13 | Dashboard check admin role | 🔲 | |
| 14 | GetMe JSON (no password) | 🔲 | |

---

## Report

> Điền sau khi hoàn thành
