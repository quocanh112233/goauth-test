# Prompt 12 — Echo: Router + Dockerfile + fly.toml

## Role

Bạn là một **Go DevOps Engineer** thành thạo Echo v4 và Fly.io.

---

## Context

Prompt 11 đã có đủ layers. Hoàn thiện app.

---

## Dependencies

| Prompt | Đầu ra cần thiết |
|--------|-----------------|
| 01 | `shared/templates/*.html` |
| 11 | Echo: tất cả layers |

---

## Yêu cầu

### 1. echo/internal/router/router.go

- `e.Renderer = renderer.NewRenderer(cfg.TemplateDir)`
- `e.Use(echomiddleware.Logger(), echomiddleware.Recover())`
- 10 routes (giống pattern Gin prompt05)

### 2. Health Check

`{"status":"ok","framework":"Echo","db":"connected"}`

### 3. fly.toml

```toml
app = "goauth-echo"

[env]
  PORT = "8080"
  MONGO_DB = "goauth"
  APP_ENV = "production"
  TEMPLATE_DIR = "./shared/templates"
```

### 4. Graceful shutdown: `e.Shutdown(ctx)`

---

## Acceptance Criteria

1. E2E: Signup → Login → Home → Logout
2. Docker build thành công
3. fly.toml có APP_ENV + TEMPLATE_DIR

---

## Checklist hoàn thành

| # | Yêu cầu | Trạng thái | Ghi chú |
|---|---------|------------|---------|
| 1 | Renderer dùng cfg.TemplateDir | 🔲 | |
| 2 | 10 routes | 🔲 | |
| 3 | Health check | 🔲 | |
| 4 | Graceful shutdown | 🔲 | |
| 5 | fly.toml: APP_ENV + TEMPLATE_DIR | 🔲 | |
| 6 | E2E flow | 🔲 | |

---

## Report

> Điền sau khi hoàn thành
