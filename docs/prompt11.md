# Prompt 11 — Echo: Config + Model + Repository + Service + Handler + Middleware

## Role

Bạn là một **Go Backend Engineer** thành thạo Echo v4. Port Gin → Echo.

---

## Context

Phase 5. Echo — framework cuối. Conventions xem `docs/conventions.md`.

**Điểm khác biệt Echo:**
- Handler: `func(c echo.Context) error` (return error)
- Template: implement `echo.Renderer` interface
- Cookie: `c.SetCookie(&http.Cookie{...})`
- Context: `c.Set()`, `c.Get()`

---

## Dependencies

| Prompt | Đầu ra cần thiết |
|--------|-----------------|
| 01 | `shared/templates/*.html` |
| 02 | MongoDB Atlas đã seed |
| 03–05 | Gin implementation (reference) |

---

## Yêu cầu

### 1. echo/go.mod

- `github.com/labstack/echo/v4 v4.12.0`
- + mongodb, crypto, jwt, godotenv

### 2. Config: `Framework = "Echo"`, `Port = "8084"`, **TemplateDir + IsProduction**

### 3–6. db + model + repository + service — copy từ Gin

### 7. echo/internal/renderer/renderer.go

```go
type TemplateRenderer struct {
    templates map[string]*template.Template
}

func NewRenderer(templateDir string) *TemplateRenderer {
    r := &TemplateRenderer{templates: make(map[string]*template.Template)}
    pages := []string{"login", "signup", "home", "dashboard", "error"}
    for _, page := range pages {
        r.templates[page] = template.Must(template.ParseFiles(
            filepath.Join(templateDir, "base.html"),
            filepath.Join(templateDir, page+".html"),
        ))
    }
    return r
}

// Implement echo.Renderer interface
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    return t.templates[name].ExecuteTemplate(w, "base", data)
}
```

> Dùng `cfg.TemplateDir` — KHÔNG hardcode path.

### 8. Handler

Cookie dùng `cfg.IsProduction`:
```go
c.SetCookie(&http.Cookie{
    Name: "access_token", Value: accessToken,
    MaxAge: 1800, Path: "/",
    HttpOnly: true, Secure: cfg.IsProduction, SameSite: http.SameSiteLaxMode,
})
```

### 9. Middleware

Echo signature: `func(next echo.HandlerFunc) echo.HandlerFunc`
- Auto-refresh, HTML/API phân biệt — giống Gin logic

---

## Anti-Patterns

❌ Không quên implement `echo.Renderer`
❌ Không quên return error trong handler
❌ Không hardcode `Secure: true`

---

## Acceptance Criteria

1. `cd echo && go build ./...` pass
2. TemplateRenderer implement `echo.Renderer`
3. Cookie Secure = cfg.IsProduction

---

## Checklist hoàn thành

| # | Yêu cầu | Trạng thái | Ghi chú |
|---|---------|------------|---------|
| 1 | Config: Framework="Echo", Port="8084" | 🔲 | |
| 2 | Config: TemplateDir + IsProduction | 🔲 | |
| 3 | Model/Repo/Service giống Gin | 🔲 | |
| 4 | TemplateRenderer (cfg.TemplateDir, 5 pages) | 🔲 | |
| 5 | Cookie Secure = cfg.IsProduction | 🔲 | |
| 6 | Middleware auto-refresh | 🔲 | |
| 7 | Handler return error | 🔲 | |

---

## Report

> Điền sau khi hoàn thành
