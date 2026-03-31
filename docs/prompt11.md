# Prompt 11 — Echo: Router + Dockerfile + fly.toml

## Role

Bạn là một **Go DevOps Engineer** thành thạo Echo v4 và Fly.io.

---

## Context

Prompt 10 đã có đủ layers cho Echo. Prompt này hoàn thiện app.

---

## Dependencies (Prompt phụ thuộc)

| Prompt | Đầu ra cần thiết |
|--------|-----------------|
| 01 | `shared/templates/*.html` |
| 10 | Echo: tất cả layers |

---

## Yêu cầu

### 1. echo/internal/router/router.go

- `e.Renderer = renderer.NewRenderer("../shared/templates")`
- `e.Use(echomiddleware.Logger(), echomiddleware.Recover())`
- Đăng ký routes:
  ```
  GET  /              → redirect /login
  GET  /health        → healthHandler
  GET  /login         → authHandler.ShowLogin
  POST /login         → authHandler.Login
  GET  /signup        → authHandler.ShowSignup
  POST /signup        → authHandler.Signup

  // Protected
  GET  /home          → homeHandler.ShowHome
  GET  /dashboard     → dashHandler.ShowDashboard
  POST /logout        → authHandler.Logout
  GET  /api/me        → apiHandler.GetMe
  ```

### 2. Health Check

`{"status":"ok","framework":"Echo","db":"connected"}`

### 3–6. main.go + Dockerfile + fly.toml + .dockerignore

- fly.toml: `app = "goauth-echo"`
- Graceful shutdown: `e.Shutdown(ctx)`

---

## Acceptance Criteria

1. `cd echo && go build ./...` pass
2. E2E: Signup → Login → Home → Logout
3. Admin Login → Dashboard
4. `/api/me` → JSON
5. Docker build thành công

---

## Checklist hoàn thành

| # | Yêu cầu | Trạng thái | Ghi chú |
|---|---------|------------|---------|
| 1 | Set e.Renderer | 🔲 | |
| 2 | Logger + Recover middleware | 🔲 | |
| 3 | 10 routes | 🔲 | |
| 4 | Health check | 🔲 | |
| 5 | Auth middleware trên protected routes | 🔲 | |
| 6 | Graceful shutdown | 🔲 | |
| 7 | Dockerfile + fly.toml | 🔲 | |
| 8 | E2E flow | 🔲 | |

---

## Report

> Điền sau khi hoàn thành
