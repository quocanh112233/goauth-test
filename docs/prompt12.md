# Prompt 12 — GitHub Actions CI/CD + Makefile hoàn chỉnh

## Role

Bạn là một **DevOps Engineer** thành thạo GitHub Actions và Fly.io CLI.

---

## Context

Phase 6. Tất cả 4 apps hoàn chỉnh. Setup CI/CD tự động deploy khi push + Makefile cho dev workflow.

---

## Dependencies (Prompt phụ thuộc)

| Prompt | Đầu ra cần thiết |
|--------|-----------------|
| 05 | Gin Dockerfile + fly.toml |
| 07 | Fiber Dockerfile + fly.toml |
| 09 | stdlib Dockerfile + fly.toml |
| 11 | Echo Dockerfile + fly.toml |

---

## Yêu cầu

### 1–4. GitHub Actions Workflows

4 files: `deploy-gin.yml`, `deploy-fiber.yml`, `deploy-stdlib.yml`, `deploy-echo.yml`

Mỗi file:
- Trigger: push to `main`, paths: `<approach>/**` + `shared/**`
- Token riêng: `FLY_API_TOKEN_<APPROACH>`
- `flyctl deploy --remote-only`

### 5. Makefile

```makefile
# Local Development
run-gin:     ## Port 8081
run-fiber:   ## Port 8082
run-stdlib:  ## Port 8083
run-echo:    ## Port 8084
run-all:     ## Song song + trap 'kill 0' EXIT

# Seed
seed:        ## Seed admin (scripts/seed.go)

# Deploy
deploy-gin deploy-fiber deploy-stdlib deploy-echo deploy-all

# Secrets
secrets-gin secrets-fiber secrets-stdlib secrets-echo secrets-all

# Test
test-gin test-fiber test-stdlib test-echo test-all test-coverage

# Utility
tidy:        ## go mod tidy cho tất cả modules
help:        ## Hiển thị help
```

### 6. README.md — section CI/CD

- Hướng dẫn set GitHub Secrets
- Lấy Fly.io token: `fly tokens create deploy -a <app>`

---

## Anti-Patterns (KHÔNG được làm)

❌ Không dùng chung 1 FLY_API_TOKEN
❌ Không quên `shared/**` trong paths filter
❌ Không dùng space thay tab trong Makefile

---

## Acceptance Criteria

1. Push `gin/` → chỉ trigger deploy-gin
2. Push `shared/` → trigger cả 4
3. `make help` hiển thị targets
4. `make run-all` chạy 4 server song song

---

## Checklist hoàn thành

| # | Yêu cầu | Trạng thái | Ghi chú |
|---|---------|------------|---------|
| 1 | 4 workflow files đúng path filter | 🔲 | |
| 2 | Mỗi workflow token riêng | 🔲 | |
| 3 | shared/** trigger cả 4 | 🔲 | |
| 4 | Makefile đủ targets (dev, deploy, test, utility) | 🔲 | |
| 5 | Makefile run-all parallel + trap | 🔲 | |
| 6 | README CI/CD section | 🔲 | |

---

## Report

> Điền sau khi hoàn thành
