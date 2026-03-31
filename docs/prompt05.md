# Prompt 05 — Gin: Router + Dockerfile + fly.toml

## Role

Bạn là một **Go DevOps Engineer** kiêm Backend Developer.

---

## Context

Prompt 03 + 04 đã có đủ layers. Prompt này kết nối tất cả: Router, Template Loading, Health Check, Dockerfile, fly.toml.

Route mapping xem `docs/api-spec.md`, cookie/env conventions xem `docs/conventions.md`.

---

## Dependencies

| Prompt | Đầu ra cần thiết |
|--------|-----------------|
| 01 | `shared/templates/*.html` (6 files) |
| 03 | Config (có TemplateDir), DB, Model, Repository |
| 04 | Service, Handler, Middleware |

---

## Yêu cầu

### 1. gin/internal/router/router.go — Template Loading

Gin dùng `r.SetHTMLTemplate()` với custom `html/template` (KHÔNG dùng `gin.LoadHTMLGlob`):

```go
import "html/template"

func Setup(cfg *config.Config, database *mongo.Database) *gin.Engine {
    r := gin.Default()

    // ✅ ĐÚNG: parse từng cặp, dùng template.Must
    tmpl := template.New("")

    // Parse từng page riêng với base.html
    pages := []string{"login", "signup", "home", "dashboard", "error"}
    for _, page := range pages {
        // Clone base template, rồi parse page template vào
        t := template.Must(template.ParseFiles(
            filepath.Join(cfg.TemplateDir, "base.html"),
            filepath.Join(cfg.TemplateDir, page+".html"),
        ))
        // Add vào template set với tên = page name
        tmpl = template.Must(tmpl.AddParseTree(page, t.Tree))
    }

    // HOẶC: dùng cách đơn giản hơn — map[string]*template.Template
    templates := make(map[string]*template.Template)
    for _, page := range pages {
        templates[page] = template.Must(template.ParseFiles(
            filepath.Join(cfg.TemplateDir, "base.html"),
            filepath.Join(cfg.TemplateDir, page+".html"),
        ))
    }

    // Gin cần: r.SetHTMLTemplate(tmpl)
    // Hoặc implement gin.HTMLRender interface cho map approach
    // ...
}
```

> ⚠️ **KHÔNG dùng**: `r.LoadHTMLGlob("*.html")` — blocks `{{define "content"}}` sẽ bị override bởi file cuối cùng.

> **Template path**: dùng `cfg.TemplateDir` (không hardcode) — xem `docs/conventions.md` mục 7.

### 2. Route Registration

```
GET  /                → redirect /login (302)
GET  /health          → healthHandler.HealthCheck
GET  /login           → authHandler.ShowLogin
POST /login           → authHandler.Login
GET  /signup          → authHandler.ShowSignup
POST /signup          → authHandler.Signup

// Protected (RequireAuth middleware)
GET  /home            → homeHandler.ShowHome
GET  /dashboard       → dashboardHandler.ShowDashboard
POST /logout          → authHandler.Logout
GET  /api/me          → apiHandler.GetMe
```

### 3. Health Check

`GET /health` → `{"status":"ok","framework":"Gin","db":"connected"}`
- Ping MongoDB, fail → `{"db":"disconnected"}` status 503

### 4. gin/cmd/main.go — Graceful shutdown

### 5. gin/Dockerfile

```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /build
COPY gin/ ./gin/
COPY shared/ ./shared/
WORKDIR /build/gin
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /build/shared/templates/ ./shared/templates/
EXPOSE 8080
CMD ["./server"]
```

> Binary chạy từ `/app`, templates ở `/app/shared/templates/`. Trong fly.toml set `TEMPLATE_DIR=./shared/templates`.

### 6. gin/fly.toml

```toml
app = "goauth-gin"
primary_region = "sin"

[env]
  PORT = "8080"
  MONGO_DB = "goauth"
  APP_ENV = "production"
  TEMPLATE_DIR = "./shared/templates"

[http_service]
  internal_port = 8080
  force_https = true
  auto_stop_machines = true
  auto_start_machines = true

[[http_service.checks]]
  grace_period = "10s"
  interval = "30s"
  method = "GET"
  path = "/health"
  timeout = "5s"
```

> Secrets: `fly secrets set MONGO_URI="..." JWT_SECRET="..." -a goauth-gin`

---

## Anti-Patterns

❌ **KHÔNG dùng `gin.LoadHTMLGlob()`** — blocks bị override
❌ Không hardcode template path — dùng `cfg.TemplateDir`
❌ Không đặt secrets trong fly.toml
❌ Không quên `APP_ENV=production` + `TEMPLATE_DIR` trong fly.toml

---

## Acceptance Criteria

1. `cd gin && go build ./...` pass
2. Template render đúng (login, signup, home, dashboard, error)
3. E2E: Signup → Login → Home → Logout
4. Admin Login → Dashboard
5. `/api/me` → JSON (no password) / 401
6. Docker build thành công

---

## Checklist hoàn thành

| # | Yêu cầu | Trạng thái | Ghi chú |
|---|---------|------------|---------|
| 1 | Template load parse từng cặp (KHÔNG dùng LoadHTMLGlob) | 🔲 | |
| 2 | Template dùng cfg.TemplateDir | 🔲 | |
| 3 | 10 routes đăng ký | 🔲 | |
| 4 | Health check ping DB | 🔲 | |
| 5 | Graceful shutdown | 🔲 | |
| 6 | Dockerfile copy templates đúng path | 🔲 | |
| 7 | fly.toml có APP_ENV + TEMPLATE_DIR | 🔲 | |
| 8 | fly.toml KHÔNG chứa secrets | 🔲 | |
| 9 | E2E flow | 🔲 | |

---

## Report

> Điền sau khi hoàn thành
