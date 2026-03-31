# Prompt 09 — net/http (stdlib): Service + Handler + Middleware + Template

## Role

Bạn là một **Go Backend Engineer** chuyên Go standard library. Bạn tự viết template renderer, middleware chaining, request handling — KHÔNG dùng framework.

---

## Context

Tiếp Prompt 08 (data layer). Prompt này implement **tầng logic** — phần phức tạp nhất của stdlib vì phải tự viết mọi thứ.

Quy tắc chung xem `docs/conventions.md`.

---

## Dependencies

| Prompt | Đầu ra cần thiết |
|--------|-----------------|
| 01 | `shared/templates/*.html` |
| 08 | Config, DB, Model, Repository |

---

## Yêu cầu

### 1. stdlib/internal/service/auth.go

5 methods — copy logic từ Gin (Prompt 04).

### 2. stdlib/internal/middleware/middleware.go

```go
type Middleware func(http.Handler) http.Handler

func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
    for i := len(middlewares) - 1; i >= 0; i-- {
        h = middlewares[i](h)
    }
    return h
}
```

### 3. stdlib/internal/middleware/auth.go

```go
// Typed context key — bắt buộc, KHÔNG dùng string literal
type contextKey string
const UserKey contextKey = "user"

func RequireAuth(service service.AuthService, jwtSecret string, isProduction bool) Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // 1. Check access_token cookie
            // 2. Parse JWT → valid → set context
            // 3. Expired → check refresh_token → auto-refresh
            // 4. Set new cookie khi auto-refresh:
            http.SetCookie(w, &http.Cookie{
                Name: "access_token", Value: newToken,
                MaxAge: 1800, Path: "/",
                HttpOnly: true,
                Secure:   isProduction,  // ← cfg.IsProduction
                SameSite: http.SameSiteLaxMode,
            })
            // 5. Set user context:
            ctx := context.WithValue(r.Context(), UserKey, user)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

### 4. stdlib/internal/handler/template.go

Template renderer tự viết:

```go
type TemplateMap struct {
    templates map[string]*template.Template
}

func NewTemplateMap(templateDir string) *TemplateMap {
    tm := &TemplateMap{templates: make(map[string]*template.Template)}
    pages := []string{"login", "signup", "home", "dashboard", "error"}
    for _, page := range pages {
        tm.templates[page] = template.Must(template.ParseFiles(
            filepath.Join(templateDir, "base.html"),
            filepath.Join(templateDir, page+".html"),
        ))
    }
    return tm
}

func (tm *TemplateMap) Render(w http.ResponseWriter, name string, data interface{}) {
    tmpl, ok := tm.templates[name]
    if !ok {
        http.Error(w, "template not found: "+name, 500)
        return
    }
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    tmpl.ExecuteTemplate(w, "base", data)
}
```

> Dùng `cfg.TemplateDir` — KHÔNG hardcode path.

### 5. stdlib/internal/handler/auth.go

net/http handler API:
```go
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    email := r.FormValue("email")
    password := r.FormValue("password")

    accessToken, refreshToken, role, err := h.service.Login(r.Context(), email, password)
    // ...

    // Set 2 cookies with cfg.IsProduction
    http.SetCookie(w, &http.Cookie{
        Name: "access_token", Value: accessToken,
        MaxAge: 1800, Path: "/",
        HttpOnly: true, Secure: h.cfg.IsProduction, SameSite: http.SameSiteLaxMode,
    })
    http.SetCookie(w, &http.Cookie{
        Name: "refresh_token", Value: refreshToken,
        MaxAge: 604800, Path: "/",
        HttpOnly: true, Secure: h.cfg.IsProduction, SameSite: http.SameSiteLaxMode,
    })

    // Redirect theo role
    if role == "admin" {
        http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
    } else {
        http.Redirect(w, r, "/home", http.StatusSeeOther)
    }
}
```

### 6. home.go + dashboard.go + api.go

- Home: `r.Context().Value(middleware.UserKey).(*model.User)` → render home
- Dashboard: check admin → render / redirect `/home`
- API: `json.NewEncoder(w).Encode(...)`

---

## Anti-Patterns

❌ KHÔNG import third-party web framework
❌ Không dùng string literal làm context key — dùng typed key
❌ Không hardcode template path — dùng `cfg.TemplateDir`
❌ Không hardcode `Secure: true` — dùng `cfg.IsProduction`
❌ Không parse templates mỗi request — parse 1 lần trong `NewTemplateMap()`

---

## Acceptance Criteria

1. `cd stdlib && go build ./...` pass
2. `go.mod` KHÔNG chứa web framework
3. Middleware Chain + typed context key
4. TemplateMap parse 5 pages đúng
5. Cookie Secure = cfg.IsProduction

---

## Checklist hoàn thành

| # | Yêu cầu | Trạng thái | Ghi chú |
|---|---------|------------|---------|
| 1 | AuthService 5 methods (giống Gin) | 🔲 | |
| 2 | Middleware Chain function | 🔲 | |
| 3 | Middleware RequireAuth typed context key | 🔲 | |
| 4 | Middleware auto-refresh + HTML/API phân biệt | 🔲 | |
| 5 | TemplateMap parse 5 pages (dùng cfg.TemplateDir) | 🔲 | |
| 6 | Login set 2 cookies (Secure=cfg.IsProduction) | 🔲 | |
| 7 | Login redirect theo role (303) | 🔲 | |
| 8 | Logout POST: xóa 2 cookies + session | 🔲 | |
| 9 | Dashboard check admin | 🔲 | |
| 10 | GetMe JSON (no password) | 🔲 | |

---

## Report

> Điền sau khi hoàn thành
