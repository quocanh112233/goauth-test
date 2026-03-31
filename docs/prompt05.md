# Prompt 05 — Gin: Router + Dockerfile + fly.toml

## Role

Bạn là một **Go DevOps Engineer** kiêm Backend Developer. Bạn wire tất cả components thành app hoàn chỉnh, đóng gói Docker, và config deploy Fly.io.

---

## Context

Prompt 03 + 04 đã có đủ layers. Prompt này kết nối tất cả: Router, Health Check, main.go, Dockerfile, fly.toml.

Route mapping chi tiết xem `docs/api-spec.md` (section "Route Summary").

---

## Dependencies (Prompt phụ thuộc)

| Prompt | Đầu ra cần thiết |
|--------|-----------------|
| 01 | `shared/templates/*.html` (6 files) |
| 03 | Config, DB, Model (User + Session), Repository (User + Session) |
| 04 | Service, Handler (Auth + Home + Dashboard + API), Middleware |

---

## Yêu cầu

### 1. gin/internal/router/router.go

Hàm `Setup(cfg *config.Config, database *mongo.Database) *gin.Engine`:

- Tạo Gin engine: `gin.Default()`
- Load HTML templates: parse mỗi page riêng với base.html
- Khởi tạo: userRepo → sessionRepo → service → handlers → middleware
- Đăng ký routes:

  ```
  GET  /                      → redirect /login (302)
  GET  /health                → healthHandler.HealthCheck
  GET  /login                 → authHandler.ShowLogin
  POST /login                 → authHandler.Login
  GET  /signup                → authHandler.ShowSignup
  POST /signup                → authHandler.Signup
  GET  /api/me                → apiHandler.GetMe              [auth middleware]

  // Protected group — áp dụng RequireAuth middleware
  GET  /home                  → homeHandler.ShowHome           [auth middleware]
  GET  /dashboard             → dashboardHandler.ShowDashboard [auth middleware]
  POST /logout                → authHandler.Logout             [auth middleware]
  ```

> **Lưu ý**: `/api/me` cũng dùng auth middleware nhưng middleware sẽ return 401 JSON thay vì redirect (xem Prompt 04 middleware).

### 2. Health Check Endpoint

GET `/health` trả về JSON:

```json
{
  "status": "ok",
  "framework": "Gin",
  "db": "connected"
}
```

- Ping MongoDB → nếu fail: `"db": "disconnected"`, status 503
- Endpoint này **không cần auth**

### 3. gin/cmd/main.go

- Load config → `db.Connect(cfg)` → setup router
- Graceful shutdown:
  - Goroutine chạy server
  - Main goroutine chờ SIGINT/SIGTERM
  - `srv.Shutdown(ctx)` timeout 5 giây
  - `db.Disconnect(client)`
- Log: `"✓ Connected to MongoDB"`, `"✓ Gin server running on :PORT"`

### 4. gin/Dockerfile

Multi-stage build:

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

### 5. gin/fly.toml

```toml
app = "goauth-gin"
primary_region = "sin"

[build]

[env]
  PORT = "8080"
  MONGO_DB = "goauth"

[http_service]
  internal_port = 8080
  force_https = true
  auto_stop_machines = true
  auto_start_machines = true
  min_machines_running = 0

[[http_service.checks]]
  grace_period = "10s"
  interval = "30s"
  method = "GET"
  path = "/health"
  timeout = "5s"

[[vm]]
  memory = "256mb"
  cpu_kind = "shared"
  cpus = 1
```

> Secrets: `fly secrets set MONGO_URI="..." JWT_SECRET="..."`

### 6. gin/.dockerignore

```
.env
*.md
.git
tmp/
```

---

## Anti-Patterns (KHÔNG được làm)

❌ Không dùng `gin.LoadHTMLGlob("*.html")` — blocks bị override
❌ Không để secrets trong fly.toml
❌ Không quên copy templates vào Docker image
❌ Không quên `/api/me` route dùng auth middleware

---

## Acceptance Criteria

1. `cd gin && go build ./...` pass
2. Truy cập `http://localhost:8081` → redirect `/login`
3. `/health` → JSON status ok
4. Signup → Login → Home (user) hoạt động
5. Login admin → Dashboard
6. Logout → redirect /login, không truy cập được /home
7. `/api/me` + valid cookie → JSON user info (không có password)
8. `/api/me` + no cookie → 401 JSON
9. `docker build -f gin/Dockerfile .` → thành công

---

## Checklist hoàn thành

| # | Yêu cầu | Trạng thái | Ghi chú |
|---|---------|------------|---------|
| 1 | Router setup đủ 10 routes | 🔲 | |
| 2 | / redirect /login | 🔲 | |
| 3 | /health JSON + ping DB | 🔲 | |
| 4 | /home + /dashboard + /logout + /api/me dùng auth middleware | 🔲 | |
| 5 | Load templates đúng cách | 🔲 | |
| 6 | main.go graceful shutdown | 🔲 | |
| 7 | Dockerfile multi-stage + copy templates | 🔲 | |
| 8 | fly.toml app = "goauth-gin" | 🔲 | |
| 9 | fly.toml health check config | 🔲 | |
| 10 | fly.toml KHÔNG chứa secrets | 🔲 | |
| 11 | E2E: Signup → Login → Home → Logout | 🔲 | |
| 12 | E2E: Admin Login → Dashboard | 🔲 | |
| 13 | API: /api/me trả JSON (no password) | 🔲 | |

---

## Report

> Điền sau khi hoàn thành
