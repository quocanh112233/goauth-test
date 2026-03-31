# Prompt 09 — net/http (stdlib): Router + Dockerfile + fly.toml

## Role

Bạn là một **Go DevOps Engineer** thành thạo `net/http` server và Fly.io.

---

## Context

Prompt 08 đã có đủ layers cho stdlib. Prompt này hoàn thiện: Router (Go 1.22+ ServeMux), Health Check, Dockerfile, fly.toml.

---

## Dependencies (Prompt phụ thuộc)

| Prompt | Đầu ra cần thiết |
|--------|-----------------|
| 01 | `shared/templates/*.html` |
| 08 | stdlib: tất cả layers |

---

## Yêu cầu

### 1. stdlib/internal/router/router.go

Dùng Go 1.22+ method routing:

```go
mux := http.NewServeMux()

// Public routes
mux.HandleFunc("GET /", redirectLogin)
mux.HandleFunc("GET /health", healthHandler.HealthCheck)
mux.HandleFunc("GET /login", authHandler.ShowLogin)
mux.HandleFunc("POST /login", authHandler.Login)
mux.HandleFunc("GET /signup", authHandler.ShowSignup)
mux.HandleFunc("POST /signup", authHandler.Signup)

// Protected routes — wrap với middleware.Chain
authMW := middleware.RequireAuth(service, jwtSecret)
mux.Handle("GET /home", middleware.Chain(http.HandlerFunc(homeHandler.ShowHome), authMW))
mux.Handle("GET /dashboard", middleware.Chain(http.HandlerFunc(dashHandler.ShowDashboard), authMW))
mux.Handle("POST /logout", middleware.Chain(http.HandlerFunc(authHandler.Logout), authMW))
mux.Handle("GET /api/me", middleware.Chain(http.HandlerFunc(apiHandler.GetMe), authMW))
```

### 2. Health Check

`{"status":"ok","framework":"net/http","db":"connected"}`

### 3. Graceful shutdown chuẩn Go

```go
srv := &http.Server{Addr: ":"+cfg.Port, Handler: router}
go func() { srv.ListenAndServe() }()
// signal → srv.Shutdown(ctx)
```

### 4–6. Dockerfile + fly.toml + .dockerignore

- fly.toml: `app = "goauth-stdlib"`

---

## Acceptance Criteria

1. `cd stdlib && go build ./...` pass
2. `go.mod` KHÔNG chứa framework
3. E2E: Signup → Login → Home → Logout
4. Admin Login → Dashboard
5. `/api/me` → JSON
6. Docker build thành công

---

## Checklist hoàn thành

| # | Yêu cầu | Trạng thái | Ghi chú |
|---|---------|------------|---------|
| 1 | Router Go 1.22+ method routing | 🔲 | |
| 2 | 10 routes đăng ký | 🔲 | |
| 3 | Protected routes wrap middleware.Chain | 🔲 | |
| 4 | Health check | 🔲 | |
| 5 | Graceful shutdown | 🔲 | |
| 6 | Dockerfile + fly.toml | 🔲 | |
| 7 | E2E flow | 🔲 | |

---

## Report

> Điền sau khi hoàn thành
