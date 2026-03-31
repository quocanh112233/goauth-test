# Prompt 10 — net/http (stdlib): Router + Dockerfile + fly.toml

## Role

Bạn là một **Go DevOps Engineer** thành thạo Go `net/http` server.

---

## Context

Prompt 08 (data) + 09 (logic) đã có đủ layers. Prompt này hoàn thiện: Router (Go 1.22+ ServeMux), Health Check, Dockerfile, fly.toml.

---

## Dependencies

| Prompt | Đầu ra cần thiết |
|--------|-----------------|
| 01 | `shared/templates/*.html` |
| 08 | Config, DB, Model, Repository |
| 09 | Service, Handler, Middleware, TemplateMap |

---

## Yêu cầu

### 1. stdlib/internal/router/router.go

Go 1.22+ method routing:

```go
mux := http.NewServeMux()

// Public
mux.HandleFunc("GET /", redirectLogin)
mux.HandleFunc("GET /health", healthHandler.HealthCheck)
mux.HandleFunc("GET /login", authHandler.ShowLogin)
mux.HandleFunc("POST /login", authHandler.Login)
mux.HandleFunc("GET /signup", authHandler.ShowSignup)
mux.HandleFunc("POST /signup", authHandler.Signup)

// Protected — wrap middleware
authMW := middleware.RequireAuth(service, jwtSecret, cfg.IsProduction)
mux.Handle("GET /home", middleware.Chain(http.HandlerFunc(homeHandler.ShowHome), authMW))
mux.Handle("GET /dashboard", middleware.Chain(http.HandlerFunc(dashHandler.ShowDashboard), authMW))
mux.Handle("POST /logout", middleware.Chain(http.HandlerFunc(authHandler.Logout), authMW))
mux.Handle("GET /api/me", middleware.Chain(http.HandlerFunc(apiHandler.GetMe), authMW))
```

### 2. Health Check

`{"status":"ok","framework":"net/http","db":"connected"}`

### 3. Graceful shutdown

```go
srv := &http.Server{Addr: ":" + cfg.Port, Handler: mux}
go func() { srv.ListenAndServe() }()
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
srv.Shutdown(ctx)
db.Disconnect(client)
```

### 4. Dockerfile + fly.toml

fly.toml:
```toml
app = "goauth-stdlib"

[env]
  PORT = "8080"
  MONGO_DB = "goauth"
  APP_ENV = "production"
  TEMPLATE_DIR = "./shared/templates"
```

---

## Acceptance Criteria

1. `cd stdlib && go build ./...` pass
2. `go.mod` KHÔNG chứa framework
3. E2E: Signup → Login → Home → Logout
4. Docker build thành công

---

## Checklist hoàn thành

| # | Yêu cầu | Trạng thái | Ghi chú |
|---|---------|------------|---------|
| 1 | Go 1.22+ method routing | 🔲 | |
| 2 | 10 routes | 🔲 | |
| 3 | Protected routes wrap middleware.Chain | 🔲 | |
| 4 | Health check | 🔲 | |
| 5 | Graceful shutdown | 🔲 | |
| 6 | fly.toml: APP_ENV + TEMPLATE_DIR | 🔲 | |
| 7 | E2E flow | 🔲 | |

---

## Report

> Điền sau khi hoàn thành
