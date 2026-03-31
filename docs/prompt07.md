# Prompt 07 — Fiber: Router + Dockerfile + fly.toml

## Role

Bạn là một **Go DevOps Engineer** thành thạo Fiber framework và Fly.io.

---

## Context

Prompt 06 đã có đủ layers. Hoàn thiện app. Conventions xem `docs/conventions.md`.

---

## Dependencies

| Prompt | Đầu ra cần thiết |
|--------|-----------------|
| 01 | `shared/templates/*.html` |
| 06 | Fiber: tất cả layers |

---

## Yêu cầu

### 1. fiber/internal/router/router.go

Template engine:
```go
// Dùng cfg.TemplateDir thay vì hardcode
engine := html.New(cfg.TemplateDir, ".html")
app := fiber.New(fiber.Config{Views: engine})
```

> ⚠️ **Docker path**: Local = `../shared/templates` (default), Docker = `./shared/templates`. Đã xử lý qua `cfg.TemplateDir` + env `TEMPLATE_DIR` (xem `conventions.md` mục 7).

Routes (10 routes — giống pattern Gin prompt05).

### 2. Health Check

`{"status":"ok","framework":"Fiber","db":"connected"}`

### 3. Dockerfile + fly.toml

fly.toml:
```toml
app = "goauth-fiber"

[env]
  PORT = "8080"
  MONGO_DB = "goauth"
  APP_ENV = "production"
  TEMPLATE_DIR = "./shared/templates"
```

> `TEMPLATE_DIR = "./shared/templates"` — giải quyết vấn đề path trong Docker.

---

## Anti-Patterns

❌ Không hardcode template path `"../shared/templates"` trong code — dùng `cfg.TemplateDir`
❌ Không quên `TEMPLATE_DIR` trong fly.toml

---

## Acceptance Criteria

1. E2E: Signup → Login → Home → Logout
2. Docker build thành công
3. Template path giải quyết đúng cả local + Docker

---

## Checklist hoàn thành

| # | Yêu cầu | Trạng thái | Ghi chú |
|---|---------|------------|---------|
| 1 | Template engine dùng cfg.TemplateDir | 🔲 | |
| 2 | 10 routes | 🔲 | |
| 3 | fly.toml có TEMPLATE_DIR + APP_ENV | 🔲 | |
| 4 | E2E flow | 🔲 | |

---

## Report

> Điền sau khi hoàn thành
