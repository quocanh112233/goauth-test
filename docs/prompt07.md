# Prompt 07 — Fiber: Router + Dockerfile + fly.toml

## Role

Bạn là một **Go DevOps Engineer** thành thạo Fiber framework và Fly.io.

---

## Context

Prompt 06 đã có đủ layers cho Fiber. Prompt này hoàn thiện: Router, Health Check, Dockerfile, fly.toml.

Route mapping chi tiết xem `docs/api-spec.md`.

---

## Dependencies (Prompt phụ thuộc)

| Prompt | Đầu ra cần thiết |
|--------|-----------------|
| 01 | `shared/templates/*.html` |
| 06 | Fiber: tất cả layers |

---

## Yêu cầu

### 1. fiber/internal/router/router.go

- Template engine: `html.New("../shared/templates", ".html")`
- Đăng ký routes:
  ```
  GET  /              → redirect /login
  GET  /health        → healthHandler
  GET  /login         → authHandler.ShowLogin
  POST /login         → authHandler.Login
  GET  /signup        → authHandler.ShowSignup
  POST /signup        → authHandler.Signup
  GET  /api/me        → apiHandler.GetMe           [auth middleware]
  GET  /home          → homeHandler.ShowHome        [auth middleware]
  GET  /dashboard     → dashHandler.ShowDashboard   [auth middleware]
  POST /logout        → authHandler.Logout          [auth middleware]
  ```

### 2. Health Check

GET `/health` → `{"status": "ok", "framework": "Fiber", "db": "connected"}`

### 3. fiber/cmd/main.go

- Graceful shutdown: `app.ShutdownWithTimeout(5 * time.Second)`

### 4–6. Dockerfile + fly.toml + .dockerignore

- Giống pattern Gin (Prompt 05)
- fly.toml: `app = "goauth-fiber"`, `primary_region = "sin"`

---

## Acceptance Criteria

1. `cd fiber && go build ./...` pass
2. E2E: Signup → Login → Home → Logout
3. Admin Login → Dashboard
4. `/api/me` → JSON user info
5. Docker build thành công

---

## Checklist hoàn thành

| # | Yêu cầu | Trạng thái | Ghi chú |
|---|---------|------------|---------|
| 1 | Router đủ 10 routes | 🔲 | |
| 2 | Health check ping DB | 🔲 | |
| 3 | Auth routes dùng middleware | 🔲 | |
| 4 | Graceful shutdown | 🔲 | |
| 5 | Dockerfile + copy templates | 🔲 | |
| 6 | fly.toml health check + no secrets | 🔲 | |
| 7 | E2E flow hoạt động | 🔲 | |

---

## Report

> Điền sau khi hoàn thành
